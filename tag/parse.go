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
	output := make([]string, len(matches))

	for i := 0; i < len(matches); i++ {
		output[i] = matches[i][1] // use the subgroup 1
	}

	return output
}
