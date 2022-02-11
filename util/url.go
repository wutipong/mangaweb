package util

import (
	"path"

	html_template "html/template"
	text_template "text/template"
)

var prefix string

func SetPrefix(p string) {
	prefix = p
}

func CreateURL(p ...string) string {
	args := append([]string{prefix}, p...)

	return path.Join(args...)
}

func HtmlTemplateFuncMap() html_template.FuncMap {
	return html_template.FuncMap{
		"CreateURL": CreateURL,
	}
}

func TextTemplateFuncMap() text_template.FuncMap {
	return text_template.FuncMap{
		"CreateURL": CreateURL,
	}
}
