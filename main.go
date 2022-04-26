package main

import (
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/tag"

	"net/http"
	"os"
	"path"

	"github.com/joho/godotenv"
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

	router := httprouter.New()

	scheduler.Init(scheduler.Options{
		MetaProviderFactory: newMetaProvider,
		TagProviderFactory:  newTagProvider,
	})

	RegisterHandler(router, *prefix)
	scheduler.Start()

	log.Info("Server starts.")
	log.Fatal(http.ListenAndServe(*address, router))

	log.Info("shutting down the server")
	scheduler.Stop()
}

func RegisterHandler(router *httprouter.Router, pathPrefix string) {
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
		PathTagThumbnail:    pathTagThumb,
	})
	// Routes
	router.GET(pathRoot, root)
	router.GET(path.Join(pathBrowse, "*tag"), browse.Handler)
	router.GET(path.Join(pathView, "*item"), view.Handler)
	router.GET(path.Join(pathGetImage, "*item"), handler.GetImage)
	router.GET(path.Join(pathUpdateCover, "*item"), view.UpdateCover)
	router.GET(path.Join(pathThumbnail, "*item"), browse.ThumbnailHandler)
	router.GET(path.Join(pathFavorite, "*item"), view.SetFavoriteHandler)
	router.GET(path.Join(pathDownload, "*item"), view.Download)
	router.GET(pathRescanLibrary, handler.RescanLibraryHandler)
	router.GET(path.Join(pathTagFavorite, "*tag"), handlertag.SetFavoriteHandler)
	router.GET(pathTagList, handlertag.TagListHandler)
	router.GET(path.Join(pathTagThumb, "*tag"), handlertag.ThumbnailHandler)

	router.ServeFiles(path.Join(pathStatic, "*filepath"), http.Dir("static"))
}

func root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, handler.CreateBrowseTagURL(""), http.StatusPermanentRedirect)
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
