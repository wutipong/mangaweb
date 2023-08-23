package scheduler

import (
	"github.com/wutipong/mangaweb/log"
	"github.com/wutipong/mangaweb/meta"
)

func MigrateMeta() error {
	allMeta, err := meta.ReadAll()
	if err != nil {
		return err
	}

	for _, m := range allMeta {
		if m.Version != meta.CurrentVersion {
			m, err = meta.Migrate(m)
			if err != nil {
				return err
			}

			err = meta.Write(m)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ScheduleMigrateMeta() {
	scheduler.Every(1).Millisecond().LimitRunsTo(1).Do(func() {
		log.Get().Sugar().Info("Upgrading metadata.")
		MigrateMeta()
	})
}
