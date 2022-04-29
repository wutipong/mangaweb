package handler

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type URLPatternTestSuite struct {
	suite.Suite
}

func (suite *URLPatternTestSuite) SetupSuite() {

	Init(Options{
		PathPrefix:        "",
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

func (suite *URLPatternTestSuite) TestCreateViewURLPattern() {
	u := CreateViewURLPattern()

	suite.Assert().Equal(
		"/view/*item",
		u)
}

func (suite *URLPatternTestSuite) TestCreateRescanURLPattern() {
	u := CreateRescanURLPattern()
	suite.Assert().Equal("/rescan_library", u)
}

func (suite *URLPatternTestSuite) TestCreateGetImageURLPattern() {
	u := CreateGetImageURLPattern()
	suite.Assert().Equal(
		"/get_image/*item",
		u)
}

func (suite *URLPatternTestSuite) TestCreateUpdateCoverURLPattern() {
	u := CreateUpdateCoverURLPattern()
	suite.Assert().Equal(
		"/update_cover/*item",
		u,
	)
}

func (suite *URLPatternTestSuite) TestCreateDownloadURLPattern() {
	u := CreateDownloadURLPattern()
	suite.Assert().Equal(
		"/download/*item",
		u,
	)
}

func (suite *URLPatternTestSuite) TestCreateSetFavoriteURLPattern() {
	u := CreateSetFavoriteURLPattern()
	suite.Assert().Equal(
		"/favorite/*item",
		u,
	)
}

func (suite *URLPatternTestSuite) TestCreateThumbnailURLPattern() {
	u := CreateThumbnailURLPattern()
	suite.Assert().Equal(
		"/thumbnail/*item",
		u,
	)
}

func (suite *URLPatternTestSuite) TestCreateBrowseURLPattern() {
	u := CreateBrowseURLPattern()
	suite.Assert().Equal(
		"/browse/*tag",
		u,
	)
}

func TestURLPatternTestSuite(t *testing.T) {
	suite.Run(t, new(URLPatternTestSuite))
}
