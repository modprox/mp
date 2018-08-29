package repository

import "fmt"

type ModInfo struct {
	Source  string `json:"source"`
	Version string `json:"version"`
}

func (mi ModInfo) String() string {
	return fmt.Sprintf("(%s @ %s)", mi.Source, mi.Version)
}
