package repository

import (
	"archive/zip"
	"bytes"
	"io/ioutil"

	"github.com/pkg/errors"
)

// A Blob is an in-memory zip archive, representative of
// a repository that was downloaded from upstream.
//
// There might not be a go.mod file.
type Blob []byte

func (b Blob) ModFile() (string, bool, error) {
	r := bytes.NewReader(b)
	unzip, err := zip.NewReader(r, int64(len(b)))
	if err != nil {
		return "", false, errors.Wrap(err, "failed to open blob")
	}

	for _, f := range unzip.File {
		if f.Name == "go.mod" {
			rc, err := f.Open()
			if err != nil {
				return "", false, err
			}

			bs, err := ioutil.ReadAll(rc)
			if err != nil {
				return "", false, errors.Wrap(err, "failed to read go.mod file from blob")
			}

			return string(bs), true, rc.Close()
		}
	}

	return "", false, nil
}
