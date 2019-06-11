package stats

import (
	"fmt"
	"time"

	"oss.indeed.com/go/modprox/pkg/since"

	"github.com/cactus/go-statsd-client/statsd"
)

type Service string

func (s Service) String() string {
	return string(s)
}

const (
	Proxy    Service = "modprox-proxy"
	Registry Service = "modprox-registry"
)

//go:generate minimock -g -i Sender -s _mock.go

// A Sender is used to emit statsd type metrics.
type Sender interface {
	Count(metric string, i int)
	Gauge(metric string, n int)
	GaugeMS(metric string, t time.Time)
}

// New creates a new Sender which will send metrics to the receiver described
// by the cfg.Agent configuration. All metrics will be emittited under the
// application named by Service s.
func New(s Service, cfg Statsd) (Sender, error) {
	address := fmt.Sprintf("%s:%d", cfg.Agent.Address, cfg.Agent.Port)
	emitter, err := statsd.NewClient(address, s.String())
	return &sender{
		emitter: emitter,
	}, err
}

type discard struct{}

func (d *discard) Count(string, int)         {}
func (d *discard) Gauge(string, int)         {}
func (d *discard) GaugeMS(string, time.Time) {}

func Discard() Sender {
	return &discard{}
}

type sender struct {
	emitter statsd.Statter
}

func (s *sender) Count(metric string, n int) {
	_ = s.emitter.Inc(metric, 1, 1)
}

func (s *sender) Gauge(metric string, n int) {
	_ = s.emitter.Gauge(metric, int64(n), 1)
}

// GaugeMS gauges the amount of time that has elapsed since t in milliseconds.
func (s *sender) GaugeMS(metric string, t time.Time) {
	_ = s.emitter.Gauge(metric, since.MS(t), 1)
}
