package proxies

import (
	"errors"
	"testing"
	"time"

	"github.com/modprox/mp/pkg/clients/payloads"
	"github.com/modprox/mp/pkg/netservice"
	"github.com/modprox/mp/registry/internal/data/datatest"

	"github.com/stretchr/testify/require"
)

func Test_Prune(t *testing.T) {
	store := &datatest.Store{}
	defer store.AssertExpectations(t)

	now := time.Date(2018, 9, 27, 14, 48, 0, 0, time.UTC)
	oneMinuteAgo := now.Add(-1 * time.Minute)
	oneHourAgo := now.Add(-1 * time.Hour)

	store.On("ListHeartbeats").Return(
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
	).Once()

	// should only purge 2.2.2.2
	store.On("PurgeProxy", netservice.Instance{
		Address: "2.2.2.2",
		Port:    34567,
	}).Return(nil).Once()

	p := NewPruner(3*time.Minute, store)

	err := p.Prune(now)
	require.NoError(t, err)
}

func Test_Prune_list_fail(t *testing.T) {
	store := &datatest.Store{}
	defer store.AssertExpectations(t)

	now := time.Date(2018, 9, 27, 14, 48, 0, 0, time.UTC)

	store.On("ListHeartbeats").Return(
		[]payloads.Heartbeat{}, errors.New("db list fail"),
	).Once()

	p := NewPruner(3*time.Minute, store)

	err := p.Prune(now)
	require.Error(t, err)
}

func Test_Prune_purge_fail(t *testing.T) {
	store := &datatest.Store{}
	defer store.AssertExpectations(t)

	now := time.Date(2018, 9, 27, 14, 48, 0, 0, time.UTC)
	oneMinuteAgo := now.Add(-1 * time.Minute)
	oneHourAgo := now.Add(-1 * time.Hour)

	store.On("ListHeartbeats").Return(
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
	).Once()

	// should only purge 2.2.2.2
	store.On("PurgeProxy", netservice.Instance{
		Address: "2.2.2.2",
		Port:    34567,
	}).Return(errors.New("db purge fail")).Once()

	p := NewPruner(3*time.Minute, store)

	err := p.Prune(now)
	require.Error(t, err)
}
