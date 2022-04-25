package tag

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/wutipong/mangaweb/handler"
	"hash/fnv"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/wutipong/mangaweb/tag"
)

const (
	ItemPerPage = 40
)

func init() {
	var err error
	templateObj, err = template.New("tag_list.gohtml").
		Funcs(handler.HtmlTemplateFuncMap()).
		ParseFiles(
			"template/tag_list.gohtml",
			"template/header.gohtml",
		)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
}

var templateObj *template.Template

type PageData struct {
	Title string
	Tags  []ItemData
}

type ItemData struct {
	ID       uint64
	Name     string
	Favorite bool
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
				ID:       id,
				Name:     t.Name,
				Favorite: t.Favorite,
			}
		}
	}

	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].Name < allItems[j].Name
	})
	return allItems
}

func TagListHandler(c echo.Context) error {
	favOnly := false
	if f, e := strconv.ParseBool(c.QueryParam("favorite")); e == nil {
		favOnly = f
	}

	tagProvider, err := handler.CreateTagProvider()
	if err != nil {
		log.Error(err)
		return err
	}

	allTags, err := tagProvider.ReadAll()
	if err != nil {
		log.Error(err)
		return err
	}

	tagData := createItems(allTags, favOnly)

	data := PageData{
		Title: "Tag list",
		Tags:  tagData,
	}

	builder := strings.Builder{}
	err = templateObj.Execute(&builder, data)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return c.HTML(http.StatusOK, builder.String())
}
