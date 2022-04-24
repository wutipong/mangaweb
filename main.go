package main

import (
	"flag"
	"fmt"
	"github.com/wutipong/mangaweb/tag"

	"net/http"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/handler/browse"
	handlertag "github.com/wutipong/mangaweb/handler/tag"
	"github.com/wutipong/mangaweb/handler/view"
	"github.com/wutipong/mangaweb/meta"
	metamongo "github.com/wutipong/mangaweb/meta/mongo"
	"github.com/wutipong/mangaweb/scheduler"
	tagmongo "github.com/wutipong/mangaweb/tag/mongo"
)

// Recreate the static resource file.
//go:generate npm install
//go:generate npm run build

func setupFlag(flagName, defValue, variable, description string) *string {
	varValue := os.Getenv(variable)
	if varValue != "" {
		defValue = varValue
	}

	return flag.String(flagName, defValue, description)
}

var versionString string = "development"

const (
	pathRoot          = "/"
	pathBrowse        = "/browse"
	pathView          = "/view"
	pathStatic        = "/static"
	pathGetImage      = "/get_image"
	pathUpdateCover   = "/update_cover"
	pathThumbnail     = "/thumbnail"
	pathFavorite      = "/favorite"
	pathDownload      = "/download"
	pathRescanLibrary = "/rescan_library"
	pathTagFavorite   = "/tag_favorite"
	pathTagList       = "/tag_list"
	pathTagThumb      = "/tag_thumb"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Info("Use .env file.")
	}

	address := setupFlag("address", ":80", "MANGAWEB_ADDRESS", "The server address")
	dataPath := setupFlag("data", "./data", "MANGAWEB_DATA_PATH", "Manga source path")
	database := setupFlag("database", "mongodb://root:password@localhost", "MANGAWEB_DB", "Specify the database connection string")
	dbName := setupFlag("database_name", "manga", "MANGAWEB_DB_NAME", "Specify the database name")
	prefix := setupFlag("prefix", "", "MANGAWEB_PREFIX", "URL prefix")

	flag.Parse()

	meta.BaseDirectory = *dataPath
	printBanner()
	log.Infof("MangaWeb version:%s", versionString)

	log.Infof("Data source Path: %s", *dataPath)
	if err := metamongo.Init(*database, *dbName); err != nil {
		log.Fatal(err)
	}

	if err := tagmongo.Init(*database, *dbName); err != nil {
		log.Fatal(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if *prefix != "" {
		pattern := path.Join(*prefix, "*")
		e.Pre(middleware.Rewrite(map[string]string{
			*prefix: "/",
			pattern: "/$1",
		}))
	}

	scheduler.Init(scheduler.Options{
		MetaProviderFactory: newMetaProvider,
		TagProviderFactory:  newTagProvider,
	})

	e.Pre(middleware.RemoveTrailingSlash())

	RegisterHandler(e, *prefix)
	scheduler.Start()

	log.Info("Server starts.")
	if err := e.Start(*address); err != http.ErrServerClosed {
		log.Error(err)
	}
	log.Info("shutting down the server")
	scheduler.Stop()
}

func RegisterHandler(e *echo.Echo, pathPrefix string) {
	handler.Init(handler.Options{
		MetaProviderFactory: newMetaProvider,
		TagProviderFactory:  newTagProvider,
		VersionString:       versionString,
		PathPrefix:          pathPrefix,
		PathRoot:            pathRoot,
		PathBrowse:          pathBrowse,
		PathView:            pathView,
		PathStatic:          pathStatic,
		PathGetImage:        pathGetImage,
		PathUpdateCover:     pathUpdateCover,
		PathThumbnail:       pathThumbnail,
		PathFavorite:        pathFavorite,
		PathDownload:        pathDownload,
		PathRescanLibrary:   pathRescanLibrary,
		PathTagFavorite:     pathTagFavorite,
		PathTagList:         pathTagList,
		PathTagThumb:        pathTagThumb,
	})
	// Routes
	e.GET(pathRoot, root)
	e.GET(pathBrowse, browse.Handler)
	e.GET(path.Join(pathBrowse, "*"), browse.Handler)
	e.GET(path.Join(pathView, "*"), view.Handler)
	e.GET(path.Join(pathGetImage, "*"), handler.GetImage)
	e.GET(path.Join(pathUpdateCover, "*"), view.UpdateCover)
	e.GET(path.Join(pathThumbnail, "*"), browse.ThumbnailHandler)
	e.GET(path.Join(pathFavorite, "*"), view.SetFavoriteHandler)
	e.GET(path.Join(pathFavorite, "*"), view.SetFavoriteHandler)
	e.GET(path.Join(pathDownload, "*"), view.Download)
	e.GET(pathRescanLibrary, handler.RescanLibraryHandler)
	e.GET(path.Join(pathTagFavorite, "*"), handlertag.SetFavoriteHandler)
	e.GET(pathTagList, handlertag.Handler)
	e.GET(path.Join(pathTagThumb, "*"), handlertag.ThumbnailHandler)

	e.Static(pathStatic, "static")
}

// Handler
func root(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, handler.CreateBrowseURL(""))
}

func newMetaProvider() (p meta.Provider, err error) {
	mp, e := metamongo.New()
	p = &mp
	err = e
	return
}

func newTagProvider() (p tag.Provider, err error) {
	mp, e := tagmongo.New()
	p = &mp
	err = e
	return
}

func printBanner() {
	fmt.Printf(
		`
	███╗░░░███╗░█████╗░███╗░░██╗░██████╗░░█████╗░░██╗░░░░░░░██╗███████╗██████╗░
	████╗░████║██╔══██╗████╗░██║██╔════╝░██╔══██╗░██║░░██╗░░██║██╔════╝██╔══██╗
	██╔████╔██║███████║██╔██╗██║██║░░██╗░███████║░╚██╗████╗██╔╝█████╗░░██████╦╝
	██║╚██╔╝██║██╔══██║██║╚████║██║░░╚██╗██╔══██║░░████╔═████║░██╔══╝░░██╔══██╗
	██║░╚═╝░██║██║░░██║██║░╚███║╚██████╔╝██║░░██║░░╚██╔╝░╚██╔╝░███████╗██████╦╝
	╚═╝░░░░░╚═╝╚═╝░░╚═╝╚═╝░░╚══╝░╚═════╝░╚═╝░░╚═╝░░░╚═╝░░░╚═╝░░╚══════╝╚═════╝░
	Version: %s`, versionString)
	fmt.Println()
}
