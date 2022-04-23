package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/handler/browse"
	"github.com/wutipong/mangaweb/handler/view"
	"github.com/wutipong/mangaweb/meta"
	"github.com/wutipong/mangaweb/meta/mongo"
	"github.com/wutipong/mangaweb/scheduler"
	"github.com/wutipong/mangaweb/util"
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
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Info("Use .env file.")
	}

	address := setupFlag("address", ":80", "MANGAWEB_ADDRESS", "The server address")
	dataPath := setupFlag("data", "./data", "MANGAWEB_DATA_PATH", "Manga source path")
	database := setupFlag("database", "mongodb://root:password@localhost", "MANGAWEB_DB", "Specify the database connection string")
	prefix := setupFlag("prefix", "", "MANGAWEB_PREFIX", "URL prefix")

	flag.Parse()

	meta.BaseDirectory = *dataPath
	printBanner()
	log.Infof("MangaWeb version:%s", versionString)

	log.Infof("Data source Path: %s", *dataPath)
	if err := mongo.Init(*database); err != nil {
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

		util.SetPrefix(*prefix)
	}

	scheduler.Init(newProvider)

	e.Pre(middleware.RemoveTrailingSlash())

	RegisterHandler(e)
	scheduler.Start()

	log.Info("Server starts.")
	if err := e.Start(*address); err != http.ErrServerClosed {
		log.Error(err)
	}
	log.Info("shutting down the server")
	scheduler.Stop()
}

func RegisterHandler(e *echo.Echo) {
	handler.Init(handler.Options{
		MetaProviderFactory: newProvider,
		VersionString:       versionString,
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
	})
	// Routes
	e.GET(pathRoot, root)
	e.GET(pathBrowse, browse.Handler)
	e.GET(path.Join(pathView, "*"), view.Handler)
	e.GET(path.Join(pathGetImage, "*"), handler.GetImage)
	e.GET(path.Join(pathUpdateCover, "*"), handler.UpdateCover)
	e.GET(path.Join(pathThumbnail, "*"), handler.ThumbnailHandler)
	e.GET(path.Join(pathFavorite, "*"), handler.SetFavoriteHandler)
	e.GET(path.Join(pathFavorite, "*"), handler.SetFavoriteHandler)
	e.GET(path.Join(pathDownload, "*"), handler.Download)
	e.GET(pathRescanLibrary, handler.RescanLibraryHandler)

	e.Static(pathStatic, "static")
}

// Handler
func root(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, util.CreateURL(pathBrowse))
}

func newProvider() (p meta.Provider, err error) {
	mp, e := mongo.New()
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
