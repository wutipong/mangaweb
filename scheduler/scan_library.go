package scheduler

import (
	"github.com/labstack/gommon/log"
	"github.com/wutipong/mangaweb/meta"
)

func ScanLibrary() error {
	provider, err := metaProviderFactory()

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

		log.Printf("Creating metadata for %s", file)

		item, err := meta.NewItem(file)
		if err != nil {
			log.Printf("Failed to create meta data : %v", err)
		}

		err = provider.Write(item)
		if err != nil {
			log.Printf("Failed to write meta data : %v", err)
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

		log.Printf("Deleting metadata for %s", m.Name)
		if err := provider.Delete(m); err != nil {
			log.Printf("Failed to delete meta for %s", m.Name)
		}

	}

	return nil
}

func ScheduleScanLibrary() {
	scheduler.Every(1).Millisecond().LimitRunsTo(1).Do(func() {
		log.Info("Scanning Library.")
		ScanLibrary()
	})
}
