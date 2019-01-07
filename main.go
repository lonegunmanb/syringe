package main

import (
	"flag"
	"github.com/lonegunmanb/syringe/ast"
	"github.com/lonegunmanb/syringe/codegen"
	"github.com/lonegunmanb/syringe/rover"
	"io"
	"os"
)

func main() {
	clean := flag.Bool("c", false, "clean generated code")
	identName := flag.String("i", "container", "generated container identifier")
	flag.Parse()
	currentPath, err := os.Getwd()
	codegen.ContainerIdentName = *identName
	if err != nil {
		println(err.Error())
		return
	}

	if *clean {
		remove(currentPath)
	} else {
		create(currentPath)
	}
	//create("/Users/byers/go/src/github.com/lonegunmanb/blender")
}

func create(startingPath string) {
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
