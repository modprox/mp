package zips

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/repository"
)

const (
	goModFile     = "go.mod"
	hgArchiveFile = ".hg_archival.txt"
)

var log = loggy.New("zips")

// Rewrite the zip we downloaded from the upstream VCS into a zip file in the format
// required by the go/cmd tooling. Fundamentally the zip file the Go tooling requires
// re-namespacing the content under the directory path of the module.
//
// Additionally, the go/cmd does the following
// - limits the size of upstream .zip file to 500 MiB
// - limits the size of upstream LICENSE to 16 MiB
// - limits the size of upstream go.mod file to 16MiB
// -
// - removes other modules living in the same repo
// -
//
// The only complete "documentation" for the format of the new zip archive is in
// the go tool source code: go/src/cmd/go/internal/modfetch/coderepo.go
//
// The zip may contain multiple modules, each with its own go.mod.  We need to prune everything
// that isn't in our module.
func Rewrite(mod coordinates.Module, b repository.Blob) (repository.Blob, error) {
	in := bytes.NewReader(b)
	unZip, err := zip.NewReader(in, int64(len(b)))
	if err != nil {
		return nil, err
	}

	majorVersion, err := majorVersion(mod.Version)
	if err != nil {
		return nil, err
	}

	out := bytes.NewBuffer([]byte{})
	reZip := zip.NewWriter(out)

	// everything before the first / in the top-level entry of the zip file
	topPrefix := ""

	goModPath := make(map[string]string) // path => module specified in that directory's go.mod file, if that file exists

	// whether there's a LICENSE file in the top-level directory
	haveTopLicense := false

	// contents of LICENSE in top-level directory, if that file exists
	var topLicenseBytes []byte

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
			rc, err := zf.Open()
			if err != nil {
				return nil, err
			}
			goModBytes, err := ioutil.ReadAll(rc)
			if err != nil {
				return nil, errors.Wrapf(err, "error reading %s", zf.Name)
			}
			modulePath := ModulePath(goModBytes)
			if modulePath == "" {
				return nil, errors.Errorf("unable to extract module path from %s", zf.Name)
			}
			goModPath[dir] = modulePath
		} else if file == "LICENSE" && dir == topPrefix {
			rc, err := zf.Open()
			if err != nil {
				return nil, err
			}
			bytes, err := ioutil.ReadAll(rc)
			if err != nil {
				return nil, errors.Wrapf(err, "error reading %s", zf.Name)
			}
			topLicenseBytes = bytes
			haveTopLicense = true
		}

	}

	// If the requested module is in a subdirectory of the repo, we'll need to strip that subdirectory
	// name from each of the module's files.  Example: if the version is v2.0.4, strip the "v2/" prefix
	// from each of module's filenames.
	//
	// versionPrefix is that prefix.
	var versionPrefix string
	if majorVersion != "" && strings.HasSuffix(mod.Source, majorVersion) {
		versionPrefix = majorVersion + "/"
	}

	haveModLicense := false
	for _, zf := range unZip.File {
		if strings.HasSuffix(zf.Name, "/") {
			// drop directory dummy entries
			continue
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

		if len(goModPath) > 0 {
			if modPath := moduleOf(goModPath, zf.Name); modPath != mod.Source {
				continue
			}
		}

		base := path.Base(name)
		if strings.ToLower(base) == goModFile && base != goModFile {
			return nil, errors.Errorf("upstream zip file contains %s, want all lower-case go.mod", zf.Name)
		}

		if name == "LICENSE" {
			haveModLicense = true
		}

		rc, err := zf.Open()
		if err != nil {
			return nil, err
		}

		unversionedName := strings.TrimPrefix(name, versionPrefix)
		w, err := reZip.Create(mod.Source + "@" + mod.Version + "/" + unversionedName) // source@version/path
		if err != nil {
			return nil, err
		}

		if _, err := io.Copy(w, rc); err != nil {
			return nil, err
		}
	}

	// If the module doesn't have a LICENSE file but the top-level directory does, copy it to the module.
	if !haveModLicense && haveTopLicense {
		w, err := reZip.Create(mod.Source + "@" + mod.Version + "/LICENSE")
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(topLicenseBytes); err != nil {
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

// Given a mapping from dir path to module path, determine the module path for the specified file.  It's going to be the
// module for the longest dir path that prefixes this file's path.  You could probably use a clever data structure for this
// but there aren't likely to be many modules, so why bother?  If none match, return an empty string.
func moduleOf(goModPath map[string]string, filepath string) string {
	var longestDirPath string
	for dirpath := range goModPath {
		if strings.HasPrefix(filepath, dirpath) {
			if len(dirpath) > len(longestDirPath) {
				longestDirPath = dirpath
			}
		}
	}
	if longestDirPath == "" {
		return ""
	}

	return goModPath[longestDirPath]
}

func majorVersion(version string) (string, error) {
	if version[0] != 'v' {
		return "", errors.Errorf("version string does not begin with 'v': %s", version)
	}

	i := 1
	for i < len(version) && version[i] != '.' {
		i = i + 1
	}
	if i >= len(version) || version[i] != '.' {
		return "", errors.Errorf("version string does not contain a dot: %s", version)
	}

	m := version[:i]
	if m == "v0" || m == "v1" {
		return "", nil
	}

	return m, nil
}

// Copied from cmd/go/internal/modfile/read.go
var (
	slashSlash = []byte("//")
	moduleStr  = []byte("module")
)

// Copied from cmd/go/internal/modfile/read.go
//
// ModulePath returns the module path from the gomod file text.
// If it cannot find a module path, it returns an empty string.
// It is tolerant of unrelated problems in the go.mod file.
func ModulePath(mod []byte) string {
	for len(mod) > 0 {
		line := mod
		mod = nil
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line, mod = line[:i], line[i+1:]
		}
		if i := bytes.Index(line, slashSlash); i >= 0 {
			line = line[:i]
		}
		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, moduleStr) {
			continue
		}
		line = line[len(moduleStr):]
		n := len(line)
		line = bytes.TrimSpace(line)
		if len(line) == n || len(line) == 0 {
			continue
		}

		if line[0] == '"' || line[0] == '`' {
			p, err := strconv.Unquote(string(line))
			if err != nil {
				return "" // malformed quoted string or multiline module path
			}
			return p
		}

		return string(line)
	}
	return "" // missing module path
}

