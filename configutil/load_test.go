package configutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetConfigFilename_ok(t *testing.T) {
	args := []string{"executable", "foo.json"}
	filename, err := GetConfigFilename(args)
	require.NoError(t, err)
	require.Equal(t, "foo.json", filename)
}

func Test_GetConfigFilename_too_few(t *testing.T) {
	args := []string{"executable"}
	_, err := GetConfigFilename(args)
	require.Error(t, err)
}

func Test_GetConfigFilename_too_many(t *testing.T) {
	args := []string{"executable", "bar.json", "baz.json"}
	_, err := GetConfigFilename(args)
	require.Error(t, err)
}

type c struct {
	Foo string `json:"foo"`
}

func setup(t *testing.T, content string) {
	err := ioutil.WriteFile("foo.txt", []byte(content), 0600)
	require.NoError(t, err)
}

func cleanup(t *testing.T) {
	err := os.Remove("foo.txt")
	require.NoError(t, err)
}

func Test_LoadConfig_ok(t *testing.T) {
	setup(t, `{"foo":"bar"}`)
	defer cleanup(t)

	var config c
	err := LoadConfig("foo.txt", &config)
	require.NoError(t, err)
	require.Equal(t, "bar", config.Foo)
}

func Test_LoadConfig_unparsable(t *testing.T) {
	setup(t, "{{{{")
	defer cleanup(t)

	var config c
	err := LoadConfig("", &config)
	require.Error(t, err)
}

func Test_LoadConfig_no_file(t *testing.T) {
	var config c
	err := LoadConfig("/does/not/exist", &config)
	require.Error(t, err)
}
