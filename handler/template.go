package handler

import (
	"encoding/json"
	"html/template"
)

func marshal(i any) template.JS {
	str, _ := json.Marshal(i)

	return template.JS(str)
}

func HtmlTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"CreateURL":               CreateURL,
		"CreateViewURL":           CreateViewURL,
		"CreateThumbnailURL":      CreateThumbnailURL,
		"CreateRescanURL":         CreateRescanURL,
		"CreateGetImageURL":       CreateGetImageURL,
		"CreateUpdateCoverURL":    CreateUpdateCoverURL,
		"CreateDownloadURL":       CreateDownloadURL,
		"CreateSetFavoriteURL":    CreateSetFavoriteURL,
		"CreateBrowseURL":         CreateBrowseURL,
		"CreateSetTagFavoriteURL": CreateSetTagFavoriteURL,
		"CreateTagListURL":        CreateTagListURL,
		"CreateBrowseTagURL":      CreateBrowseTagURL,
		"CreateTagThumbnailURL":   CreateTagThumbnailURL,
		"CreateRootURL":           CreateRootURL,
		"Marshal":                 marshal,
	}
}
