package coordinates

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func sMod(id int64, source, version string) SerialModule {
	return SerialModule{
		SerialID: id,
		Module: Module{
			Source:  source,
			Version: version,
		},
	}
}

func Test_ModsByVersion(t *testing.T) {
	mods := []SerialModule{
		sMod(3, "foo", "v1.2.13"),
		sMod(3, "bar", "v0.0.3"),
		sMod(3, "bar", "v1.0.10"),
		sMod(3, "bar", "v1.2.1"),
		sMod(3, "baz", "v1.2.3"),
		sMod(3, "bar", "v11.3.3"),
		sMod(3, "foo", "v2.0.14"),
		sMod(3, "bar", "v11.2.3"),
		sMod(3, "baz", "v3.2.3"),
		sMod(3, "baz", "v3.12.1"),
		sMod(3, "bar", "v1.20.3"),
	}

	sort.Sort(ModsByVersion(mods))

	require.Equal(t, []SerialModule{
		sMod(3, "bar", "v0.0.3"),
		sMod(3, "bar", "v1.0.10"),
		sMod(3, "bar", "v1.2.1"),
		sMod(3, "bar", "v1.20.3"),
		sMod(3, "bar", "v11.2.3"),
		sMod(3, "bar", "v11.3.3"),
		sMod(3, "baz", "v1.2.3"),
		sMod(3, "baz", "v3.2.3"),
		sMod(3, "baz", "v3.12.1"),
		sMod(3, "foo", "v1.2.13"),
		sMod(3, "foo", "v2.0.14"),
	}, mods)
}
