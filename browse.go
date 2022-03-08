package main

import (
	"hash/fnv"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/wutipong/mangaweb/meta"
	"github.com/wutipong/mangaweb/util"
)

const (
	ItemPerPage = 10
)

func init() {
	var err error
	broseTemplate, err = template.New("browse.gohtml").
		Funcs(util.HtmlTemplateFuncMap()).
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
	Version      string
	FavoriteOnly bool
	Items        []item
	AllItemCount int64
	ItemPerPage  int
	PageIndex    int
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

		urlStr = util.CreateURL("/view/", url.PathEscape(m.Name))
		thumbURL = util.CreateURL("/thumbnail/", url.PathEscape(m.Name))

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

	page := 0
	if i, e := strconv.ParseInt(c.QueryParam("page"), 10, 0); e == nil {
		page = int(i)
	}

	search := c.QueryParam("search")
	searchCriteria := make([]meta.SearchCriteria, 0)
	if search != "" {
		searchCriteria = append(searchCriteria, meta.SearchCriteria{
			Field: meta.SearchFieldName,
			Value: search,
		})
	}

	if favOnly {
		searchCriteria = append(searchCriteria, meta.SearchCriteria{
			Field: meta.SearchFieldFavorite,
			Value: true,
		})
	}

	var sort meta.SortField
	switch sortBy {
	case "name":
		sort = meta.SortFieldName
	case "date":
		sort = meta.SortFieldCreateTime
	}

	order := meta.SortOrderAscending

	if descending {
		order = meta.SortOrderDescending
	}

	allMeta, err := p.Search(searchCriteria, sort, order, ItemPerPage, page)
	if err != nil {
		return err
	}

	items, err := createItems(allMeta)
	if err != nil {
		return err
	}

	count, err := p.Count(searchCriteria)
	if err != nil {
		return err
	}

	data := browseData{
		Title:        "Manga - Browsing",
		Version:      "0", //TODO: versionString,
		FavoriteOnly: favOnly,
		Items:        items,
		AllItemCount: count,
		ItemPerPage:  ItemPerPage,
		PageIndex:    page,
	}
	err = broseTemplate.Execute(&builder, data)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return c.HTML(http.StatusOK, builder.String())
}
