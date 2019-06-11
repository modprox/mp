package proxies

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"oss.indeed.com/go/modprox/pkg/clients/payloads"
	"oss.indeed.com/go/modprox/pkg/netservice"
	"oss.indeed.com/go/modprox/registry/internal/data"
)

func Test_Prune(t *testing.T) {
	store := data.NewStoreMock(t)
	defer store.MinimockFinish()

	now := time.Date(2018, 9, 27, 14, 48, 0, 0, time.UTC)
	oneMinuteAgo := now.Add(-1 * time.Minute)
	oneHourAgo := now.Add(-1 * time.Hour)

	store.ListHeartbeatsMock.Return(
		[]payloads.Heartbeat{
			{
				Self: netservice.Instance{
					Address: "1.1.1.1",
					Port:    23456,
				},
				Timestamp: int(oneMinuteAgo.Unix()),
			},
			{
				Self: netservice.Instance{
					Address: "2.2.2.2",
					Port:    34567,
				},
				Timestamp: int(oneHourAgo.Unix()),
			},
		}, nil,
	)

	// should only purge 2.2.2.2
	store.PurgeProxyMock.When(netservice.Instance{
		Address: "2.2.2.2",
		Port:    34567,
	}).Then(nil)

	p := NewPruner(3*time.Minute, store)

	err := p.Prune(now)
	require.NoError(t, err)
}

func Test_Prune_list_fail(t *testing.T) {
	store := data.NewStoreMock(t)
	defer store.MinimockFinish()

	now := time.Date(2018, 9, 27, 14, 48, 0, 0, time.UTC)

	store.ListHeartbeatsMock.Return([]payloads.Heartbeat{}, errors.New("db list fail"))

	p := NewPruner(3*time.Minute, store)

	err := p.Prune(now)
	require.Error(t, err)
}

func Test_Prune_purge_fail(t *testing.T) {
	store := data.NewStoreMock(t)
	defer store.MinimockFinish()

	now := time.Date(2018, 9, 27, 14, 48, 0, 0, time.UTC)
	oneMinuteAgo := now.Add(-1 * time.Minute)
	oneHourAgo := now.Add(-1 * time.Hour)

	store.ListHeartbeatsMock.Return(
		[]payloads.Heartbeat{
			{
				Self: netservice.Instance{
					Address: "1.1.1.1",
					Port:    23456,
				},
				Timestamp: int(oneMinuteAgo.Unix()),
			},
			{
				Self: netservice.Instance{
					Address: "2.2.2.2",
					Port:    34567,
				},
				Timestamp: int(oneHourAgo.Unix()),
			},
		}, nil,
	)

	// should only purge 2.2.2.2
	store.PurgeProxyMock.When(netservice.Instance{
		Address: "2.2.2.2",
		Port:    34567,
	}).Then(errors.New("db purge fail"))

	p := NewPruner(3*time.Minute, store)

	err := p.Prune(now)
	require.Error(t, err)
}
