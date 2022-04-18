package scheduler

import (
	"github.com/go-co-op/gocron"
	"github.com/labstack/gommon/log"
	"github.com/wutipong/mangaweb/meta"
	"time"
)

var metaProviderFactory meta.MetaProviderFactory
var scheduler *gocron.Scheduler

func Init(factory meta.MetaProviderFactory) {
	metaProviderFactory = factory

	scheduler = gocron.NewScheduler(time.UTC)
	scheduler.Every(30).Minutes().Do(func() {
		log.Info("Update metadata set.")
		ScanLibrary()
	})
}

func Start() {
	scheduler.StartAsync()
}

func Stop() {
	scheduler.Stop()
}
