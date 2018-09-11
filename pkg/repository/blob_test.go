package repository

import (
	"archive/zip"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func createFakeZip(t *testing.T, hasModFile bool) []byte {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"go.mod", "module github.com/modprox/libmodprox"},
		{"todo.txt", "Get animal handling licence.\nWrite more examples."},
	}
	for _, file := range files {
		if (file.Name != "go.mod") || hasModFile {
			f, err := w.Create(file.Name)
			require.NoError(t, err)
			_, err = f.Write([]byte(file.Body))
			require.NoError(t, err)
		}
	}

	err := w.Close()
	require.NoError(t, err)
	return buf.Bytes()
}

func Test_Blob_ModFile(t *testing.T) {
	b := createFakeZip(t, true)
	blob := Blob(b)

	content, exists, err := blob.ModFile()
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, "module github.com/modprox/libmodprox", content)
}

func Test_Blob_ModFile_none(t *testing.T) {
	b := createFakeZip(t, false)
	blob := Blob(b)

	_, exists, err := blob.ModFile()
	require.NoError(t, err)
	require.False(t, exists)
}
