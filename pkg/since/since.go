// Package since provides convenience functions for computing elapsed time.
package since

import (
	"time"
)

// MS returns the number of milliseconds that have passed since t as int64.
func MS(t time.Time) int64 {
	return int64(1000 * time.Since(t).Seconds())
}

// MSFrom returns the number of milliseconds that have passed between t and now
// as an int64.
func MSFrom(t, now time.Time) int64 {
	return int64(1000 * now.Sub(t).Seconds())
}
