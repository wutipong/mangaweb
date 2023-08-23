package view

import (
	"fmt"
	"hash/fnv"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/log"
	"github.com/wutipong/mangaweb/meta"
	"github.com/wutipong/mangaweb/tag"
)

const (
	maxPageWidth  = 1600
	maxPageHeight = 1600
)

func init() {
	var err error
	viewTemplate, err = template.New("view.gohtml").
		Funcs(handler.HtmlTemplateFuncMap()).
		ParseFiles("template/view.gohtml")
	if err != nil {
		log.Get().Sugar().Panic(err)
		os.Exit(-1)
	}
}

var viewTemplate *template.Template

type viewData struct {
	Name             string
	Title            string
	Version          string
	BrowseURL        string
	Favorite         bool
	ImageURLs        []string
	UpdateCoverURLs  []string
	DownloadPageURLs []string
	Tags             []tagData
	DownloadURL      string
	SetFavoriteURL   string
}

type tagData struct {
	Name string
	URL  string
}

func Handler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	item := handler.ParseParam(params, "item")
	item = filepath.FromSlash(item)

	query := r.URL.Query()

	m, err := meta.Read(r.Context(), item)
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
			meta.Write(r.Context(), m)
		}
	}

	if !m.IsRead {
		m.IsRead = true
		meta.Write(r.Context(), m)
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

	tags := make([]tagData, 0)

	for _, tagStr := range m.Tags {
		t, err := tag.Read(r.Context(), tagStr)
		if err != nil {
			log.Get().Sugar().Fatal(err)
			handler.WriteError(w, err)
			return
		}

		if !t.Hidden {
			tags = append(tags, tagData{
				Name: t.Name,
				URL:  handler.CreateBrowseTagURL(t.Name),
			})
		}
	}

	log.Get().Info("View Item", zap.String("item_name", item))

	data := viewData{
		Name:             item,
		Title:            fmt.Sprintf("Manga - Viewing [%s]", item),
		Version:          handler.CreateVersionString(),
		BrowseURL:        browseUrl,
		ImageURLs:        createImageURLs(item, pages),
		UpdateCoverURLs:  createUpdateCoverURLs(item, pages),
		DownloadPageURLs: createDownloadImageURLs(item, pages),
		Favorite:         m.Favorite,
		Tags:             tags,
		DownloadURL:      handler.CreateDownloadURL(item),
		SetFavoriteURL:   handler.CreateSetFavoriteURL(item),
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
