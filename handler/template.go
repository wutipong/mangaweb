package handler

import (
	"html/template"

	"github.com/wutipong/mangaweb/util"
)

func HtmlTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"CreateURL": util.CreateURL,
	}
}
