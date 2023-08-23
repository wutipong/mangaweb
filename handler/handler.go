package handler

var options Options

type Options struct {
	VersionString string

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
	PathTagFavorite   string
	PathTagList       string
	PathTagThumbnail  string
}

func Init(o Options) {
	options = o
}

func CreateVersionString() string {
	return options.VersionString
}
