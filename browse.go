package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func init() {
	var err error
	broseTemplate, err = template.New("browse.gohtml").
		ParseFiles(
			"template/browse.gohtml",
			"template/header.gohtml",
		)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
}

var broseTemplate *template.Template = nil

type browseData struct {
	Title string
	Items []item
}

type item struct {
	Name     string
	LinkURL  string
	ThumbURL string
}

func createItems(files []string) []item {
	output := make([]item, len(files))
	for i, f := range files {
		var url string
		var thumbURL string

		url = "/view/" + f
		thumbURL = "/get_image/" + f + "?i=0;width=200"

		output[i] = item{
			Name:     f,
			LinkURL:  url,
			ThumbURL: thumbURL,
		}
	}
	return output
}

// Handler
func browse(c echo.Context) error {
	builder := strings.Builder{}

	files, err := ListDir()
	if err != nil {
		return err
	}
	data := browseData{
		Title: fmt.Sprintf("Manga - Browsing"),
		Items: createItems(files),
	}
	err = broseTemplate.Execute(&builder, data)
	if err != nil {
		log.Println(err)
		return err
	}

	return c.HTML(http.StatusOK, builder.String())
}
