package scheduler

import (
	"github.com/wutipong/mangaweb/log"
	"github.com/wutipong/mangaweb/meta"
)

func ScanLibrary() error {
	provider, err := createMetaProvider()

	if err != nil {
		return err
	}
	defer provider.Close()

	allMeta, err := provider.ReadAll()
	if err != nil {
		return err
	}

	files, err := meta.ListDir("")
	if err != nil {
		return err
	}

	for _, file := range files {
		found := false
		for _, m := range allMeta {
			if m.Name == file {
				found = true
				break
			}
		}
		if found {
			continue
		}

		log.Get().Sugar().Infof("Creating metadata for %s", file)

		item, err := meta.NewItem(file)
		if err != nil {
			log.Get().Sugar().Errorf("Failed to create meta data : %v", err)
		}

		err = provider.Write(item)
		if err != nil {
			log.Get().Sugar().Errorf("Failed to write meta data : %v", err)
		}
	}

	for _, m := range allMeta {
		found := false
		for _, file := range files {
			if m.Name == file {
				found = true
				break
			}
		}
		if found {
			continue
		}

		log.Get().Sugar().Infof("Deleting metadata for %s", m.Name)
		if err := provider.Delete(m); err != nil {
			log.Get().Sugar().Infof("Failed to delete meta for %s", m.Name)
		}

	}

	return nil
}

func ScheduleScanLibrary() {
	scheduler.Every(1).Millisecond().LimitRunsTo(1).Do(func() {
		log.Get().Sugar().Infof("Scanning Library.")
		ScanLibrary()
	})
}
