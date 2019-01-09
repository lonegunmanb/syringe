package ast

import (
	"os"
	"path/filepath"
)

type FileRetriever interface {
	GetFiles(startingPath string) ([]FileInfo, error)
}

func NewFileRetriever() FileRetriever {
	return &fileRetriever{}
}

type fileRetriever struct {
}

func (r *fileRetriever) GetFiles(startingPath string) ([]FileInfo, error) {
	var files []FileInfo
	err := filepath.Walk(startingPath, func(path string, info os.FileInfo, err error) error {
		files = append(files, &fileInfo{
			FileInfo: info,
			filePath: filepath.Dir(path),
		})
		return nil
	})
	return files, err
}
