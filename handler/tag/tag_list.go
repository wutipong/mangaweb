package tag

import (
	"hash/fnv"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/log"

	"github.com/wutipong/mangaweb/tag"
)

const (
	ItemPerPage = 40
)

func init() {
	var err error
	templateObj, err = template.New("tag_list.gohtml").
		Funcs(handler.HtmlTemplateFuncMap()).
		ParseFiles("template/tag_list.gohtml")
	if err != nil {
		log.Get().Sugar().Panic(err)
		os.Exit(-1)
	}
}

var templateObj *template.Template

type PageData struct {
	Title      string
	Version    string
	BrowseURL  string
	TagListURL string
	Tags       []ItemData
}

type ItemData struct {
	ID           uint64
	Name         string
	Favorite     bool
	URL          string
	ThumbnailURL string
}

func createItems(allTags []tag.Tag, favoriteOnly bool) []ItemData {
	allItems := make([]ItemData, len(allTags))

	for i, t := range allTags {
		isAdding := true
		if favoriteOnly {
			isAdding = t.Favorite
		}

		if isAdding {
			hash := fnv.New64()
			hash.Write([]byte(t.Name))
			id := hash.Sum64()

			allItems[i] = ItemData{
				ID:           id,
				Name:         t.Name,
				Favorite:     t.Favorite,
				URL:          handler.CreateBrowseTagURL(t.Name),
				ThumbnailURL: handler.CreateTagThumbnailURL(t.Name),
			}
		}
	}

	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].Name < allItems[j].Name
	})
	return allItems
}

func TagListHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	query := r.URL.Query()

	log.Get().Info("Tag list")

	favOnly := false
	if f, e := strconv.ParseBool(query.Get("favorite")); e == nil {
		favOnly = f
	}

	allTags, err := tag.ReadAll()
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	tagData := createItems(allTags, favOnly)

	data := PageData{
		Title:      "Tag list",
		Version:    handler.CreateVersionString(),
		Tags:       tagData,
		TagListURL: handler.CreateTagListURL(),
		BrowseURL:  handler.CreateBrowseURL(""),
	}

	builder := strings.Builder{}
	err = templateObj.Execute(&builder, data)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	handler.WriteHtml(w, builder.String())
}
