package zips

import (
	"archive/zip"
	"bytes"
	"io"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/modprox/mp/pkg/coordinates"
	"github.com/modprox/mp/pkg/repository"
)

const (
	goModFile     = "go.mod"
	hgArchiveFile = ".hg_archival.txt"
)

// Rewrite the zip we downloaded from the upstream VCS into a zip file in the format
// required by the go/cmd tooling. Fundamentally the zip file the Go tooling requires
// re-namespaces the content under the directory path of the module.
//
// Additionally, the go/cmd does the following
// - limits the size of upstream .zip file to 500 MiB
// - limits the size of upstream LICENSE to 16 MiB
// - limits the size of upstream go.mod file to 16MiB
// -
// - removes submodules
// -
//
// The only complete "documentation" for the format of the new zip archive is in
// the go tool source code: go/src/cmd/go/internal/modfetch/coderepo.go
func Rewrite(mod coordinates.Module, b repository.Blob) (repository.Blob, error) {
	in := bytes.NewReader(b)
	unZip, err := zip.NewReader(in, int64(len(b)))
	if err != nil {
		return nil, err
	}

	out := bytes.NewBuffer([]byte{})
	reZip := zip.NewWriter(out)

	topPrefix := ""
	hasGoMod := make(map[string]bool) // path => has go.mod file
	for _, zf := range unZip.File {
		if topPrefix == "" {
			i := strings.Index(zf.Name, "/")
			if i < 0 {
				return nil, errors.Errorf("upstream zip missing top-level directory prefix")
			}
			topPrefix = zf.Name[:i+1]
		}
		if !strings.HasPrefix(zf.Name, topPrefix) {
			return nil, errors.Errorf("upstream zip contains multiple top-level directories")
		}
		dir, file := path.Split(zf.Name)
		if file == goModFile {
			hasGoMod[dir] = true
		}
	}

	root := topPrefix
	inSubModule := func(name string) bool {
		for {
			dir, _ := path.Split(name)
			if len(dir) <= len(root) {
				return false
			}

			if hasGoMod[dir] {
				return true
			}

			name = dir[:len(dir)-1]
		}
	}

	for _, zf := range unZip.File {
		if topPrefix == "" {
			i := strings.Index(zf.Name, "/")
			if i < 0 {
				return nil, errors.Errorf("upstream missing top-level directory prefix")
			}
			topPrefix = zf.Name[:i+1]
		}

		if strings.HasSuffix(zf.Name, "/") {
			// drop directory dummy entries
			continue
		}

		if !strings.HasPrefix(zf.Name, topPrefix) {
			return nil, errors.Errorf("upstream zip file contains multiple top-level directories")
		}

		name := strings.TrimPrefix(zf.Name, topPrefix)
		if name == hgArchiveFile {
			// no hg stuff
			continue
		}

		if isVendorPath(name) {
			// no vendor directories
			continue
		}

		if inSubModule(zf.Name) {
			// no submodule directories
			continue
		}

		base := path.Base(name)
		if strings.ToLower(base) == goModFile && base != goModFile {
			return nil, errors.Errorf("upstream zip file contains %s, want all lower-case go.mod", zf.Name)
		}

		rc, err := zf.Open()
		if err != nil {
			return nil, err
		}

		w, err := reZip.Create(mod.Source + "@" + mod.Version + "/" + name) // source@version/path
		if err != nil {
			return nil, err
		}

		if _, err := io.Copy(w, rc); err != nil {
			return nil, err
		}
	}

	if err := reZip.Close(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func isVendorPath(name string) bool {
	var i int
	if strings.HasPrefix(name, "vendor/") {
		i += len("vendor/")
	} else if j := strings.Index(name, "/vendor/"); j >= 0 {
		i += len("/vendor/")
	} else {
		return false
	}
	return strings.Contains(name[i:], "/")
}
