package handler

import (
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
)

type URLTestSuite struct {
	suite.Suite
	FilePath string
}

func (suite *URLTestSuite) SetupSuite() {
	Init(Options{
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

func (suite *URLTestSuite) TestCreateURL() {
	u := CreateURL("/browse/?abcdefg")

	suite.Assert().Equal("/browse/?abcdefg", u)
}

func (suite *URLTestSuite) TestCreateURLTwoParam() {
	u := CreateURL("/browse", "abcdefg")

	suite.Assert().Equal(u, "/browse/abcdefg")
}

func (suite *URLTestSuite) TestCreateViewURL() {
	u := CreateViewURL(suite.FilePath)

	suite.Assert().Equal(
		"/view/%E5%90%8C%E4%BA%BA%E3%83%9E%E3%83%B3%E3%82%AC/%E3%82%A8%E3%83%AD%E3%81%84%E3%81%BE%E3%82%93%E3%81%8C.zip",
		u)
}

func (suite *URLTestSuite) TestCreateRescanURL() {
	u := CreateRescanURL()
	suite.Assert().Equal("/rescan_library", u)
}

func (suite *URLTestSuite) TestCreateGetImageURL() {
	u := CreateGetImageURL(suite.FilePath, 50)
	suite.Assert().Equal(
		"/get_image/%E5%90%8C%E4%BA%BA%E3%83%9E%E3%83%B3%E3%82%AC/%E3%82%A8%E3%83%AD%E3%81%84%E3%81%BE%E3%82%93%E3%81%8C.zip?i=50",
		u)
}

func (suite *URLTestSuite) TestCreateUpdateCoverURL() {
	u := CreateUpdateCoverURL(suite.FilePath, 50)
	suite.Assert().Equal(
		"/update_cover/%E5%90%8C%E4%BA%BA%E3%83%9E%E3%83%B3%E3%82%AC/%E3%82%A8%E3%83%AD%E3%81%84%E3%81%BE%E3%82%93%E3%81%8C.zip?i=50",
		u,
	)
}

func (suite *URLTestSuite) TestCreateDownloadURL() {
	u := CreateDownloadURL(suite.FilePath)
	suite.Assert().Equal(
		"/download/%E5%90%8C%E4%BA%BA%E3%83%9E%E3%83%B3%E3%82%AC/%E3%82%A8%E3%83%AD%E3%81%84%E3%81%BE%E3%82%93%E3%81%8C.zip",
		u,
	)
}

func (suite *URLTestSuite) TestCreateSetFavoriteURL() {
	u := CreateSetFavoriteURL(suite.FilePath)
	suite.Assert().Equal(
		"/favorite/%E5%90%8C%E4%BA%BA%E3%83%9E%E3%83%B3%E3%82%AC/%E3%82%A8%E3%83%AD%E3%81%84%E3%81%BE%E3%82%93%E3%81%8C.zip",
		u,
	)
}

func (suite *URLTestSuite) TestCreateThumbnailURL() {
	u := CreateThumbnailURL(suite.FilePath)
	suite.Assert().Equal(
		"/thumbnail/%E5%90%8C%E4%BA%BA%E3%83%9E%E3%83%B3%E3%82%AC/%E3%82%A8%E3%83%AD%E3%81%84%E3%81%BE%E3%82%93%E3%81%8C.zip",
		u,
	)
}

func (suite *URLTestSuite) TestCreateBrowseURL() {
	u := CreateBrowseURL("aabbcc")
	suite.Assert().Equal(
		"/browse#aabbcc",
		u,
	)
}

func TestURLTestSuite(t *testing.T) {
	suite.Run(t, new(URLTestSuite))
}
