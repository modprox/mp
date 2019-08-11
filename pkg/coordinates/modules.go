package coordinates

import (
	"fmt"

	"gophers.dev/pkgs/semantic"
)

type Module struct {
	Source  string `json:"source"`
	Version string `json:"version"`
}

// String representation intended for human consumption.
//
// Includes surrounding parenthesis and some whitespace to pop in logs.
func (m Module) String() string {
	return fmt.Sprintf("(%s @ %s)", m.Source, m.Version)
}

// AtVersion representation intended for machine consumption.
//
// Format is source@version.
func (m Module) AtVersion() string {
	return fmt.Sprintf(
		"%s@%s",
		m.Source,
		m.Version,
	)
}

// Bytes is AtVersion but converted to Bytes for use in a data-store which
// stores bytes.
func (m Module) Bytes() []byte {
	return []byte(m.AtVersion())
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

	tagX, okX := semantic.Parse(modX.Version)
	tagY, okY := semantic.Parse(modY.Version)

	if !okX && !okY {
		return false
	} else if okX && !okY {
		return false
	} else if !okX && okY {
		return true
	}

	return tagX.Less(tagY)
}
