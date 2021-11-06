package main

import (
	"hash/fnv"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/wutipong/mangaweb/meta"

	"github.com/labstack/echo/v4"

	urlutil "github.com/wutipong/go-utils/url"
)

func init() {
	var err error
	broseTemplate, err = template.New("browse.gohtml").
		Funcs(urlutil.HtmlTemplateFuncMap()).
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
	Items        []item
}

type item struct {
	ID         uint64
	Name       string
	LinkURL    string
	ThumbURL   string
	CreateTime time.Time
	Favorite   bool
	IsRead     bool
}

func createItems(allMeta []meta.Item) (allItems []item, err error) {
	allItems = make([]item, len(allMeta))

	for i, m := range allMeta {
		var urlStr string
		var thumbURL string

		urlStr = urlutil.CreateURL("/view/", url.PathEscape(m.Name))
		thumbURL = urlutil.CreateURL("/thumbnail/", url.PathEscape(m.Name))

		hash := fnv.New64()
		hash.Write([]byte(m.Name))
		id := hash.Sum64()

		allItems[i] = item{
			ID:         id,
			Name:       m.Name,
			LinkURL:    urlStr,
			ThumbURL:   thumbURL,
			CreateTime: m.CreateTime,
			Favorite:   m.Favorite,
			IsRead:     m.IsRead,
		}
	}
	return
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
	p, err := newProvider()
	if err != nil {
		return err
	}
	defer p.Close()

	builder := strings.Builder{}

	favOnly := false
	if f, e := strconv.ParseBool(c.QueryParam("favorite")); e == nil {
		favOnly = f
	}

	sortBy := c.QueryParam("sort")
	if sortBy == "" {
		sortBy = "date"
	}

	descending := false
	if f, e := strconv.ParseBool(c.QueryParam("descending")); e == nil {
		descending = f
	}

	var allMeta []meta.Item
	search := c.QueryParam("search")
	if search == "" {
		allMeta, err = p.ReadAll()
		if err != nil {
			return err
		}
	} else {
		allMeta, err = p.Find(search)
		if err != nil {
			return err
		}
	}

	items, err := createItems(allMeta)
	if err != nil {
		return err
	}

	if favOnly {
		var tempItems []item
		for _, item := range items {
			if item.Favorite {
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
		Title:        "Manga - Browsing",
		FavoriteOnly: favOnly,
		Items:        items,
	}
	err = broseTemplate.Execute(&builder, data)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return c.HTML(http.StatusOK, builder.String())
}
