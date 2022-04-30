package tag

import "github.com/wutipong/mangaweb/errors"

var ErrTagNotFound = errors.New(1_000_000, "tag '%s' not found.")
