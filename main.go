package main

import (
	"flag"
	"fmt"

	"net/http"
	"os"
	"path"
	"time"

	"github.com/wutipong/mangaweb/meta"
	"github.com/wutipong/mangaweb/meta/mongo"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/gommon/log"

	urlutil "github.com/wutipong/go-utils/url"

	"github.com/go-co-op/gocron"

	"github.com/joho/godotenv"
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

		urlutil.SetPrefix(*prefix)
	}

	e.Pre(middleware.RemoveTrailingSlash())

	// Routes
	e.GET("/", root)
	e.GET("/browse", browse)
	e.GET("/browse/*", browse)

	e.GET("/view", view)
	e.GET("/view/*", view)

	e.Static("/static", "static")

	e.GET("/get_image/*", GetImage)

	e.GET("/thumbnail/*", thumbnail)

	e.GET("/view", view)
	e.GET("/view/*", view)

	// Schedule the update metadata task to run every 30 minutes.
	s := gocron.NewScheduler(time.UTC)
	s.Every(30).Minutes().Do(func() {
		log.Info("Update metadata set.")
		synchronizeMetaData()
	})
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
	return c.Redirect(http.StatusPermanentRedirect, urlutil.CreateURL("/browse"))
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
		_, err := provider.New(file)
		if err != nil {
			log.Printf("Failed to create meta data : %v", err)
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
