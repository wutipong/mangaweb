package scheduler

import (
	"github.com/go-co-op/gocron"
	"github.com/labstack/gommon/log"
	"github.com/wutipong/mangaweb/meta"
	"github.com/wutipong/mangaweb/tag"
	"time"
)

var metaProviderFactory meta.ProviderFactory
var tagProviderFactory tag.ProviderFactory

var scheduler *gocron.Scheduler

type Options struct {
	MetaProviderFactory meta.ProviderFactory
	TagProviderFactory  tag.ProviderFactory
}

func Init(options Options) {
	metaProviderFactory = options.MetaProviderFactory
	tagProviderFactory = options.TagProviderFactory

	scheduler = gocron.NewScheduler(time.UTC)
	scheduler.Every(30).Minutes().Do(func() {
		log.Info("Update metadata set.")
		ScanLibrary()
		log.Info("Update tag list.")
		UpdateTags()
	})
}

func Start() {
	scheduler.StartAsync()
}

func Stop() {
	scheduler.Stop()
}

func createMetaProvider() (p meta.Provider, err error) {
	return metaProviderFactory()
}

func createTagProvider() (p tag.Provider, err error) {
	return tagProviderFactory()
}
