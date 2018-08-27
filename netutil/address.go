package netutil

import "fmt"

type Service struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func (s Service) String() string {
	return fmt.Sprintf("<%s:%d>", s.Address, s.Port)
}
