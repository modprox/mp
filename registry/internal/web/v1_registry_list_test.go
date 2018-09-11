package web

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/modprox/mp/pkg/coordinates"
)

type Range = coordinates.RangeID

type Ranges = coordinates.RangeIDs

func Test_inRange(t *testing.T) {
	try := func(i int64, rangeID Range, exp bool) {
		output := inRange(i, rangeID)
		require.Equal(t, exp, output)
	}

	try(1, Range{1, 1}, true)
	try(1, Range{1, 5}, true)
	try(1, Range{2, 5}, false)
	try(2, Range{1, 3}, true)
	try(2, Range{2, 3}, true)
	try(2, Range{1, 1}, false)
	try(2, Range{3, 5}, false)
	try(10, Range{3, 9}, false)
	try(10, Range{3, 13}, true)
	try(10, Range{11, 30}, false)
}

func Test_inListButNotRange(t *testing.T) {
	try := func(ids []int64, ranges Ranges, exp []int64) {
		output := inListButNotRange(ids, ranges)
		require.Equal(t, exp, output)
	}

	try(
		[]int64{1, 2, 3},
		Ranges{{1, 2}, {5, 6}},
		[]int64{3},
	)

	try(
		[]int64{4, 5, 6, 10, 11, 12},
		Ranges{{4, 4}, {11, 12}},
		[]int64{5, 6, 10},
	)
}
