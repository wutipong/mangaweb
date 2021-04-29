package urlutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateURL(t *testing.T) {
	SetPrefix("manga")

	u := CreateURL("/browse/?abcdefg")

	assert.Equal(t, u, "manga/browse/?abcdefg")
}

func TestCreateURLTwoParam(t *testing.T) {
	SetPrefix("manga")

	u := CreateURL("/browse", "abcdefg")

	assert.Equal(t, u, "manga/browse/abcdefg")
}

func TestCreateURLNoPrefix(t *testing.T) {
	SetPrefix("")

	u := CreateURL("/browse", "abcdefg")

	assert.Equal(t, u, "/browse/abcdefg")
}
