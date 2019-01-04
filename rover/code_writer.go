package rover

import (
	"fmt"
	"github.com/lonegunmanb/syrinx/ast"
	"github.com/lonegunmanb/syrinx/codegen"
	"io"
	"path/filepath"
	"strings"
)

func GenerateCode(startingPath string, osEnv ast.GoPathEnv, writerFactory func(filePath string) (io.Writer, error)) error {
	if !filepath.IsAbs(startingPath) {
		absPath, err := filepath.Abs(startingPath)
		if err != nil {
			return err
		}
		startingPath = absPath
	}
	pkgPath, err := ast.GetPkgPath(osEnv, startingPath)
	pkgName := getPkgName(pkgPath)

	if err != nil {
		return err
	}
	rover := newCodeRover(startingPath)
	typeInfos, err := rover.getStructTypes()
	if err != nil {
		return err
	}
	for _, typeInfo := range typeInfos {
		createFileName := fmt.Sprintf("gen_%s.go", typeInfo.GetName())
		writer, err := writerFactory(osEnv.ConcatFileNameWithPath(typeInfo.GetPhysicalPath(), createFileName))
		if err != nil {
			return err
		}
		productCodegen := codegen.NewProductCodegen(typeInfo, writer)
		err = productCodegen.GenerateCode()
		if err != nil {
			return err
		}
	}
	registerFileName := fmt.Sprintf("%s/gen_register_ioc.go", startingPath)
	writer, err := writerFactory(registerFileName)
	if err != nil {
		return err
	}
	registerCodegen := codegen.NewRegisterCodegen(writer, typeInfos, pkgName, pkgPath)
	err = registerCodegen.GenerateCode()
	if err != nil {
		return err
	}
	return nil
}

func getPkgName(goPath string) string {
	s := strings.Split(goPath, "/")
	return s[len(s)-1]
}
