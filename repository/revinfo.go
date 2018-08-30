package repository

import (
	"encoding/json"
	"time"
)

type RevInfo struct {
	Version string    `json:"version,omitempty"` // version string
	Name    string    `json:"name,omitempty"`    // complete ID in underlying repository
	Short   string    `json:"short,omitempty"`   // shortened ID, for use in pseudo-version
	Time    time.Time `json:"time,omitempty"`    // commit time
}

func (ri RevInfo) String() string {
	bs, err := json.MarshalIndent(ri, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bs)
}
