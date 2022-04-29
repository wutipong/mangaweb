package scheduler

import "github.com/wutipong/mangaweb/log"

func RebuildThumbnail() error {
	provider, err := createMetaProvider()

	if err != nil {
		return err
	}
	defer provider.Close()

	allMeta, err := provider.ReadAll()
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

		provider.Write(m)
	}

	return nil
}

func ScheduleRebuildThumbnail() {
	scheduler.Every(1).Millisecond().LimitRunsTo(1).Do(func() {
		log.Get().Sugar().Info("Force updating thumbnail")
		RebuildThumbnail()
	})
}
