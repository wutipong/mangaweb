package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTagOne(t *testing.T) {
	assert.ElementsMatch(t, ParseTag("[Test]Some weird name"), []string{"Test"})
}

func TestParseTagNoTag(t *testing.T) {
	assert.Empty(t, ParseTag("Hello World"))
}

func TestParseTagMultiple(t *testing.T) {
	assert.ElementsMatch(t, ParseTag("[Test]Some weird name [Download]"), []string{"Test", "Download"})
}

func TestParseTagDuplicate(t *testing.T) {
	assert.ElementsMatch(t, ParseTag("[Test]something[Download]/[Test]Some weird name [Download]"), []string{"Test", "Download"})
}
