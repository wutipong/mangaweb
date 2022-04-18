package scheduler

import "github.com/labstack/gommon/log"

func RebuildThumbnail() error {
	provider, err := metaProviderFactory()

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
		log.Printf("Generating new thumbnail for %s", m.Name)
		if e != nil {
			log.Printf("Failed to generate thumbnail for %s", m.Name)
			continue
		}

		provider.Write(m)
	}

	return nil
}

func ScheduleRebuildThumbnail() {
	scheduler.Every(1).Millisecond().LimitRunsTo(1).Do(func() {
		log.Info("Force updating thumbnail")
		RebuildThumbnail()
	})
}
