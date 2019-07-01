package proxies

import (
	"time"

	"go.gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/registry/internal/data"
)

//go:generate go run github.com/gojuno/minimock/cmd/minimock -g -i Pruner -s _mock.go

type Pruner interface {
	Prune(time.Time) error
}

type pruner struct {
	maxAge time.Duration
	store  data.Store
	log    loggy.Logger
}

func NewPruner(maxAge time.Duration, store data.Store) Pruner {
	return &pruner{
		maxAge: maxAge,
		store:  store,
		log:    loggy.New("proxy-prune"),
	}
}

func (p *pruner) Prune(now time.Time) error {
	heartbeats, err := p.store.ListHeartbeats()
	if err != nil {
		return err
	}

	p.log.Tracef("looking through proxy heartbeats for removable instances")
	for _, heartbeat := range heartbeats {
		then := time.Unix(int64(heartbeat.Timestamp), 0)
		elapsed := now.Sub(then)
		if elapsed > p.maxAge {
			p.log.Warnf("purging M.I.A. proxy %s", heartbeat.Self)
			if err := p.store.PurgeProxy(heartbeat.Self); err != nil {
				p.log.Errorf("failed to purge proxy: %s: %v", heartbeat.Self, err)
				return err
			}
		} else {
			p.log.Tracef("not purging proxy of age %v (max %v)", elapsed, p.maxAge)
		}
	}

	return nil
}
