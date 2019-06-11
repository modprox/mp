package heartbeat

import (
	"time"

	"oss.indeed.com/go/modprox/pkg/metrics/stats"

	"github.com/shoenig/toolkit"

	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/proxy/internal/modules/store"
)

type PokeLooper interface {
	Loop()
}

func NewLooper(
	interval time.Duration,
	index store.Index,
	emitter stats.Sender,
	sender Sender,
) PokeLooper {
	return &looper{
		interval: interval,
		index:    index,
		emitter:  emitter,
		sender:   sender,
		log:      loggy.New("heartbeat-looper"),
	}
}

type looper struct {
	interval time.Duration
	index    store.Index
	sender   Sender
	emitter  stats.Sender
	log      loggy.Logger
}

// Loop will block and run forever, sending heartbeats
// at the configured interval, to whichever of the specified
// registry instances accepts the heartbeat first.
func (l *looper) Loop() {
	_ = toolkit.Interval(l.interval, l.loop)
}

func (l *looper) loop() error {
	modules, versions, err := l.index.Summary()
	if err != nil {
		return err
	}

	l.emitter.Gauge("index-num-modules", modules)   // really packages
	l.emitter.Gauge("index-num-versions", versions) // really modules

	if err := l.sender.Send(
		modules,
		versions,
	); err != nil {
		l.emitter.Count("heartbeat-send-failure", 1)

		l.log.Warnf("could not send heartbeat, will try again later: %v", err)
		return nil // always nil, never stop
	}

	l.emitter.Count("heartbeat-send-ok", 1)
	return nil
}
