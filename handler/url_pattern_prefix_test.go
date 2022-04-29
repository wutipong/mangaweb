package handler

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type URLPatternPrefixTestSuite struct {
	suite.Suite
}

func (suite *URLPatternPrefixTestSuite) SetupSuite() {

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
}

func (suite *URLPatternPrefixTestSuite) TestCreateViewURLPattern() {
	u := CreateViewURLPattern()

	suite.Assert().Equal(
		"/manga/view/*item",
		u)
}

func (suite *URLPatternPrefixTestSuite) TestCreateRescanURLPattern() {
	u := CreateRescanURLPattern()
	suite.Assert().Equal("/manga/rescan_library", u)
}

func (suite *URLPatternPrefixTestSuite) TestCreateGetImageURLPattern() {
	u := CreateGetImageURLPattern()
	suite.Assert().Equal(
		"/manga/get_image/*item",
		u)
}

func (suite *URLPatternPrefixTestSuite) TestCreateUpdateCoverURLPattern() {
	u := CreateUpdateCoverURLPattern()
	suite.Assert().Equal(
		"/manga/update_cover/*item",
		u,
	)
}

func (suite *URLPatternPrefixTestSuite) TestCreateDownloadURLPattern() {
	u := CreateDownloadURLPattern()
	suite.Assert().Equal(
		"/manga/download/*item",
		u,
	)
}

func (suite *URLPatternPrefixTestSuite) TestCreateSetFavoriteURLPattern() {
	u := CreateSetFavoriteURLPattern()
	suite.Assert().Equal(
		"/manga/favorite/*item",
		u,
	)
}

func (suite *URLPatternPrefixTestSuite) TestCreateThumbnailURLPattern() {
	u := CreateThumbnailURLPattern()
	suite.Assert().Equal(
		"/manga/thumbnail/*item",
		u,
	)
}

func (suite *URLPatternPrefixTestSuite) TestCreateBrowseURLPattern() {
	u := CreateBrowseURLPattern()
	suite.Assert().Equal(
		"/manga/browse/*tag",
		u,
	)
}

func TestURLPatternPrefixTestSuite(t *testing.T) {
	suite.Run(t, new(URLPatternPrefixTestSuite))
}
