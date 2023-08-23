package main

import (
	"context"
	"flag"
	"fmt"

	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"

	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/handler/browse"
	handlertag "github.com/wutipong/mangaweb/handler/tag"
	"github.com/wutipong/mangaweb/handler/view"
	"github.com/wutipong/mangaweb/log"
	"github.com/wutipong/mangaweb/meta"
	"github.com/wutipong/mangaweb/scheduler"
	"github.com/wutipong/mangaweb/tag"
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
	connectionStr := setupFlag("database", "postgres://postgres:password@localhost:5432/manga", "MANGAWEB_DB", "Specify the database connection string")
	prefix := setupFlag("prefix", "", "MANGAWEB_PREFIX", "URL prefix")

	flag.Parse()

	meta.BaseDirectory = *dataPath
	printBanner()
	log.Get().Sugar().Infof("MangaWeb version:%s", versionString)

	log.Get().Sugar().Infof("Data source Path: %s", *dataPath)

	router := httprouter.New()

	conn, err := pgxpool.New(context.Background(), *connectionStr)
	if err != nil {
		log.Get().Sugar().Fatal(err)

		return
	}

	defer conn.Close()

	tag.Init(context.Background(), conn)
	meta.Init(context.Background(), conn)

	scheduler.Init(scheduler.Options{})

	RegisterHandler(router, *prefix)
	scheduler.Start()

	log.Get().Info("Server starts.")
	log.Get().Sugar().Fatal(http.ListenAndServe(*address, router))

	log.Get().Sugar().Info("shutting down the server")
	scheduler.Stop()
}

func RegisterHandler(router *httprouter.Router, pathPrefix string) {
	handler.Init(handler.Options{
		VersionString:     versionString,
		PathPrefix:        pathPrefix,
		PathRoot:          pathRoot,
		PathBrowse:        pathBrowse,
		PathView:          pathView,
		PathStatic:        pathStatic,
		PathGetImage:      pathGetImage,
		PathUpdateCover:   pathUpdateCover,
		PathThumbnail:     pathThumbnail,
		PathFavorite:      pathFavorite,
		PathDownload:      pathDownload,
		PathRescanLibrary: pathRescanLibrary,
		PathTagFavorite:   pathTagFavorite,
		PathTagList:       pathTagList,
		PathTagThumbnail:  pathTagThumb,
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
