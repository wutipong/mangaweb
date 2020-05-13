package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"flag"
	"log"
)

func main() {
	address := flag.String("address", ":80", "The server address")
	path := flag.String("path", "/data", "Image source path")
	prefix := flag.String("prefix", "*", "Url prefix")

	flag.Parse()

	BaseDirectory = *path

	log.Printf("Image Source Path: %s", *path)
	log.Printf("using prefix %s", *prefix)

	dbx, err := initDatabase()
	defer dbx.Close()

	if err != nil {
		log.Fatal(err)
	}
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(middleware.Rewrite(map[string]string{
		*prefix: "$1",
	}))

	e.Pre(addDB(dbx))

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

	metaStop, metaDone := updateMetaRoutine(dbx)

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

func addDB(db *sqlx.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	}
}

// Handler
func root(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "/browse")
}

func updateMetaRoutine(db *sqlx.DB) (stop func(), done chan bool) {
	updateInterval := 30 * time.Minute
	checkInterval := 30 * time.Second
	lastUpdate := time.Now()

	isRunning := true

	stop = func() {
		isRunning = false
	}

	done = make(chan bool)

	go func() {
		for isRunning {
			<-time.After(checkInterval)

			now := time.Now()
			sub := now.Sub(lastUpdate)

			if sub < updateInterval {
				continue
			}

			lastUpdate = now

			files, err := ListDir()
			if err != nil {
				continue
			}

			for _, file := range files {
				if isMetaFileExist(db, file) {
					continue
				}

				NewMeta(db, file)

			}
		}
		done <- true
	}()

	return
}
