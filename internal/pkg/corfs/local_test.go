package corfs

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalImplementsFileSystem(t *testing.T) {
	backend := LocalFilesystem{}
	var fileSystem FileSystem
	fileSystem = &backend

	assert.NotNil(t, fileSystem)
}

func TestLocalListFiles(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test")
	defer os.RemoveAll(tmpdir)
	assert.Nil(t, err)

	ioutil.WriteFile(path.Join(tmpdir, "tmpfile"), []byte("foo"), 0777)

	fs := LocalFilesystem{}

	files, err := fs.ListFiles(tmpdir)
	assert.Nil(t, err)

	assert.Len(t, files, 1)
	assert.Equal(t, "tmpfile", files[0].Name)
}

func TestLocalOpenReader(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test")
	defer os.RemoveAll(tmpdir)
	assert.Nil(t, err)

	ioutil.WriteFile(path.Join(tmpdir, "tmpfile"), []byte("foo bar baz"), 0777)

	fs := LocalFilesystem{}

	path := filepath.Join(tmpdir, "tmpfile")

	// Test reader that begins at beginning of file
	reader, err := fs.OpenReader(path, 0)
	assert.Nil(t, err)

	contents, err := ioutil.ReadAll(reader)
	assert.Nil(t, err)
	assert.Equal(t, []byte("foo bar baz"), contents)
	err = reader.Close()
	assert.Nil(t, err)

	// Test reader that begins in the middle of a file
	reader, err = fs.OpenReader(path, 4)
	assert.Nil(t, err)

	contents, err = ioutil.ReadAll(reader)
	assert.Nil(t, err)
	assert.Equal(t, []byte("bar baz"), contents)
	err = reader.Close()
	assert.Nil(t, err)
}

func TestLocalOpenWriter(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test")
	defer os.RemoveAll(tmpdir)
	assert.Nil(t, err)

	fs := LocalFilesystem{}

	path := filepath.Join(tmpdir, "tmpfile")

	writer, err := fs.OpenWriter(path)
	assert.Nil(t, err)

	n, err := writer.Write([]byte("foo bar baz"))
	assert.Equal(t, 11, n)
	assert.Nil(t, err)

	contents, err := ioutil.ReadFile(path)
	assert.Nil(t, err)
	assert.Equal(t, []byte("foo bar baz"), contents)
}

func TestLocalStat(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test")
	defer os.RemoveAll(tmpdir)
	assert.Nil(t, err)

	path := path.Join(tmpdir, "tmpfile")

	ioutil.WriteFile(path, []byte("foo"), 0777)

	fs := LocalFilesystem{}

	fInfo, err := fs.Stat(path)
	assert.Nil(t, err)

	assert.Equal(t, path, fInfo.Name)
	assert.Equal(t, int64(3), fInfo.Size)
}