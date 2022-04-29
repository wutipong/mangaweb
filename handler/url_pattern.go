package handler

func CreateViewURLPattern() string {
	urlStr := CreateURL(options.PathView, "*item")

	return urlStr
}

func CreateThumbnailURLPattern() string {
	urlStr := CreateURL(options.PathThumbnail, "*item")

	return urlStr
}

func CreateRescanURLPattern() string {
	return CreateURL(options.PathRescanLibrary)
}

func CreateGetImageURLPattern() string {
	url := CreateURL(options.PathGetImage, "*item")

	return url
}

func CreateUpdateCoverURLPattern() string {
	url := CreateURL(options.PathUpdateCover, "*item")

	return url
}

func CreateDownloadURLPattern() string {
	urlStr := CreateURL(options.PathDownload, "*item")

	return urlStr
}

func CreateSetFavoriteURLPattern() string {
	urlStr := CreateURL(options.PathFavorite, "*item")

	return urlStr
}

func CreateBrowseURLPattern() string {
	return CreateURL(options.PathBrowse, "*tag")
}

func CreateSetTagFavoriteURLPattern() string {
	urlStr := CreateURL(options.PathTagFavorite, "*item")

	return urlStr
}

func CreateTagListURLPattern() string {
	return CreateURL(options.PathTagList)
}

func CreateTagThumbnailURLPattern() string {
	urlStr := CreateURL(options.PathTagThumbnail, "*tag")

	return urlStr
}
