package heartbeat

import (
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/shoenig/toolkit"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/proxy/internal/modules/store"
)

type PokeLooper interface {
	Loop()
}

func NewLooper(
	interval time.Duration,
	index store.Index,
	statter statsd.Statter,
	sender Sender,
) PokeLooper {
	return &looper{
		interval: interval,
		index:    index,
		sender:   sender,
		statter:  statter,
		log:      loggy.New("heartbeat-looper"),
	}
}

type looper struct {
	interval time.Duration
	index    store.Index
	sender   Sender
	statter  statsd.Statter
	log      loggy.Logger
}

// Loop will block and run forever, sending heartbeats
// at the configured interval, to whichever of the specified
// registry instances accepts the heartbeat first.
func (l *looper) Loop() {
	toolkit.Interval(l.interval, l.loop)
}

func (l *looper) loop() error {
	modules, versions, err := l.index.Summary()
	if err != nil {
		return err
	}

	l.statter.Gauge("index-num-modules", int64(modules), 1)
	l.statter.Gauge("index-num-versions", int64(versions), 1)

	if err := l.sender.Send(
		modules,
		versions,
	); err != nil {
		l.statter.Inc("heartbeat-send-failure", 1, 1)
		l.log.Warnf("could not send heartbeat, will try again later: %v", err)
		return nil // always nil, never stop
	}

	l.statter.Inc("heartbeat-send-ok", 1, 1)
	return nil
}
