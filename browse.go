package main

import (
	"fmt"
	"hash/fnv"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
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
	Rows  [][]item
}

type item struct {
	ID       uint64
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

		hash := fnv.New64()
		hash.Write([]byte(f))
		id := hash.Sum64()

		output[i] = item{
			ID:       id,
			Name:     f,
			LinkURL:  url,
			ThumbURL: thumbURL,
		}
	}
	return output
}

func makeRows(items []item, col int) [][]item {
	var rows [][]item

	for i, it := range items {
		if i%col == 0 {
			rows = append(rows, make([]item, 0))
		}

		r := i / col
		rows[r] = append(rows[r], it)
	}

	return rows
}

// Handler
func browse(c echo.Context) error {
	builder := strings.Builder{}

	files, err := ListDir()
	if err != nil {
		return err
	}

	sort.Strings(files)
	items := createItems(files)
	data := browseData{
		Title: fmt.Sprintf("Manga - Browsing"),
		Rows:  makeRows(items, 2),
	}
	err = broseTemplate.Execute(&builder, data)
	if err != nil {
		log.Println(err)
		return err
	}

	return c.HTML(http.StatusOK, builder.String())
}
