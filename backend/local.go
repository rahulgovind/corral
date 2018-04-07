package backend

import (
	"io"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type LocalBackend struct {
	basePath string
}

func (l *LocalBackend) ListFiles() ([]FileInfo, error) {
	files := make([]FileInfo, 0)

	err := filepath.Walk(l.basePath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			log.Error(err)
			return err
		}
		if f.IsDir() {
			return nil
		}
		files = append(files, FileInfo{
			Name: f.Name(),
			Size: f.Size(),
		})
		return err
	})

	return files, err
}

func (l *LocalBackend) OpenReader(filePath string, startAt int64) (io.ReadCloser, error) {
	file, err := os.OpenFile(path.Join(l.basePath, filePath), os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	_, err = file.Seek(startAt, io.SeekStart)
	return file, err
}

func (l *LocalBackend) OpenWriter(filePath string) (io.WriteCloser, error) {
	filePath = path.Join(l.basePath, filePath)
	return os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
}

func (l *LocalBackend) Stat(filePath string) (FileInfo, error) {
	fInfo, err := os.Stat(path.Join(l.basePath, filePath))
	if err != nil {
		return FileInfo{}, err
	}
	return FileInfo{
		Name: fInfo.Name(),
		Size: fInfo.Size(),
	}, nil
}

func (l *LocalBackend) Init(basePath string) {
	l.basePath = basePath
}
