package util

import (
	"net/url"
	"os"
	"path"
	"strings"

	html_template "html/template"
)

var prefix string

func SetPrefix(p string) {
	prefix = p
}

func CreateFilePathURL(p string) string {
	parts := strings.Split(p, string(os.PathSeparator))
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}

	return path.Join(parts...)
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
