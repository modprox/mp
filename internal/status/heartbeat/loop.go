package heartbeat

import (
	"time"

	"github.com/shoenig/toolkit"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/modprox-proxy/internal/modules/store"
)

type PokeLooper interface {
	Loop()
}

func NewLooper(
	interval time.Duration,
	index store.Index,
	sender Sender,
) PokeLooper {
	return &looper{
		interval: interval,
		index:    index,
		sender:   sender,
		log:      loggy.New("heartbeat-looper"),
	}
}

type looper struct {
	interval time.Duration
	index    store.Index
	sender   Sender
	log      loggy.Logger
}

// Loop will block and run forever, sending heartbeats
// at the configured interval, to whichever of the specified
// registry instances accepts the heartbeat first.
func (l *looper) Loop() {
	toolkit.Interval(l.interval, l.loop)
}

func (l *looper) loop() error {
	// todo: get real information, depends on proxy #14
	// todo: when we have a boltdb and can maintain a separate
	// todo: table with this information indexed
	numPackages := 1
	numModules := 2

	if err := l.sender.Send(
		numPackages,
		numModules,
	); err != nil {
		l.log.Warnf("could not send heartbeat, will try again later: %v", err)
		return nil // always nil, never stop
	}

	return nil
}
