package handler

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
)

func CreateFilePathURL(p string) string {
	parts := strings.Split(p, string(os.PathSeparator))
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}

	return path.Join(parts...)
}

func CreateURL(p ...string) string {
	args := append([]string{options.PathPrefix}, p...)

	return path.Join(args...)
}

func CreateViewURL(filepath string) string {
	urlStr := CreateFilePathURL(filepath)
	urlStr = CreateURL(options.PathView, urlStr)

	return urlStr
}

func CreateThumbnailURL(filepath string) string {
	urlStr := CreateFilePathURL(filepath)
	urlStr = CreateURL(options.PathThumbnail, urlStr)

	return urlStr
}

func CreateRescanURL() string {
	return CreateURL(options.PathRescanLibrary)
}

func CreateGetImageURL(filepath string, index int) string {
	filePart := CreateFilePathURL(filepath)
	url := CreateURL(options.PathGetImage, fmt.Sprintf("%s?i=%v", filePart, index))

	return url
}

func CreateGetImageWithSizeURL(filepath string, index int, width int, height int) string {
	filePart := CreateFilePathURL(filepath)
	url := CreateURL(options.PathGetImage, fmt.Sprintf("%s?i=%v&width=%v&height=%v", filePart, index, width, height))

	return url
}

func CreateUpdateCoverURL(filepath string, index int) string {
	filePart := CreateFilePathURL(filepath)
	baseURL := path.Join(options.PathUpdateCover, filePart)
	url := CreateURL(fmt.Sprintf("%s?i=%v", baseURL, index))

	return url
}

func CreateDownloadURL(filepath string) string {
	urlStr := CreateFilePathURL(filepath)
	urlStr = CreateURL(options.PathDownload, urlStr)

	return urlStr
}

func CreateSetFavoriteURL(filepath string) string {
	urlStr := CreateFilePathURL(filepath)
	urlStr = CreateURL(options.PathFavorite, urlStr)

	return urlStr
}

func CreateBrowseURL(id string) string {
	return CreateURL(fmt.Sprintf("%s#%v", options.PathBrowse, id))
}

func CreateBrowseTagURL(tag string) string {
	tagUrl := CreateFilePathURL(tag)
	return CreateURL(options.PathBrowse, tagUrl)
}

func CreateSetTagFavoriteURL(filepath string) string {
	urlStr := CreateFilePathURL(filepath)
	urlStr = CreateURL(options.PathTagFavorite, urlStr)

	return urlStr
}

func CreateTagListURL() string {
	return CreateURL(options.PathTagList)
}

func CreateTagThumbnailURL(filepath string) string {
	urlStr := CreateFilePathURL(filepath)
	urlStr = CreateURL(options.PathTagThumbnail, urlStr)

	return urlStr
}

func CreateRootURL() string {
	return CreateURL(options.PathRoot)
}
