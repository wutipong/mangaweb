package browse

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"hash/fnv"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/log"
	"github.com/wutipong/mangaweb/meta"
)

const (
	ItemPerPage = 40
)

func init() {
	var err error
	broseTemplate, err = template.New("browse.gohtml").
		Funcs(handler.HtmlTemplateFuncMap()).
		ParseFiles("template/browse.gohtml")
	if err != nil {
		log.Get().Sugar().Panic(err)
		os.Exit(-1)
	}
}

var broseTemplate *template.Template

type browseData struct {
	Title        string
	Version      string
	FavoriteOnly bool
	SortBy       string
	SortOrder    string
	Tag          string
	TagFavorite  bool
	BrowseURL    string
	TagListURL   string

	Items []item
	Pages []pageItem
}

type item struct {
	ID         uint64
	Name       string
	CreateTime time.Time
	Favorite   bool
	IsRead     bool
}

type pageItem struct {
	Content         string
	LinkURL         url.URL
	IsActive        bool
	IsEnabled       bool
	IsHiddenOnSmall bool
}

func createItems(allMeta []meta.Meta) (allItems []item, err error) {
	allItems = make([]item, len(allMeta))

	for i, m := range allMeta {
		hash := fnv.New64()
		hash.Write([]byte(m.Name))
		id := hash.Sum64()

		allItems[i] = item{
			ID:         id,
			Name:       m.Name,
			CreateTime: m.CreateTime,
			Favorite:   m.Favorite,
			IsRead:     m.IsRead,
		}
	}
	return
}

func Handler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	query := r.URL.Query()
	tagStr := handler.ParseParam(params, "tag")

	p, err := handler.CreateMetaProvider()
	if err != nil {
		handler.WriteError(w, err)
		return
	}
	defer p.Close()

	query.Get("favorite")

	favOnly := false
	if f, e := strconv.ParseBool(query.Get("favorite")); e == nil {
		favOnly = f
	}

	page := 0
	if i, e := strconv.ParseInt(query.Get("page"), 10, 0); e == nil {
		page = int(i)
	}

	search := query.Get("search")
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

	if tagStr != "" {
		searchCriteria = append(searchCriteria, meta.SearchCriteria{
			Field: meta.SearchFieldTag,
			Value: tagStr,
		})
	}

	sort := parseSortField(query.Get("sort"))
	order := parseSortOrder(query.Get("order"))

	allMeta, err := p.Search(searchCriteria, sort, order, ItemPerPage, page)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	items, err := createItems(allMeta)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	count, err := p.Count(searchCriteria)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	pageCount := int(count / ItemPerPage)
	if count%ItemPerPage > 0 {
		pageCount++
	}

	if page > pageCount || page < 0 {
		page = 0
	}

	log.Get().Info("Browse",
		zap.Int("page", page),
		zap.String("tag", tagStr),
		zap.Bool("favorite_only", favOnly),
		zap.String("sort_by", string(sort)),
		zap.String("sort_order", string(order)))

	data := browseData{
		Title:        "Browse - All items",
		Version:      handler.CreateVersionString(),
		FavoriteOnly: favOnly,
		SortBy:       string(sort),
		SortOrder:    string(order),
		Items:        items,
		Pages:        createPageItems(page, pageCount, *r.URL),
		BrowseURL:    handler.CreateBrowseURL(""),
		TagListURL:   handler.CreateTagListURL(),
	}

	if tagStr != "" {
		data.Title = fmt.Sprintf("Browse - %s", tagStr)
		data.Tag = tagStr

		tagProvider, err := handler.CreateTagProvider()
		if err != nil {
			handler.WriteError(w, err)
			return
		}

		tagObj, err := tagProvider.Read(tagStr)
		if err != nil {
			handler.WriteError(w, err)
			return
		}

		data.TagFavorite = tagObj.Favorite
	}

	builder := strings.Builder{}
	err = broseTemplate.Execute(&builder, data)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	handler.WriteHtml(w, builder.String())
}

func parseSortOrder(orderStr string) meta.SortOrder {
	order := meta.SortOrder(orderStr)

	switch order {
	case meta.SortOrderAscending:
		return order
	case meta.SortOrderDescending:
		return order
	}

	return meta.SortOrderDescending
}

func parseSortField(sortBy string) meta.SortField {
	sort := meta.SortField(sortBy)

	switch sort {
	case meta.SortFieldName:
		return sort
	case meta.SortFieldCreateTime:
		return sort

	default:
		return meta.SortFieldCreateTime
	}
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
		Content:         First,
		LinkURL:         changePageParam(baseUrl, firstPage),
		IsActive:        false,
		IsEnabled:       true,
		IsHiddenOnSmall: false,
	})

	enablePrevious := previousPage >= firstPage
	output = append(output, pageItem{
		Content:         Previous,
		LinkURL:         changePageParam(baseUrl, previousPage),
		IsActive:        false,
		IsEnabled:       enablePrevious,
		IsHiddenOnSmall: false,
	})

	for i := current - HalfDisplayPageCount; i <= current+HalfDisplayPageCount; i++ {
		if i < firstPage {
			continue
		}
		if i > lastPage {
			continue
		}

		output = append(output, pageItem{
			Content:         strconv.Itoa(i),
			LinkURL:         changePageParam(baseUrl, i),
			IsActive:        i == current,
			IsEnabled:       true,
			IsHiddenOnSmall: !(i == current),
		})
	}

	enableNext := nextPage < count
	output = append(output, pageItem{
		Content:         Next,
		LinkURL:         changePageParam(baseUrl, nextPage),
		IsActive:        false,
		IsEnabled:       enableNext,
		IsHiddenOnSmall: false,
	})

	output = append(output, pageItem{
		Content:         Last,
		LinkURL:         changePageParam(baseUrl, lastPage),
		IsActive:        false,
		IsEnabled:       true,
		IsHiddenOnSmall: false,
	})

	return output
}
