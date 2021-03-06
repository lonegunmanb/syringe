package rover

import (
	"fmt"
	"github.com/lonegunmanb/syringe/codegen"
	"github.com/lonegunmanb/syringe/ioc"
	"github.com/lonegunmanb/syringe/util"
	"github.com/lonegunmanb/varys/ast"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func GenerateCode(startingPath string, ignorePattern string, preferredPkgName *string) error {
	startingPath, err := toAbsPath(startingPath)
	if err != nil {
		return err
	}
	osEnv := getOsEnv()
	pkgPath, err := osEnv.GetPkgPath(startingPath)
	if err != nil {
		return err
	}
	rover := newCodeRover(startingPath)
	rover.ignorePattern = ignorePattern
	typeInfos, err := rover.getStructTypes()
	if err != nil {
		return err
	}
	fileOperator := getFileOperator()
	for _, typeInfo := range typeInfos {
		err := generateProductAssembleFile(typeInfo, osEnv, fileOperator)
		if err != nil {
			return err
		}
	}
	err = generateRegisterFile(startingPath, typeInfos, pkgPath, fileOperator, preferredPkgName)
	return err
}

func CleanGeneratedCodeFiles(startingPath string) error {
	r := newCodeRover(startingPath)
	fileRetriever := getFileRetriever()
	osEnv := getOsEnv()
	files, err := fileRetriever.GetFiles(r.roverStartingPath)
	fileOperator := getFileOperator()
	if err != nil {
		return err
	}
	for _, fileInfo := range files {
		if !isGoSrcFile(fileInfo) {
			continue
		}
		filePath := osEnv.ConcatFileNameWithPath(fileInfo.Dir(), fileInfo.Name())
		isGeneratedFile, err := isGeneratedFile(filePath, fileOperator)
		if err != nil {
			return err
		}
		if isGeneratedFile {
			err := fileOperator.Del(filePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var osEnvKey = (*ast.GoPathEnv)(nil)

func getOsEnv() ast.GoPathEnv {
	return roverContainer.GetOrRegister(osEnvKey, func(ioc ioc.Container) interface{} {
		return ast.NewGoPathEnv()
	}).(ast.GoPathEnv)
}

var fileRetrieverKey = (*ast.FileRetriever)(nil)

func generateProductAssembleFile(typeInfo ast.TypeInfo, osEnv ast.GoPathEnv, fileOperator util.FileOperator) error {
	fileName := fmt.Sprintf("%s.go", strings.ToLower(typeInfo.GetName()))
	createFileName := fmt.Sprintf("gen_%s", fileName)
	filePath := osEnv.ConcatFileNameWithPath(typeInfo.GetPhysicalPath(), createFileName)
	filePath, err := getUniqueFileName(filePath, createFileName, typeInfo, osEnv)
	if err != nil {
		return err
	}
	writer, err := fileOperator.Open(filePath)
	if err != nil {
		return err
	}
	err = writeHead(writer)
	if err != nil {
		return err
	}
	productCodegen := codegen.NewProductCodegen(typeInfo, writer)
	err = productCodegen.GenerateCode()
	if err != nil {
		return err
	}
	safeClose(writer)
	return nil
}

func generateRegisterFile(startingPath string,
	typeInfos []ast.TypeInfo,
	pkgPath string,
	fo util.FileOperator,
	preferPkgName *string) error {
	pkgName := getPkgName(pkgPath)
	if preferPkgName != nil {
		pkgName = *preferPkgName
	}
	for _, typeInfo := range typeInfos {
		if typeInfo.GetPhysicalPath() == startingPath && typeInfo.GetPkgName() != pkgName {
			pkgName = typeInfo.GetPkgName()
		}
	}
	registerFileName := fmt.Sprintf("%s/gen_register_ioc.go", startingPath)
	writer, err := fo.Open(registerFileName)
	if err != nil {
		return err
	}
	defer safeClose(writer)
	err = writeHead(writer)
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

func safeClose(writer io.Writer) {
	if c, ok := writer.(io.Closer); ok {
		_ = c.Close()
	}
}

func toAbsPath(startingPath string) (string, error) {
	if !filepath.IsAbs(startingPath) {
		absPath, err := filepath.Abs(startingPath)
		if err != nil {
			return "", err
		}
		return absPath, nil
	}
	return startingPath, nil
}

func getUniqueFileName(path string, fileName string, typeInfo ast.TypeInfo, osEnv ast.GoPathEnv) (string, error) {
	count := 1
	for {
		needRename, err := nonGeneratedFileExisted(path)
		if err != nil {
			return "", err
		}
		if !needRename {
			break
		}
		fileName = fmt.Sprintf("gen_%s_%d.go", strings.ToLower(typeInfo.GetName()), count)
		path = osEnv.ConcatFileNameWithPath(typeInfo.GetPhysicalPath(), fileName)
		count++
	}
	return path, nil
}

func getFileRetriever() ast.FileRetriever {
	return roverContainer.GetOrRegister(fileRetrieverKey, func(ioc ioc.Container) interface{} {
		return ast.NewFileRetriever()
	}).(ast.FileRetriever)
}

var fileOperatorKey = (*util.FileOperator)(nil)

func getFileOperator() util.FileOperator {
	return roverContainer.GetOrRegister(fileOperatorKey, func(ioc ioc.Container) interface{} {
		return &util.OsFileOperator{}
	}).(util.FileOperator)
}

func isGeneratedFile(filePath string, fileOperator util.FileOperator) (bool, error) {
	firstLine, err := fileOperator.FirstLine(filePath)
	return firstLine == commentHead, err
}

const commentHead = "// Code generated by syringe. DO NOT EDIT."

var firstLine = fmt.Sprintf("%s\r\n", commentHead)

func writeHead(writer io.Writer) error {
	_, err := writer.Write([]byte(firstLine))
	return err
}

func getPkgName(goPath string) string {
	s := strings.Split(goPath, "/")
	return s[len(s)-1]
}

func nonGeneratedFileExisted(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		isGenerated, err := isGeneratedFile(filePath, getFileOperator())
		return !isGenerated, err
	}
	return false, nil
}

func isGoSrcFile(info os.FileInfo) bool {
	return strings.HasSuffix(info.Name(), ".go")
}
