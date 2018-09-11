package netservice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_String(t *testing.T) {
	try := func(input Instance, exp string) {
		output := input.String()
		require.Equal(t, exp, output)
	}

	try(Instance{
		Address: "1.1.1.1",
		Port:    1111,
	}, "[1.1.1.1:1111]")
}
