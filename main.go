package main

import (
	"flag"
	"fmt"
	"github.com/wutipong/mangaweb/scheduler"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"path"

	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/handler/browse"
	"github.com/wutipong/mangaweb/handler/view"
	"github.com/wutipong/mangaweb/meta"
	"github.com/wutipong/mangaweb/meta/mongo"
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

	handler.Init(handler.Options{
		MetaProviderFactory: newProvider,
		VersionString:       versionString,
	})

	scheduler.Init(newProvider)

	e.Pre(middleware.RemoveTrailingSlash())

	// Routes
	e.GET("/", root)
	e.GET("/browse", browse.Handler)
	e.GET("/browse/*", browse.Handler)

	e.GET("/view", view.Handler)
	e.GET("/view/*", view.Handler)

	e.Static("/static", "static")

	e.GET("/get_image/*", handler.GetImage)
	e.GET("/update_cover/*", handler.UpdateCover)

	e.GET("/thumbnail/*", handler.ThumbnailHandler)

	e.GET("/favorite", handler.SetFavoriteHandler)
	e.GET("/favorite/*", handler.SetFavoriteHandler)

	e.GET("/download/*", handler.Download)

	e.GET("/rescan_library", handler.RescanLibraryHandler)

	scheduler.Start()

	log.Info("Server starts.")
	if err := e.Start(*address); err != http.ErrServerClosed {
		log.Error(err)
	}
	log.Info("shutting down the server")
	scheduler.Stop()
}

// Handler
func root(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, util.CreateURL("/browse"))
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
