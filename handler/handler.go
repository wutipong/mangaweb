package handler

import (
	"github.com/wutipong/mangaweb/meta"
	"github.com/wutipong/mangaweb/tag"
)

var options Options

type Options struct {
	MetaProviderFactory meta.ProviderFactory
	TagProviderFactory  tag.ProviderFactory
	VersionString       string

	PathPrefix        string
	PathRoot          string
	PathBrowse        string
	PathView          string
	PathStatic        string
	PathGetImage      string
	PathUpdateCover   string
	PathThumbnail     string
	PathFavorite      string
	PathDownload      string
	PathRescanLibrary string
}

func Init(o Options) {
	options = o
}

func CreateMetaProvider() (provider meta.Provider, err error) {
	return options.MetaProviderFactory()
}

func CreateVersionString() string {
	return options.VersionString
}
