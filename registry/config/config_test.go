package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Configuration_CSRF_nokey(t *testing.T) {
	var c Configuration
	_, _, err := c.CSRF.Key()
	require.Error(t, err)
}

func Test_Configuration_CSRF_badkey(t *testing.T) {
	var c Configuration
	c.CSRF.AuthenticationKey = "foobar"
	_, _, err := c.CSRF.Key()
	require.Error(t, err)
}

func Test_Configuration_CSRF_devmode(t *testing.T) {
	key := "12345678901234567890123456789012"
	var c Configuration
	c.CSRF.DevelopmentMode = true
	c.CSRF.AuthenticationKey = key
	bs, devMode, err := c.CSRF.Key()
	require.NoError(t, err)
	require.True(t, devMode)
	require.Equal(t, []byte(key), bs)
}

func Test_Configuration_CSRF(t *testing.T) {
	key := "12345678901234567890123456789012"
	var c Configuration
	c.CSRF.AuthenticationKey = key
	bs, devMode, err := c.CSRF.Key()
	require.NoError(t, err)
	require.False(t, devMode)
	require.Equal(t, []byte(key), bs)
}
