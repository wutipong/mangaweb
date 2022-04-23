package handler

import (
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
)

type URLPrefixTestSuite struct {
	suite.Suite
	FilePath string
}

func (suite *URLPrefixTestSuite) SetupSuite() {

	Init(Options{
		PathPrefix:        "/manga",
		PathRoot:          "/",
		PathBrowse:        "/browse",
		PathView:          "/view",
		PathStatic:        "/static",
		PathGetImage:      "/get_image",
		PathUpdateCover:   "/update_cover",
		PathThumbnail:     "/thumbnail",
		PathFavorite:      "/favorite",
		PathDownload:      "/download",
		PathRescanLibrary: "/rescan_library",
	})

	suite.FilePath = filepath.Join("同人マンガ", "エロいまんが.zip")
}

func (suite *URLPrefixTestSuite) TestCreateURL() {
	u := CreateURL("/browse/?abcdefg")

	suite.Assert().Equal("/manga/browse/?abcdefg", u)
}

func (suite *URLPrefixTestSuite) TestCreateURLTwoParam() {
	u := CreateURL("/browse", "abcdefg")

	suite.Assert().Equal(u, "/manga/browse/abcdefg")
}

func (suite *URLPrefixTestSuite) TestCreateViewURL() {
	u := CreateViewURL(suite.FilePath)

	suite.Assert().Equal(
		"/manga/view/%E5%90%8C%E4%BA%BA%E3%83%9E%E3%83%B3%E3%82%AC/%E3%82%A8%E3%83%AD%E3%81%84%E3%81%BE%E3%82%93%E3%81%8C.zip",
		u)
}

func (suite *URLPrefixTestSuite) TestCreateRescanURL() {
	u := CreateRescanURL()
	suite.Assert().Equal("/manga/rescan_library", u)
}

func (suite *URLPrefixTestSuite) TestCreateGetImageURL() {
	u := CreateGetImageURL(suite.FilePath, 50)
	suite.Assert().Equal(
		"/manga/get_image/%E5%90%8C%E4%BA%BA%E3%83%9E%E3%83%B3%E3%82%AC/%E3%82%A8%E3%83%AD%E3%81%84%E3%81%BE%E3%82%93%E3%81%8C.zip?i=50",
		u)
}

func (suite *URLPrefixTestSuite) TestCreateUpdateCoverURL() {
	u := CreateUpdateCoverURL(suite.FilePath, 50)
	suite.Assert().Equal(
		"/manga/update_cover/%E5%90%8C%E4%BA%BA%E3%83%9E%E3%83%B3%E3%82%AC/%E3%82%A8%E3%83%AD%E3%81%84%E3%81%BE%E3%82%93%E3%81%8C.zip?i=50",
		u,
	)
}

func (suite *URLPrefixTestSuite) TestCreateDownloadURL() {
	u := CreateDownloadURL(suite.FilePath)
	suite.Assert().Equal(
		"/manga/download/%E5%90%8C%E4%BA%BA%E3%83%9E%E3%83%B3%E3%82%AC/%E3%82%A8%E3%83%AD%E3%81%84%E3%81%BE%E3%82%93%E3%81%8C.zip",
		u,
	)
}

func (suite *URLPrefixTestSuite) TestCreateSetFavoriteURL() {
	u := CreateSetFavoriteURL(suite.FilePath)
	suite.Assert().Equal(
		"/manga/favorite/%E5%90%8C%E4%BA%BA%E3%83%9E%E3%83%B3%E3%82%AC/%E3%82%A8%E3%83%AD%E3%81%84%E3%81%BE%E3%82%93%E3%81%8C.zip",
		u,
	)
}

func (suite *URLPrefixTestSuite) TestCreateThumbnailURL() {
	u := CreateThumbnailURL(suite.FilePath)
	suite.Assert().Equal(
		"/manga/thumbnail/%E5%90%8C%E4%BA%BA%E3%83%9E%E3%83%B3%E3%82%AC/%E3%82%A8%E3%83%AD%E3%81%84%E3%81%BE%E3%82%93%E3%81%8C.zip",
		u,
	)
}

func (suite *URLPrefixTestSuite) TestCreateBrowseURL() {
	u := CreateBrowseURL("aabbcc")
	suite.Assert().Equal(
		"/manga/browse#aabbcc",
		u,
	)
}

func TestURLPrefixTestSuite(t *testing.T) {
	suite.Run(t, new(URLPrefixTestSuite))
}
