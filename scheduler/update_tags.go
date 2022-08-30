package scheduler

import (
	"github.com/wutipong/mangaweb/log"
	"github.com/wutipong/mangaweb/meta"
	"github.com/wutipong/mangaweb/tag"
)

func UpdateTags() error {
	metaProvider, err := createMetaProvider()

	if err != nil {
		return err
	}
	defer metaProvider.Close()

	allMeta, err := metaProvider.ReadAll()
	if err != nil {
		return err
	}

	tagSet := make(map[string]bool)
	for _, m := range allMeta {
		tags := tag.ParseTag(m.Name)
		for _, t := range tags {
			tagSet[t] = true
		}
	}

	tagProvider, err := createTagProvider()
	if err != nil {
		return err
	}

	defer tagProvider.Close()

	allTag, err := tagProvider.ReadAll()
	allTagSet := make(map[string]bool)
	for _, t := range allTag {
		allTagSet[t.Name] = true
	}

	findMetaWithTag := func(tag string) meta.Meta {
		for _, m := range allMeta {
			for _, t := range m.Tags {
				if t == tag {
					return m
				}
			}
		}

		return meta.Meta{}
	}

	for tagStr := range tagSet {
		if _, found := allTagSet[tagStr]; !found {
			t := tag.NewTag(tagStr)
			m := findMetaWithTag(tagStr)
			t.Thumbnail = m.Thumbnail

			err = tagProvider.Write(t)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ScheduleUpdateTags() {
	scheduler.Every(1).Millisecond().LimitRunsTo(1).Do(func() {
		log.Get().Sugar().Info("Update tags.")
		UpdateTags()
	})
}
