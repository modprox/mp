package netservice

import (
	"os"
)

// Hostname returns the hostname or panics.
func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}
