package tag

import (
	"regexp"
)

func ParseTag(name string) []string {
	const pattern = "\\[(.*?)\\]"

	regex, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	matches := regex.FindAllStringSubmatch(name, -1)
	output := make([]string, len(matches))

	for i := 0; i < len(matches); i++ {
		output[i] = matches[i][1] // use the subgroup 1
	}

	return output
}
