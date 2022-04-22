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

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/util"
)

func init() {
	var err error
	viewTemplate, err = template.New("view.gohtml").
		Funcs(util.HtmlTemplateFuncMap()).
		ParseFiles(
			"template/view.gohtml",
			"template/header.gohtml",
		)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
}

var viewTemplate *template.Template

type viewData struct {
	Name            string
	Title           string
	BrowseURL       string
	ImageURLs       []string
	UpdateCoverURLs []string
	StartIndex      int64
	Favorite        bool
	DownloadURL     string
	SetFavoriteURL  string
}

func Handler(c echo.Context) error {
	fileName := c.Param("*")
	fileName = filepath.FromSlash(fileName)
	db, err := handler.CreateMetaProvider()
	if err != nil {
		return err
	}
	defer db.Close()

	m, err := db.Read(fileName)
	if err != nil {
		return err
	}

	pages, err := ListPages(m)
	if err != nil {
		return err
	}

	hash := fnv.New64()
	hash.Write([]byte(fileName))
	id := hash.Sum64()

	if fav, e := strconv.ParseBool(c.QueryParam("favorite")); e == nil {
		if fav != m.Favorite {
			m.Favorite = fav
			db.Write(m)
		}
	}

	if !m.IsRead {
		m.IsRead = true
		db.Write(m)
	}

	browseUrl := c.Request().Referer()
	if browseUrl == "" {
		util.CreateURL(fmt.Sprintf("/browse#%v", id))
	} else {
		if u, e := url.Parse(browseUrl); e == nil {
			u.Fragment = strconv.FormatUint(id, 10)
			browseUrl = u.String()
		}
	}

	data := viewData{
		Name:            fileName,
		Title:           fmt.Sprintf("Manga - Viewing [%s]", fileName),
		BrowseURL:       browseUrl,
		ImageURLs:       createImageURLs(fileName, pages),
		UpdateCoverURLs: createUpdateCoverURLs(fileName, pages),
		Favorite:        m.Favorite,
		DownloadURL:     createDownloadURL(fileName),
		SetFavoriteURL:  createSetFavoriteURL(fileName),
	}

	builder := strings.Builder{}
	err = viewTemplate.Execute(&builder, data)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return c.HTML(http.StatusOK, builder.String())
}

func createImageURLs(file string, pages []Page) []string {
	output := make([]string, len(pages))
	for i, p := range pages {
		filePart := util.CreateFilePathURL(file)
		url := util.CreateURL(fmt.Sprintf("/get_image/%s?i=%v", filePart, p.Index))

		output[i] = url
	}
	return output
}

func createUpdateCoverURLs(file string, pages []Page) []string {
	output := make([]string, len(pages))
	for i, p := range pages {
		filePart := util.CreateFilePathURL(file)
		url := util.CreateURL(fmt.Sprintf("/update_cover/%s?i=%v", filePart, p.Index))

		output[i] = url
	}
	return output
}

func createDownloadURL(file string) string {
	filePart := util.CreateFilePathURL(file)
	return util.CreateURL(fmt.Sprintf("/download/%s", filePart))
}

func createSetFavoriteURL(file string) string {
	filePart := util.CreateFilePathURL(file)
	return util.CreateURL(fmt.Sprintf("/favorite/%s", filePart))
}
