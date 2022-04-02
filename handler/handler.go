package handler

import "github.com/wutipong/mangaweb/meta"

var CreateMetaProvider meta.MetaProviderFactory
var VersionString string

type Options struct {
	MetaProviderFactory meta.MetaProviderFactory
	VersionString       string
}

func Init(options Options) {
	CreateMetaProvider = options.MetaProviderFactory
	VersionString = options.VersionString
}
