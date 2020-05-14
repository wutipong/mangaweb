package main

import (
	"fmt"
	"hash/fnv"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
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

var broseTemplate *template.Template

type browseData struct {
	Title        string
	FavoriteOnly bool
	Rows         [][]item
}

type item struct {
	ID         uint64
	Name       string
	LinkURL    string
	ThumbURL   string
	CreateTime time.Time
	Favorite   bool
}

func createItems(db *sqlx.DB, files []string) []item {
	output := make([]item, len(files))
	for i, f := range files {
		var url string
		var thumbURL string

		url = "/view/" + f
		thumbURL = "/thumbnail/" + f

		hash := fnv.New64()
		hash.Write([]byte(f))
		id := hash.Sum64()

		meta, _ := OpenMeta(db, f)
		output[i] = item{
			ID:         id,
			Name:       f,
			LinkURL:    url,
			ThumbURL:   thumbURL,
			CreateTime: meta.CreateTime,
			Favorite:   meta.Favorite,
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
	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	builder := strings.Builder{}

	files, err := ListDir()
	if err != nil {
		return err
	}

	fav := false
	if f, e := strconv.ParseBool(c.QueryParam("favorite")); e == nil {
		fav = f
	}

	sortBy := c.QueryParam("sort")
	if sortBy == "" {
		sortBy = "name"
	}

	descending := false
	if f, e := strconv.ParseBool(c.QueryParam("descending")); e == nil {
		descending = f
	}

	items := createItems(db, files)
	if fav == true {
		var tempItems []item
		for _, item := range items {
			if item.Favorite == true {
				tempItems = append(tempItems, item)
			}
		}
		items = tempItems
	}

	switch sortBy {
	case "name":
		sort.Slice(items, func(i, j int) bool {
			return items[i].Name < items[j].Name
		})
	case "date":
		sort.Slice(items, func(i, j int) bool {
			return items[j].CreateTime.Before(items[i].CreateTime)
		})
	}

	if descending {
		for i := len(items)/2 - 1; i >= 0; i-- {
			opp := len(items) - 1 - i
			items[i], items[opp] = items[opp], items[i]
		}
	}

	data := browseData{
		Title:        fmt.Sprintf("Manga - Browsing"),
		FavoriteOnly: fav,
		Rows:         makeRows(items, 2),
	}
	err = broseTemplate.Execute(&builder, data)
	if err != nil {
		log.Println(err)
		return err
	}

	return c.HTML(http.StatusOK, builder.String())
}
