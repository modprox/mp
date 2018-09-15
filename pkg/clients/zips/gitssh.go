package zips

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/repository"
	"github.com/modprox/mp/pkg/upstream"
	"github.com/pkg/errors"

	// go-imports really hates these imports,
	// use aliases make them stick
	billyFS "gopkg.in/src-d/go-billy.v4"
	billyMemFS "gopkg.in/src-d/go-billy.v4/memfs"
	billyGit "gopkg.in/src-d/go-git.v4"
	billyMem "gopkg.in/src-d/go-git.v4/storage/memory"
)

type gitClient struct {
	options GitOptions
	fs      billyFS.Filesystem
	log     loggy.Logger
}

type GitOptions struct {
}

func NewGitClient(options GitOptions) Client {
	return &gitClient{
		fs: billyMemFS.New(),
	}
}

func (gc *gitClient) Protocols() []string {
	return []string{"git", "git+ssh"}
}

func (gc *gitClient) Get(r *upstream.Request) (repository.Blob, error) {
	if r == nil {
		return nil, errors.New("request is nil")
	}

	uri := r.URI()

	memoryStore := billyMem.NewStorage()
	if _, err := billyGit.Clone(memoryStore, gc.fs, &billyGit.CloneOptions{
		URL: uri,
	}); err != nil {
		return nil, err
	}

	// turn our in-memory filesystem into an in-memory zip archive
	return zipFS(gc.fs)
}

func zipFS(fs billyFS.Filesystem) (repository.Blob, error) {
	var buf bytes.Buffer
	zipper := zip.NewWriter(&buf)

	info, err := fs.Stat(".")
	if err != nil {
		return nil, err
	}

	if err := addFromFS(fs, zipper, ".", info); err != nil {
		return nil, err
	}

	if err := zipper.Flush(); err != nil {
		return nil, err
	}

	if err := zipper.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func addFromFS(fs billyFS.Filesystem, zipper *zip.Writer, p string, info os.FileInfo) error {
	fullName := filepath.Join(p, info.Name())
	isDir := info.IsDir()

	// fmt.Printf("addFromFS %s (isdir->%t)\n", fullName, isDir)

	if isDir {
		dirs, err := fs.ReadDir(fullName)
		if err != nil {
			return err
		}
		for _, dir := range dirs {
			if err := addFromFS(fs, zipper, fullName, dir); err != nil {
				return err
			}
		}
	} else {
		f, err := fs.Open(fullName)
		if err != nil {
			return err
		}
		w, err := zipper.Create(fullName)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, f); err != nil {
			return err
		}
	}
	return nil
}

/*
// Filesystem abstraction based on memory
fs := memfs.New()
// Git objects storer based on memory
storer := memory.NewStorage()

// Clones the repository into the worktree (fs) and storer all the .git
// content into the storer
_, err := git.Clone(storer, fs, &git.CloneOptions{
    URL: "https://github.com/git-fixtures/basic.git",
})
if err != nil {
    log.Fatal(err)
}

// Prints the content of the CHANGELOG file from the cloned repository
changelog, err := fs.Open("CHANGELOG")
if err != nil {
    log.Fatal(err)
}

io.Copy(os.Stdout, changelog)
*/
