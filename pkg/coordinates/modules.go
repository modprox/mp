package coordinates

import (
	"fmt"

	"github.com/modprox/taggit/tags"
)

type Module struct {
	Source  string `json:"source"`
	Version string `json:"version"`
}

func (mi Module) String() string {
	return fmt.Sprintf("(%s @ %s)", mi.Source, mi.Version)
}

func (mi Module) Bytes() []byte {
	return []byte(fmt.Sprintf(
		"%s@%s",
		mi.Source,
		mi.Version,
	))
}

type SerialModule struct {
	Module
	SerialID int64 `json:"id"`
}

type ModsByVersion []SerialModule

func (mods ModsByVersion) Len() int      { return len(mods) }
func (mods ModsByVersion) Swap(x, y int) { mods[x], mods[y] = mods[y], mods[x] }
func (mods ModsByVersion) Less(x, y int) bool {
	modX := mods[x]
	modY := mods[y]

	if modX.Source < modY.Source {
		return true
	} else if modX.Source > modY.Source {
		return false
	}

	tagX, okX := tags.Parse(modX.Version)
	tagY, okY := tags.Parse(modY.Version)

	if !okX && !okY {
		return false
	} else if okX && !okY {
		return false
	} else if !okX && okY {
		return true
	}

	return tagX.Less(tagY)
}
