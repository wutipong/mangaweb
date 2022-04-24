package handler

import (
	"html/template"
)

func HtmlTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"CreateURL":            CreateURL,
		"CreateViewURL":        CreateViewURL,
		"CreateThumbnailURL":   CreateThumbnailURL,
		"CreateRescanURL":      CreateRescanURL,
		"CreateGetImageURL":    CreateGetImageURL,
		"CreateUpdateCoverURL": CreateUpdateCoverURL,
		"CreateDownloadURL":    CreateDownloadURL,
		"CreateSetFavoriteURL": CreateSetFavoriteURL,
		"CreateBrowseURL":      CreateBrowseURL,
	}
}