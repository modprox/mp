package store

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/modprox/libmodprox/repository"

	"github.com/stretchr/testify/require"
)

func setupIndex(t *testing.T) (string, Index) {
	tmpDir, err := ioutil.TempDir("", "index-")
	require.NoError(t, err)

	options := IndexOptions{
		Directory: tmpDir,
	}

	index, err := NewIndex(options)
	require.NoError(t, err)

	return tmpDir, index
}

func cleanupIndex(t *testing.T, tmpDir string) {
	err := os.RemoveAll(tmpDir)
	require.NoError(t, err)
}

func Test_Index_empty(t *testing.T) {
	tmpDir, index := setupIndex(t)
	defer cleanupIndex(t, tmpDir)

	versions, err := index.Versions("github.com/pkg/errors")
	require.NoError(t, err)
	require.Equal(t, 0, len(versions))

	_, err = index.Info(repository.ModInfo{
		Source:  "github.com/pkg/errors",
		Version: "v0.8.0",
	})
	require.Error(t, err)

	_, err = index.Mod(repository.ModInfo{
		Source:  "github.com/pkg/errors",
		Version: "v0.8.0",
	})
	require.Error(t, err)

	exists, err := index.Contains(repository.ModInfo{
		Source:  "github.com/pkg/errors",
		Version: "v0.8.0",
	})
	require.NoError(t, err)
	require.False(t, exists)
}

func Test_Index_Put_1(t *testing.T) {
	tmpDir, index := setupIndex(t)
	defer cleanupIndex(t, tmpDir)

	err := index.Put(ModuleAddition{
		Mod: repository.ModInfo{
			Source:  "github.com/pkg/errors",
			Version: "v0.8.0",
		},
		ModFile:  "module github.com/pkg/errors",
		UniqueID: 1,
	})
	require.NoError(t, err)

	versions, err := index.Versions("github.com/pkg/errors")
	require.NoError(t, err)
	require.Equal(t, 1, len(versions))

	// not the module added
	versions, err = index.Versions("gitlab.com/some/other")
	require.NoError(t, err)
	require.Equal(t, 0, len(versions))

	info, err := index.Info(repository.ModInfo{
		Source:  "github.com/pkg/errors",
		Version: "v0.8.0",
	})
	require.NoError(t, err)
	require.Equal(t, repository.RevInfo{
		Version: "v0.8.0",
	}, info)

	// not the module added
	_, err = index.Info(repository.ModInfo{
		Source:  "gitlab.com/some/other",
		Version: "v0.0.0",
	})
	require.Error(t, err)

	modFile, err := index.Mod(repository.ModInfo{
		Source:  "github.com/pkg/errors",
		Version: "v0.8.0",
	})
	require.NoError(t, err)
	require.Equal(t, "module github.com/pkg/errors", modFile)

	// not the module added
	_, err = index.Mod(repository.ModInfo{
		Source:  "gitlab.com/some/other",
		Version: "v0.0.0",
	})
	require.Error(t, err)
}

// todo: test were we put several in, and test version sorting
