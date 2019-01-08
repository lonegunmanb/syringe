package util

import (
	"bufio"
	"io"
	"os"
)

type FileOperator interface {
	Open(path string) (io.Writer, error)
	Del(path string) error
	FirstLine(path string) (string, error)
}

type OsFileOperator struct {
}

func (*OsFileOperator) Open(path string) (io.Writer, error) {
	return os.Create(path)
}

func (*OsFileOperator) Del(path string) error {
	return os.Remove(path)
}

func (*OsFileOperator) FirstLine(path string) (string, error) {
	fi, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = fi.Close()
	}()

	br := bufio.NewReader(fi)
	content, _, _ := br.ReadLine()
	firstLine := string(content)
	return firstLine, nil
}
