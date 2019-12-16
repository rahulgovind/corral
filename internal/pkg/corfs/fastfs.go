package corfs

import (
	"fmt"
	"github.com/hashicorp/golang-lru"
	"github.com/rahulgovind/fastfs/helpers"
	"io"
	"net/url"
	"path/filepath"
	"strings"
)

//type FileSystem interface {
//	ListFiles(pathGlob string) ([]FileInfo, error)
//	Stat(filePath string) (FileInfo, error)
//	OpenReader(filePath string, startAt int64) (io.ReadCloser, error)
//	OpenWriter(filePath string) (io.WriteCloser, error)
//	Delete(filePath string) error
//	Join(elem ...string) string
//	Init() error
//}
var validFastFSSchemes = map[string]bool{
	"http":   true,
	"fastfs": true,
}

type FastFSFileSystem struct {
	objectCache *lru.Cache
	client      *helpers.Client
}

func parseFastFSURI(uri string) (*url.URL, error) {
	parsed, err := url.Parse(uri)
	if _, ok := validFastFSSchemes[parsed.Scheme]; !ok {

		return nil, fmt.Errorf("Invalid fastfs scheme: '%s'\tURI: %s", parsed.Scheme, uri)
	}

	if strings.HasPrefix(parsed.Path, "/") {
		parsed.Path = parsed.Path[1:]
	}

	//log.Infof("Parsed URL for %v: %v", uri, parsed)
	return parsed, err
}

func (ffs *FastFSFileSystem) ListFiles(pathGlob string) ([]FileInfo, error) {
	fastfsFiles := make([]FileInfo, 0)

	parsed, err := parseFastFSURI(pathGlob)

	if err != nil {
		return nil, err
	}

	dirGlob := pathGlob

	//if !strings.HasSuffix(pathGlob, "/") {
	//	dirGlob = pathGlob + "/*"
	//} else {
	//	dirGlob = pathGlob + "*"
	//}

	lastIndex := strings.LastIndex(pathGlob, "/")
	basePath := ""
	if lastIndex != -1 {
		basePath = pathGlob[:lastIndex+1]
	}

	parsedBasePath, err := parseFastFSURI(basePath)
	if err != nil {
		return nil, err
	}
	filelist, err := ffs.client.ListFiles(parsedBasePath.Path)

	objectPrefix := fmt.Sprintf("%s://%s/", parsed.Scheme, parsed.Host)
	for _, file := range filelist.Files {
		fullPath := objectPrefix + file.Path
		//fmt.Println(dirGlob, fullPath)
		match, _ := filepath.Match(dirGlob, fullPath)

		if match {
			ffs.objectCache.Add(fullPath, FileInfo{fullPath, file.Size})
			fastfsFiles = append(fastfsFiles, FileInfo{fullPath, file.Size})
		}
	}

	return fastfsFiles, nil
}

func (ffs *FastFSFileSystem) Stat(filePath string) (FileInfo, error) {
	if object, exists := ffs.objectCache.Get(filePath); exists {
		return object.(FileInfo), nil
	}

	parsed, err := parseFastFSURI(filePath)
	if err != nil {
		return FileInfo{}, err
	}

	fi, err := ffs.client.Stat(parsed.Path)
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{filePath, fi.Size}, nil
}

func (ffs *FastFSFileSystem) OpenReader(filePath string, startAt int64) (io.ReadCloser, error) {
	parsed, err := parseFastFSURI(filePath)
	if err != nil {
		return nil, err
	}

	r := ffs.client.OpenOffsetReader(parsed.Path, startAt)
	return r, nil
}

func (ffs *FastFSFileSystem) OpenWriter(filePath string) (io.WriteCloser, error) {
	parsed, err := parseFastFSURI(filePath)
	if err != nil {
		return nil, err
	}
	return ffs.client.BlockUploadWriter(parsed.Path)
	//return ffs.client.OpenWriter(parsed.Path)
}

func (ffs *FastFSFileSystem) Delete(filePath string) error {
	parsed, err := parseFastFSURI(filePath)
	if err != nil {
		return err
	}
	go ffs.client.Delete(parsed.Path)
	return nil
}

func (ffs *FastFSFileSystem) Join(elem ...string) string {
	stripped := make([]string, len(elem))
	for i, str := range elem {
		if strings.HasPrefix(str, "/") {
			str = str[1:]
		}
		if strings.HasSuffix(str, "/") && i != len(elem)-1 {
			str = str[:len(str)-1]
		}
		stripped[i] = str
	}

	return strings.Join(stripped, "/")
}

func (ffs *FastFSFileSystem) Init() error {
	ffs.client = helpers.New("ec2-18-220-77-181.us-east-2.compute.amazonaws.com:8100", 2, 5)
	ffs.objectCache, _ = lru.New(10000)
	return nil
}
