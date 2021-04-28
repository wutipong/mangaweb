package main

import (
	"context"
	"mangaweb/meta"
	"mangaweb/meta/mongo"
	"mangaweb/meta/postgres"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"flag"
	"log"
)

const mongoStr = "mongodb://root:password@mongo"

func setupFlag(flagName, defValue, variable, description string) *string {
	varValue := os.Getenv(variable)
	if varValue != "" {
		defValue = varValue
	}

	return flag.String(flagName, defValue, description)
}

func main() {
	address := setupFlag("address", ":80", "MANGAWEB_ADDRESS", "The server address")
	path := setupFlag("path", "/data", "MANGAWEB_IMAGE_PATH", "Image source path")
	prefix := setupFlag("prefix", "*", "MANGAWEB_URL_PREFIX", "Url prefix")
	database := setupFlag("database", "localhost:5432", "MANGAWEB_DB", "Specify the database connection string")
	useMongoStr := setupFlag("mongo", "", "MANGAWEB_USE_MONGO", "use mongo flag")

	flag.Parse()

	useMongo = *useMongoStr == "true"

	meta.BaseDirectory = *path

	log.Printf("Image Source Path: %s", *path)
	log.Printf("using prefix %s", *prefix)

	err := postgres.Init(*database)
	if err != nil {
		log.Fatal(err)
	}

	err = mongo.Init(mongoStr)
	if err != nil {
		log.Fatal(err)
	}

	err = syncToMongo()
	if err != nil {
		log.Fatal(err)
	}

	// migrateMeta(dbx)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(middleware.Rewrite(map[string]string{
		*prefix: "$1",
	}))

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

	go func() {
		if err := e.Start(*address); err != nil {
			e.Logger.Info("shutting down the server")
		}
		e.Logger.Fatal()
	}()

	metaStop, metaDone := updateMetaRoutine()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	metaStop()
	<-metaDone

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}

// Handler
func root(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "/browse")
}

func synchronizeMetaData() error {
	provider, err := postgres.New()

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

func updateMetaRoutine() (stop func(), done chan bool) {
	updateInterval := 30 * time.Minute
	isRunning := true

	stop = func() {
		isRunning = false
	}

	done = make(chan bool)

	go func() {
		for isRunning {
			log.Printf("Update metadata set.")
			synchronizeMetaData()
			<-time.After(updateInterval)
		}
		done <- true
	}()

	return
}

func syncToMongo() error {
	p1, err := postgres.New()
	if err != nil {
		return err
	}

	p2, err := mongo.New()
	if err != nil {
		return err
	}

	if c, _ := p2.NeedSetup(); c {
		return nil
	}

	items, err := p1.ReadAll()
	if err != nil {
		return err
	}

	for _, i := range items {
		err = p2.Write(i)
		if err != nil {
			return err
		}
	}

	return nil
}

var useMongo bool

func getProvider() (p meta.Provider, err error) {
	if useMongo {
		mp, e := mongo.New()
		p = &mp
		err = e
		return
	} else {
		mp, e := postgres.New()
		p = &mp
		err = e
		return
	}
}
