package ast

import "os"

type FileInfo interface {
	os.FileInfo
	Path() string
}

type fileInfo struct {
	os.FileInfo
	filePath string
}

func (f *fileInfo) Path() string {
	return f.filePath
}
