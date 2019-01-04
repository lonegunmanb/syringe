package main

import (
	"github.com/lonegunmanb/syrinx/ast"
	"github.com/lonegunmanb/syrinx/rover"
	"io"
	"os"
)

func main() {
	err := rover.GenerateCode("/Users/byers/go/src/github.com/lonegunmanb/syrinx/test_code", ast.NewGoPathEnv(), func(filePath string) (io.Writer, error) {
		return os.Create(filePath)
	})
	if err != nil {
		println(err.Error())
	}
}
