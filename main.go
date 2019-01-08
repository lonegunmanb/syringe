package main

import (
	"flag"
	"github.com/lonegunmanb/syringe/codegen"
	"github.com/lonegunmanb/syringe/rover"
	"os"
)

func main() {
	clean := flag.Bool("c", false, "clean generated code")
	containerIdentName := flag.String("cident", "container", "generated container identifier")
	productIdentName := flag.String("pident", "product", "generated product identifier")
	ignorePatten := flag.String("ignore", "", "ignore file patten")
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
		create(currentPath, *ignorePatten)
	}
}

func create(startingPath string, ignorePatten string) {
	err := rover.GenerateCode(startingPath, ignorePatten)
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
