package config

import "testing"
import "github.com/stretchr/testify/assert"

func TestConfiguration_MustCSRFAuthKey(t *testing.T) {
	var c Configuration
	assert.Panics(t, func() { c.MustCSRFAuthKey() }, "should panic with invalid csrf_auth_key")
	c.CSRFAuthKey = "this is a 32 byte long string ok"
	assert.Len(t, c.CSRFAuthKey, 32)
	assert.NotPanics(t, func() { c.MustCSRFAuthKey() }, "should not panic with valid csrf_auth_key")
}
