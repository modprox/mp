package fsutil

import (
	"strings"

	"github.com/pkg/errors"
)

func SafePath(p string) error {
	if strings.Contains(p, "..") {
		return errors.Errorf("unsafe path %q contains double dots", p)
	}
	return nil
}
