package util

import (
	"os"
	"path/filepath"
)

type FileRetriever interface {
	GetFiles(startingPath string, predicate func(os.FileInfo) bool) ([]FileInfo, error)
}

func NewFileRetriever() FileRetriever {
	return &fileRetriever{}
}

type fileRetriever struct {
}

func (r *fileRetriever) GetFiles(startingPath string, predicate func(os.FileInfo) bool) ([]FileInfo, error) {
	var files []FileInfo
	err := filepath.Walk(startingPath, func(path string, info os.FileInfo, err error) error {
		if predicate(info) {
			files = append(files, &fileInfo{
				FileInfo: info,
				filePath: filepath.Dir(path),
			})
		}
		return nil
	})
	return files, err
}
