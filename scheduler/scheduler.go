package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/wutipong/mangaweb/log"
)

var scheduler *gocron.Scheduler

type Options struct {
}

func Init(options Options) {
	scheduler = gocron.NewScheduler(time.UTC)
	scheduler.Every(30).Minutes().Do(func() {
		log.Get().Sugar().Info("Update metadata set.")
		ScanLibrary()
		log.Get().Sugar().Info("Update tag list.")
		UpdateTags()
		log.Get().Sugar().Info("Update missing thumbnails.")
		UpdateMissingThumbnail()
	})
	ScheduleMigrateMeta()
}

func Start() {
	scheduler.StartAsync()
}

func Stop() {
	scheduler.Stop()
}
