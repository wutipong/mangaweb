package scheduler

import "github.com/wutipong/mangaweb/log"

func UpdateMissingThumbnail() error {
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
		if len(m.Thumbnail) != 0 {
			continue
		}
		e := m.GenerateThumbnail(0)
		log.Get().Sugar().Infof("Re-generating new thumbnail for %s", m.Name)
		if e != nil {
			log.Get().Sugar().Errorf("Failed to generate thumbnail for %s", m.Name)
			continue
		}

		provider.Write(m)
	}

	return nil
}

func ScheduleUpdateMissingThumbnail() {
	scheduler.Every(1).Hour().Do(func() {
		log.Get().Sugar().Info("Updating missing thumbnail")
		UpdateMissingThumbnail()
	})
}
