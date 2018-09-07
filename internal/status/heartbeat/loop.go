package heartbeat

import (
	"time"

	"github.com/modprox/libmodprox/loggy"
	"github.com/shoenig/toolkit"
)

type PokeLooper interface {
	Loop()
}

func NewLooper(interval time.Duration, sender Sender) PokeLooper {
	return &looper{
		interval: interval,
		sender:   sender,
		log:      loggy.New("heartbeat-looper"),

		numPackages: 1,
		numModules:  2,
	}
}

type looper struct {
	interval time.Duration
	sender   Sender
	log      loggy.Logger

	// todo: get real information
	numPackages int
	numModules  int
}

// Loop will block and run forever, sending heartbeats
// at the configured interval, to whichever of the specified
// registry instances accepts the heartbeat first.
func (l *looper) Loop() {
	toolkit.Interval(l.interval, l.loop)
}

func (l *looper) loop() error {
	// todo: get real information

	if err := l.sender.Send(
		l.numPackages,
		l.numModules,
	); err != nil {
		l.log.Warnf("could not send heartbeat, will try again later", err)
		return nil // always nil, never stop
	}

	return nil
}
