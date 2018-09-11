package coordinates

import "fmt"

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
