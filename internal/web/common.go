package web

import (
	"strings"

	"github.com/modprox/libmodprox/repository"

	"github.com/pkg/errors"
)

// GET baseURL/module/@v/list fetches a list of all known versions, one per line.

func moduleFromPath(p string) (string, error) {
	vIdx := strings.Index(p, "@v")
	if vIdx <= 0 {
		return "", errors.Errorf("malformed path: %q", p)
	}
	return p[0:vIdx], nil
}

func modInfoFromPath(p string) (repository.ModInfo, error) {
	var mod repository.ModInfo
	split := strings.Split(p, "@v")
	if len(split) != 2 {
		return mod, errors.Errorf("malformed path: %q", p)
	}
	mod.Source = split[0]
	mod.Version = trimExt(split[1])
	return mod, nil
}

func trimExt(v string) string {
	dotIdx := strings.LastIndex(v, ".")
	if dotIdx <= 0 {
		return v
	}
	return v[:dotIdx]
}
