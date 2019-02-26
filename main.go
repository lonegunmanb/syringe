package main

import (
	"flag"
	"github.com/lonegunmanb/syringe/codegen"
	"github.com/lonegunmanb/syringe/rover"
	"os"
)

//go:generate go get github.com/golang/mock/gomock
//go:generate go install github.com/golang/mock/gomock
//go:generate go generate ./...
//go:generate go get -t -v ./...

func main() {
	clean := flag.Bool("c", false, "clean generated code")
	containerIdentName := flag.String("cident", "container", "generated container identifier")
	productIdentName := flag.String("pident", "product", "generated product identifier")
	ignorePattern := flag.String("ignore", "", "ignore file pattern")
	flag.Parse()
	currentPath, err := os.Getwd()
	codegen.ContainerIdentName = *containerIdentName
	codegen.ProductIdentName = *productIdentName
	if err != nil {
		println(err.Error())
		return
	}

	if *clean {
		remove(currentPath)
	} else {
		create(currentPath, *ignorePattern)
	}
}

func create(startingPath string, ignorePattern string) {
	err := rover.GenerateCode(startingPath, ignorePattern)
	if err != nil {
		println(err.Error())
	}
}

func remove(startingPath string) {
	err := rover.CleanGeneratedCodeFiles(startingPath)
	if err != nil {
		println(err.Error())
	}
}
