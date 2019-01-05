package main

import (
	"flag"
	"github.com/lonegunmanb/syrinx/ast"
	"github.com/lonegunmanb/syrinx/rover"
	"io"
	"os"
)

func main() {
	clean := flag.Bool("clean", false, "clean generated code")
	flag.Parse()
	currentPath, err := os.Getwd()
	if err != nil {
		println(err.Error())
		return
	}
	println(currentPath)

	if *clean {
		remove(currentPath)
	} else {
		create(currentPath)
	}
}

func create(startingPath string) {
	println(startingPath)
	err := rover.GenerateCode(startingPath, ast.NewGoPathEnv(), func(filePath string) (io.Writer, error) {
		return os.Create(filePath)
	})
	if err != nil {
		println(err.Error())
	}
}

func remove(startingPath string) {
	err := rover.CleanGeneratedCodeFiles(startingPath, ast.NewGoPathEnv())
	if err != nil {
		println(err.Error())
	}
}
