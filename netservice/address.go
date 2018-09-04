package netservice

import "fmt"

type Instance struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func (s Instance) String() string {
	return fmt.Sprintf("<%s:%d>", s.Address, s.Port)
}
