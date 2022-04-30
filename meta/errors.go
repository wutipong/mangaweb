package meta

import "github.com/wutipong/mangaweb/errors"

var ErrMetaDataNotFound = errors.New(2_000_000, "metadata for '%s' not found.")
