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
	ItemPerPage = 40
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
	Pages        []pageItem
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

type pageItem struct {
	Content   string
	LinkURL   url.URL
	IsActive  bool
	IsEnabled bool
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

	pageCount := int(count / ItemPerPage)
	if count%ItemPerPage > 0 {
		pageCount++
	}

	if page > pageCount || page < 0 {
		page = 0
	}

	data := browseData{
		Title:        "Manga - Browsing",
		Version:      versionString,
		FavoriteOnly: favOnly,
		Items:        items,
		Pages:        createPageItems(page, pageCount, *c.Request().URL),
	}

	builder := strings.Builder{}
	err = broseTemplate.Execute(&builder, data)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return c.HTML(http.StatusOK, builder.String())
}

func createPageItems(current int, count int, baseUrl url.URL) []pageItem {
	const (
		First    = "First"
		Previous = "Previous"
		Next     = "Next"
		Last     = "Last"

		DisplayPageCount     = 6
		HalfDisplayPageCount = DisplayPageCount / 2
	)

	firstPage := 0
	lastPage := count - 1
	previousPage := current - 1
	nextPage := current + 1

	changePageParam := func(baseUrl url.URL, page int) url.URL {
		query := baseUrl.Query()

		if query.Has("page") {
			query.Set("page", strconv.Itoa(page))
		} else {
			query.Add("page", strconv.Itoa(page))
		}

		baseUrl.RawQuery = query.Encode()
		return baseUrl
	}

	output := make([]pageItem, 0)
	output = append(output, pageItem{
		Content:   First,
		LinkURL:   changePageParam(baseUrl, firstPage),
		IsActive:  false,
		IsEnabled: true,
	})

	enablePrevious := previousPage >= firstPage
	output = append(output, pageItem{
		Content:   Previous,
		LinkURL:   changePageParam(baseUrl, previousPage),
		IsActive:  false,
		IsEnabled: enablePrevious,
	})

	for i := current - HalfDisplayPageCount; i <= current+HalfDisplayPageCount; i++ {
		if i < firstPage {
			continue
		}
		if i > lastPage {
			continue
		}

		output = append(output, pageItem{
			Content:   strconv.Itoa(i),
			LinkURL:   changePageParam(baseUrl, i),
			IsActive:  i == current,
			IsEnabled: true,
		})
	}

	enableNext := nextPage < count
	output = append(output, pageItem{
		Content:   Next,
		LinkURL:   changePageParam(baseUrl, nextPage),
		IsActive:  false,
		IsEnabled: enableNext,
	})

	output = append(output, pageItem{
		Content:   Last,
		LinkURL:   changePageParam(baseUrl, lastPage),
		IsActive:  false,
		IsEnabled: true,
	})

	return output
}
