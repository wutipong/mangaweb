package handler

import (
	"fmt"
	"github.com/wutipong/mangaweb/util"
	"path"
)

func CreateViewURL(filepath string) string {
	urlStr := util.CreateFilePathURL(filepath)
	urlStr = path.Join(options.PathView, urlStr)
	urlStr = util.CreateURL(urlStr)

	return urlStr
}

func CreateThumbnailURL(filepath string) string {
	urlStr := util.CreateFilePathURL(filepath)
	urlStr = path.Join(options.PathThumbnail, urlStr)
	urlStr = util.CreateURL(urlStr)

	return urlStr
}

func CreateRescanURL() string {
	return util.CreateURL(options.PathRescanLibrary)
}

func CreateGetImageURL(filepath string, index int) string {
	filePart := util.CreateFilePathURL(filepath)
	baseURL := path.Join(options.PathGetImage, filePart)
	url := util.CreateURL(fmt.Sprintf("%s?i=%v", baseURL, index))

	return url
}

func CreateUpdateCoverURL(filepath string, index int) string {
	filePart := util.CreateFilePathURL(filepath)
	baseURL := path.Join(options.PathUpdateCover, filePart)
	url := util.CreateURL(fmt.Sprintf("%s?i=%v", baseURL, index))

	return url
}

func CreateDownloadURL(filepath string) string {
	urlStr := util.CreateFilePathURL(filepath)
	urlStr = path.Join(options.PathDownload, urlStr)
	urlStr = util.CreateURL(urlStr)

	return urlStr
}

func CreateSetFavoriteURL(filepath string) string {
	urlStr := util.CreateFilePathURL(filepath)
	urlStr = path.Join(options.PathFavorite, urlStr)
	urlStr = util.CreateURL(urlStr)

	return urlStr
}

func CreateBrowseURL(id string) string {
	return util.CreateURL(fmt.Sprintf("%s#%v", options.PathBrowse, id))
}
