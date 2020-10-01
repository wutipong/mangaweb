package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
)

func init() {
	var err error
	viewTemplate, err = template.New("view.gohtml").
		ParseFiles(
			"template/view.gohtml",
			"template/header.gohtml",
		)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
}

var viewTemplate *template.Template

type viewData struct {
	Name       string
	Title      string
	BrowseURL  string
	ImageURLs  []string
	StartIndex int64
	Favorite   bool
}

func view(c echo.Context) error {
	builder := strings.Builder{}

	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	pages, err := ListPages(db, p)
	if err != nil {
		return err
	}

	hash := fnv.New64()
	hash.Write([]byte(p))
	id := hash.Sum64()

	var meta itemMeta
	meta.Read(db, p)
	if fav, e := strconv.ParseBool(c.QueryParam("favorite")); e == nil {
		if fav != meta.Favorite {
			meta.Favorite = fav
			meta.Write(db)
		}
	}

	data := viewData{
		Name:      p,
		Title:     fmt.Sprintf("Manga - Viewing [%s]", p),
		BrowseURL: fmt.Sprintf("/browse#%v", id),
		ImageURLs: createImageURLs(p, pages),
		Favorite:  meta.Favorite,
	}
	err = viewTemplate.Execute(&builder, data)
	if err != nil {
		log.Println(err)
		return err
	}

	return c.HTML(http.StatusOK, builder.String())
}

func createImageURLs(file string, pages []Page) []string {
	output := make([]string, len(pages))
	for i, p := range pages {
		url := fmt.Sprintf("/get_image/%s?i=%v", url.PathEscape(file), p.Index)

		output[i] = url
	}
	return output
}
