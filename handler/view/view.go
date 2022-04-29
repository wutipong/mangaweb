package view

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"hash/fnv"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/log"
)

const (
	maxPageWidth  = 1600
	maxPageHeight = 1600
)

func init() {
	var err error
	viewTemplate, err = template.New("view.gohtml").
		Funcs(handler.HtmlTemplateFuncMap()).
		ParseFiles(
			"template/view.gohtml",
			"template/header.gohtml",
		)
	if err != nil {
		log.Get().Sugar().Panic(err)
		os.Exit(-1)
	}
}

var viewTemplate *template.Template

type viewData struct {
	Name             string
	Title            string
	BrowseURL        string
	Favorite         bool
	ImageURLs        []string
	UpdateCoverURLs  []string
	DownloadPageURLs []string
	Tags             []string
}

func Handler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	item := handler.ParseParam(params, "item")
	item = filepath.FromSlash(item)

	query := r.URL.Query()

	db, err := handler.CreateMetaProvider()
	if err != nil {
		handler.WriteError(w, err)
		return
	}
	defer db.Close()

	m, err := db.Read(item)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	pages, err := ListPages(m)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	hash := fnv.New64()
	hash.Write([]byte(item))
	id := hash.Sum64()

	if fav, e := strconv.ParseBool(query.Get("favorite")); e == nil {
		if fav != m.Favorite {
			m.Favorite = fav
			db.Write(m)
		}
	}

	if !m.IsRead {
		m.IsRead = true
		db.Write(m)
	}

	browseUrl := r.Referer()
	if browseUrl == "" {
		browseUrl = handler.CreateBrowseURL(strconv.FormatUint(id, 16))
	} else {
		if u, e := url.Parse(browseUrl); e == nil {
			u.Fragment = strconv.FormatUint(id, 10)
			browseUrl = u.String()
		}
	}
	tagProvider, err := handler.CreateTagProvider()
	if err != nil {
		log.Get().Sugar().Fatal(err)
		handler.WriteError(w, err)
		return
	}

	tags := make([]string, 0)
	for _, tagStr := range m.Tags {
		t, err := tagProvider.Read(tagStr)
		if err != nil {
			log.Get().Sugar().Fatal(err)
			handler.WriteError(w, err)
			return
		}

		if !t.Hidden {
			tags = append(tags, t.Name)
		}
	}

	log.Get().Info("View Item", zap.String("item_name", item))

	data := viewData{
		Name:             item,
		Title:            fmt.Sprintf("Manga - Viewing [%s]", item),
		BrowseURL:        browseUrl,
		ImageURLs:        createImageURLs(item, pages),
		UpdateCoverURLs:  createUpdateCoverURLs(item, pages),
		DownloadPageURLs: createDownloadImageURLs(item, pages),
		Favorite:         m.Favorite,
		Tags:             tags,
	}

	builder := strings.Builder{}
	err = viewTemplate.Execute(&builder, data)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	handler.WriteHtml(w, builder.String())
}

func createDownloadImageURLs(file string, pages []Page) []string {
	output := make([]string, len(pages))
	for i, p := range pages {
		output[i] = handler.CreateGetImageURL(file, p.Index)
	}
	return output
}

func createImageURLs(file string, pages []Page) []string {
	output := make([]string, len(pages))
	for i, p := range pages {
		output[i] = handler.CreateGetImageWithSizeURL(file, p.Index, maxPageWidth, maxPageHeight)
	}
	return output
}

func createUpdateCoverURLs(file string, pages []Page) []string {
	output := make([]string, len(pages))
	for i, p := range pages {
		output[i] = handler.CreateUpdateCoverURL(file, p.Index)
	}
	return output
}
