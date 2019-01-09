package ast

import (
	"errors"
	"fmt"
	"github.com/ahmetb/go-linq"
	"os"
	"runtime"
	"strings"
)

type GoPathEnv interface {
	IsWindows() bool
	GetGoPath() string
	ConcatFileNameWithPath(path string, fileName string) string
}

type envImpl struct {
}

func NewGoPathEnv() GoPathEnv {
	return &envImpl{}
}

func (*envImpl) IsWindows() bool {
	return runtime.GOOS == "windows"
}

func (*envImpl) GetGoPath() string {
	return os.Getenv("GOPATH")
}

func (e *envImpl) ConcatFileNameWithPath(path string, fileName string) string {
	return concatFileNameWithPath(e.IsWindows(), path, fileName)
}

func concatFileNameWithPath(isWindows bool, path string, fileName string) string {
	if isWindows {
		return fmt.Sprintf("%s\\%s", path, fileName)
	}
	return fmt.Sprintf("%s/%s", path, fileName)
}

func GetPkgPath(env GoPathEnv, systemPath string) (string, error) {
	isWindows := env.IsWindows()
	goPaths, err := getGoPaths(env)
	if err != nil {
		return "", err
	}
	return getPkgPathFromSystemPathUsingGoPath(isWindows, goPaths, systemPath)
}

func getGoPaths(env GoPathEnv) (gopaths []string, err error) {
	sep := ":"
	if env.IsWindows() {
		sep = ";"
	}
	goPath := env.GetGoPath()
	if goPath == "" {
		err = errors.New("no go path detected")
		return
	}
	gopaths = strings.Split(goPath, sep)
	return
}

func getPkgPathFromSystemPathUsingGoPath(isWindows bool, goPaths []string, systemPath string) (pkgPath string, err error) {
	goSrcPath := linq.From(goPaths).Select(func(path interface{}) interface{} {
		srcTemplate := "%s/src/"
		if isWindows {
			srcTemplate = "%s\\src\\"
		}
		return fmt.Sprintf(srcTemplate, path)
	}).FirstWith(func(path interface{}) bool {
		return strings.HasPrefix(systemPath, path.(string))
	})
	if goSrcPath == nil {
		err = errors.New(fmt.Sprintf("%s is not in go src path", systemPath))
		return
	} else {
		pkgPath = strings.TrimPrefix(systemPath, goSrcPath.(string))
		if isWindows {
			pkgPath = strings.Replace(pkgPath, "\\", "/", -1)
		}
	}
	return
}
