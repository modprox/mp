package zips

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	// go-imports really hates these imports,
	// use aliases make them stick
	billyMemFS "gopkg.in/src-d/go-billy.v4/memfs"
	billyGit "gopkg.in/src-d/go-git.v4"
	billyPlumbing "gopkg.in/src-d/go-git.v4/plumbing"
	billyMem "gopkg.in/src-d/go-git.v4/storage/memory"
)

func Test_zipFS_file(t *testing.T) {
	fs := billyMemFS.New()
	f, err := fs.Create("foo.txt")
	require.NoError(t, err)
	f.Write([]byte("foo_bar_baz"))
	err = f.Close()
	require.NoError(t, err)

	blob, err := zipFS(fs)
	require.NoError(t, err)
	require.True(t, len(blob) > 0)
}

func Test_zipFS_dir(t *testing.T) {
	fs := billyMemFS.New()

	err := fs.MkdirAll("one/two/three", 0770)
	require.NoError(t, err)

	err = fs.MkdirAll("one/2/3", 0770)
	require.NoError(t, err)

	fOne, err := fs.Create("one/two/three/four.txt")
	require.NoError(t, err)

	contentOne := strings.NewReader("i am four.txt")
	_, err = io.Copy(fOne, contentOne)
	require.NoError(t, err)
	err = fOne.Close()
	require.NoError(t, err)

	f1, err := fs.Create("one/2/3/4.txt")
	require.NoError(t, err)

	content1 := strings.NewReader("i am 4.txt")
	_, err = io.Copy(f1, content1)
	require.NoError(t, err)
	err = f1.Close()
	require.NoError(t, err)

	blob, err := zipFS(fs)
	require.NoError(t, err)
	require.True(t, len(blob) > 0)
}

func Test_clone_tag(t *testing.T) {
	fs := billyMemFS.New()

	memoryStore := billyMem.NewStorage()
	repo, err := billyGit.Clone(memoryStore, fs, &billyGit.CloneOptions{
		URL:           "git@github.com:shoenig/toolkit.git",
		ReferenceName: "refs/tags/v1.0.0",
		Depth:         1,
		SingleBranch:  true,
	})
	require.NoError(t, err)

	iter, err := repo.Tags()
	require.NoError(t, err)

	err = iter.ForEach(func(ref *billyPlumbing.Reference) error {
		require.Equal(t, "refs/tags/v1.0.0", string(ref.Name()))
		return nil
	})
	require.NoError(t, err)
}
