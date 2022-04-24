package meta

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoMigrateFromZero(t *testing.T) {
	funcs := make(map[int]MigrateFunction)
	funcs[0] = func(m Meta) (out Meta, err error) {
		out = m
		out.Name += "1"

		return
	}

	funcs[1] = func(m Meta) (out Meta, err error) {
		out = m
		out.Name += "2"

		return
	}

	funcs[2] = func(m Meta) (out Meta, err error) {
		out = m
		out.Name += "3"

		return
	}

	m := Meta{
		Name:    "",
		Version: 0,
	}

	target := 3
	out, err := doMigrate(m, funcs, target)

	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, out.Version, target)
	assert.Equal(t, out.Name, "123")
}

func TestDoMigrateFromNonZero(t *testing.T) {
	funcs := make(map[int]MigrateFunction)
	funcs[0] = func(m Meta) (out Meta, err error) {
		out = m
		out.Name += "1"

		return
	}

	funcs[1] = func(m Meta) (out Meta, err error) {
		out = m
		out.Name += "2"

		return
	}

	funcs[2] = func(m Meta) (out Meta, err error) {
		out = m
		out.Name += "3"

		return
	}

	m := Meta{
		Name:    "",
		Version: 1,
	}

	target := 3
	out, err := doMigrate(m, funcs, target)

	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, out.Version, target)
	assert.Equal(t, out.Name, "23")
}

func TestDoMigrateWithError(t *testing.T) {
	funcs := make(map[int]MigrateFunction)
	funcs[0] = func(m Meta) (out Meta, err error) {
		out = m
		out.Name += "1"

		return
	}

	funcs[1] = func(m Meta) (out Meta, err error) {
		err = fmt.Errorf("test error")
		return
	}

	funcs[2] = func(m Meta) (out Meta, err error) {
		out = m
		out.Name += "3"

		return
	}

	m := Meta{
		Name:    "",
		Version: 1,
	}

	target := 3
	out, err := doMigrate(m, funcs, target)

	assert.NotNil(t, err)
	assert.NotNil(t, out)
}
