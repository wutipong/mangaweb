package main

import (
	"flag"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"

	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/handler/browse"
	handlertag "github.com/wutipong/mangaweb/handler/tag"
	"github.com/wutipong/mangaweb/handler/view"
	"github.com/wutipong/mangaweb/log"
	"github.com/wutipong/mangaweb/meta"
	metamongo "github.com/wutipong/mangaweb/meta/mongo"
	"github.com/wutipong/mangaweb/scheduler"
	"github.com/wutipong/mangaweb/tag"
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
	err := log.Init()
	if err != nil {
		panic(err)
	}
	defer log.Close()

	err = godotenv.Load()
	if err != nil {
		log.Get().Info("Use .env file.")
	}

	address := setupFlag("address", ":80", "MANGAWEB_ADDRESS", "The server address")
	dataPath := setupFlag("data", "./data", "MANGAWEB_DATA_PATH", "Manga source path")
	database := setupFlag("database", "mongodb://root:password@localhost", "MANGAWEB_DB", "Specify the database connection string")
	dbName := setupFlag("database_name", "manga", "MANGAWEB_DB_NAME", "Specify the database name")
	prefix := setupFlag("prefix", "", "MANGAWEB_PREFIX", "URL prefix")

	flag.Parse()

	meta.BaseDirectory = *dataPath
	printBanner()
	log.Get().Sugar().Infof("MangaWeb version:%s", versionString)

	log.Get().Sugar().Infof("Data source Path: %s", *dataPath)
	if err := metamongo.Init(*database, *dbName); err != nil {
		log.Get().Sugar().Fatal(err)
	}

	if err := tagmongo.Init(*database, *dbName); err != nil {
		log.Get().Sugar().Fatal(err)
	}

	router := httprouter.New()

	scheduler.Init(scheduler.Options{
		MetaProviderFactory: newMetaProvider,
		TagProviderFactory:  newTagProvider,
	})

	RegisterHandler(router, *prefix)
	scheduler.Start()

	log.Get().Info("Server starts.")
	log.Get().Sugar().Fatal(http.ListenAndServe(*address, router))

	log.Get().Sugar().Info("shutting down the server")
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
	router.GET(handler.CreateURL(pathRoot), root)
	router.GET(handler.CreateBrowseURLPattern(), browse.Handler)
	router.GET(handler.CreateViewURLPattern(), view.Handler)
	router.GET(handler.CreateGetImageURLPattern(), handler.GetImage)
	router.GET(handler.CreateUpdateCoverURLPattern(), view.UpdateCover)
	router.GET(handler.CreateThumbnailURLPattern(), browse.ThumbnailHandler)
	router.GET(handler.CreateSetFavoriteURLPattern(), view.SetFavoriteHandler)
	router.GET(handler.CreateDownloadURLPattern(), view.Download)
	router.GET(handler.CreateRescanURLPattern(), handler.RescanLibraryHandler)
	router.GET(handler.CreateSetTagFavoriteURLPattern(), handlertag.SetFavoriteHandler)
	router.GET(handler.CreateTagListURLPattern(), handlertag.TagListHandler)
	router.GET(handler.CreateTagThumbnailURLPattern(), handlertag.ThumbnailHandler)

	router.ServeFiles(handler.CreateURL(pathStatic, "*filepath"), http.Dir("static"))
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
