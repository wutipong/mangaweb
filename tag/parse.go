package tag

import (
	"regexp"
)

var regex *regexp.Regexp

func init() {
	const pattern = "\\[(.*?)\\]"
	var err error
	regex, err = regexp.Compile(pattern)

	if err != nil {
		panic(err)
	}
}

func ParseTag(name string) []string {
	matches := regex.FindAllStringSubmatch(name, -1)
	tagSet := make(map[string]bool)
	output := make([]string, 0)

	for i := 0; i < len(matches); i++ {
		tag := matches[i][1]
		if _, found := tagSet[tag]; !found {
			tagSet[tag] = true
			output = append(output, tag)
		}
	}

	return output
}
