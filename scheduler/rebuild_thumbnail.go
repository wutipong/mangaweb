package scheduler

import (
	"github.com/wutipong/mangaweb/log"
	"github.com/wutipong/mangaweb/meta"
)

func RebuildThumbnail() error {
	allMeta, err := meta.ReadAll()
	if err != nil {
		return err
	}

	for _, m := range allMeta {
		e := m.GenerateThumbnail(0)
		log.Get().Sugar().Infof("Generating new thumbnail for %s", m.Name)
		if e != nil {
			log.Get().Sugar().Errorf("Failed to generate thumbnail for %s", m.Name)
			continue
		}

		meta.Write(m)
	}

	return nil
}

func ScheduleRebuildThumbnail() {
	scheduler.Every(1).Millisecond().LimitRunsTo(1).Do(func() {
		log.Get().Sugar().Info("Force updating thumbnail")
		RebuildThumbnail()
	})
}
