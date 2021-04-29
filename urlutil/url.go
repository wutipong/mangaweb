package urlutil

import (
	"html/template"
	"path"
)

var prefix string

func SetPrefix(p string) {
	prefix = p
}

func CreateURL(p ...string) string {
	args := append([]string{prefix}, p...)

	return path.Join(args...)
}

func TemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"CreateURL": CreateURL,
	}
}
