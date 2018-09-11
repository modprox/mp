package configutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Format(t *testing.T) {
	config := struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}{
		Foo: "red",
		Bar: 8,
	}

	formatted := Format(config)
	require.JSONEq(t, `{"foo":"red", "bar":8}`, formatted)
}
