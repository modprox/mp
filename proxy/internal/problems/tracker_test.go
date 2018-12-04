package problems

import (
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/modprox/mp/pkg/coordinates"

	"github.com/stretchr/testify/require"
)

func Test_Tracker_empty_default(t *testing.T) {
	pt := New("foo")

	problems := pt.Problems()
	require.Equal(t, 0, len(problems))

	_, exists := pt.Problem(coordinates.Module{
		Source:  "github.com/foo/bar",
		Version: "1.2.3",
	})
	require.False(t, exists)
}

func Test_Tracker_Set_one(t *testing.T) {
	pt := New("foo")

	pt.Set(Problem{
		Module: coordinates.Module{
			Source:  "github.com/foo/bar",
			Version: "1.2.3",
		},
		Time:    time.Date(2018, 12, 2, 20, 0, 0, 0, time.UTC),
		Message: "foobar is broken",
	})

	problems := pt.Problems()
	require.Equal(t, 1, len(problems))

	problem, exists := pt.Problem(coordinates.Module{
		Source:  "github.com/foo/bar",
		Version: "1.2.3",
	})
	require.True(t, exists)
	require.Equal(t, "foobar is broken", problem.Message)
}

func Test_byName(t *testing.T) {
	mod1 := coordinates.Module{
		Source:  "github.com/zzz/bar",
		Version: "1.0.0",
	}

	mod2 := coordinates.Module{
		Source:  "github.com/aaa/bar",
		Version: "2.2.3",
	}

	mod3 := coordinates.Module{
		Source:  "github.com/foo/bar",
		Version: "1.2.3",
	}

	mod4 := coordinates.Module{
		Source:  "github.com/aaa/bar",
		Version: "1.2.3",
	}

	mod5 := coordinates.Module{
		Source:  "github.com/bbb/bar",
		Version: "1.2.3",
	}

	mod6 := coordinates.Module{
		Source:  "github.com/foo/ccc",
		Version: "1.2.3",
	}

	problems := []Problem{
		Create(mod1, errors.New("m1")),
		Create(mod2, errors.New("m2")),
		Create(mod3, errors.New("m3")),
		Create(mod4, errors.New("m4")),
		Create(mod5, errors.New("m5")),
		Create(mod6, errors.New("m6")),
	}

	sort.Sort(byName(problems))

	// mod2 comes before mod4 in time
	require.Equal(t, mod2, problems[0].Module)

	// mod1 is zzz
	require.Equal(t, mod1, problems[5].Module)
}
