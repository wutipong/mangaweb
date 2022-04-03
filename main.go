package main

import (
	"flag"
	"fmt"

	"net/http"
	"os"
	"path"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

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
	rebuild := setupFlag("rebuild_thumbnail", "false", "MANGAWEB_REBUILD_THUMBNAIL", "force rebuild all thumbnail.")

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

	// Schedule the update metadata task to run every 30 minutes.
	s := gocron.NewScheduler(time.UTC)
	s.Every(30).Minutes().Do(func() {
		log.Info("Update metadata set.")
		synchronizeMetaData()
	})

	if *rebuild == "true" {
		s.Every(1).Millisecond().LimitRunsTo(1).Do(func() {
			log.Info("Force updating thumbnail")
			rebuildThumbnail()
		})
	}

	s.StartAsync()

	log.Info("Server starts.")
	if err := e.Start(*address); err != http.ErrServerClosed {
		log.Error(err)
	}
	log.Info("shutting down the server")
	s.Stop()
}

// Handler
func root(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, util.CreateURL("/browse"))
}

func synchronizeMetaData() error {
	provider, err := newProvider()

	if err != nil {
		return err
	}
	defer provider.Close()

	allMeta, err := provider.ReadAll()
	if err != nil {
		return err
	}

	files, err := meta.ListDir("")
	if err != nil {
		return err
	}

	for _, file := range files {
		found := false
		for _, m := range allMeta {
			if m.Name == file {
				found = true
				break
			}
		}
		if found {
			continue
		}

		log.Printf("Creating metadata for %s", file)

		item, err := meta.NewItem(file)
		if err != nil {
			log.Printf("Failed to create meta data : %v", err)
		}

		err = provider.Write(item)
		if err != nil {
			log.Printf("Failed to write meta data : %v", err)
		}
	}

	for _, m := range allMeta {
		found := false
		for _, file := range files {
			if m.Name == file {
				found = true
				break
			}
		}
		if found {
			continue
		}

		log.Printf("Deleting metadata for %s", m.Name)
		if err := provider.Delete(m); err != nil {
			log.Printf("Failed to delete meta for %s", m.Name)
		}

	}

	return nil
}

func rebuildThumbnail() error {
	provider, err := newProvider()

	if err != nil {
		return err
	}
	defer provider.Close()

	allMeta, err := provider.ReadAll()
	if err != nil {
		return err
	}

	for _, m := range allMeta {
		e := m.GenerateThumbnail(0)
		log.Printf("Generating new thumbnail for %s", m.Name)
		if e != nil {
			log.Printf("Failed to generate thumbnail for %s", m.Name)
			continue
		}

		provider.Write(m)
	}

	return nil
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
