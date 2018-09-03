package zips

import (
	"archive/zip"
	"bytes"
	"io"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/modprox/libmodprox/repository"
)

// we get to do some gymnastics to turn our zip archive
// into a properly formed zip archive.
//
// most of this code is inspired from go/src/cmd/go/internal/modfetch/coderepo.go

func Rewrite(mod repository.ModInfo, b repository.Blob) (repository.Blob, error) {
	in := bytes.NewReader(b)
	unzip, err := zip.NewReader(in, int64(len(b)))
	if err != nil {
		return nil, err
	}

	out := bytes.NewBuffer([]byte{})
	rezip := zip.NewWriter(out)

	topPrefix := ""

	for _, zf := range unzip.File {
		if topPrefix == "" {
			i := strings.Index(zf.Name, "/")
			if i < 0 {
				return nil, errors.Errorf("missing top-level directory prefix")
			}
			topPrefix = zf.Name[:i+1]
		}
		if !strings.HasPrefix(zf.Name, topPrefix) {
			return nil, errors.Errorf("zip file contains multiple top-level directories")
		}
	}

	for _, zf := range unzip.File {
		if topPrefix == "" {
			i := strings.Index(zf.Name, "/")
			if i < 0 {
				return nil, errors.Errorf("missing top-level directory prefix")
			}
			topPrefix = zf.Name[:i+1]
		}
		if strings.HasSuffix(zf.Name, "/") {
			// drop directory dummy entries
			continue
		}
		if !strings.HasPrefix(zf.Name, topPrefix) {
			return nil, errors.Errorf("zip file contains multiple top-level directories")
		}
		name := strings.TrimPrefix(zf.Name, topPrefix)
		if name == ".hg_archival.txt" {
			// no hg stuff
			continue
		}
		if name == ".gitattributes" {
			// no git attributes
			continue
		}
		if isVendoredPackage(name) {
			// no vendor directories
			continue
		}
		// todo: no submodules?
		base := path.Base(name)
		if strings.ToLower(base) == "go.mod" && base != "go.mod" {
			return nil, errors.Errorf("zip file contains %s, want all lower-case go.mod", zf.Name)
		}

		rc, err := zf.Open()
		if err != nil {
			return nil, err
		}
		w, err := rezip.Create(mod.Source + "@" + mod.Version + "/" + name) // source@version/path
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(w, rc); err != nil {
			return nil, err
		}
	}

	if err := rezip.Close(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func isVendoredPackage(name string) bool {
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
