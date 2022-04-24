package scheduler

import (
	"github.com/labstack/gommon/log"
	"github.com/wutipong/mangaweb/meta"
	"github.com/wutipong/mangaweb/tag"
	"sort"
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
	sort.Slice(allTag, func(i, j int) bool {
		return allTag[i].Name < allTag[j].Name
	})

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

	for tagStr, _ := range tagSet {
		if sort.Search(len(allTag), func(i int) bool {
			return allTag[i].Name == tagStr
		}) == len(allTag) {
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
		log.Info("Update tags.")
		UpdateTags()
	})
}
