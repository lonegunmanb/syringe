package syrinx

///TODO:temp code, cleanup later
import (
	"errors"
	"github.com/ahmetb/go-linq"
	"go/ast"
	"go/token"
	"math/rand"
)

func Err() (int, error) {
	if rand.Int() > 1 {
		return 1, nil
	}
	return 0, errors.New("")
}

//noinspection GoUnusedExportedFunction
func GetAllStructNames(file *ast.File) []ast.Node {
	ast.Walk(nil, file)
	var tags = make([]ast.Node, 0)
	linq.From(file.Decls).Where(func(decl interface{}) bool {
		gd, ok := decl.(*ast.GenDecl)
		return ok && gd.Tok == token.TYPE
	}).SelectMany(func(gd interface{}) linq.Query {
		return linq.From(gd.(*ast.GenDecl).Specs)
	}).Where(func(spec interface{}) bool {
		ts, ok := spec.(*ast.TypeSpec)
		if !ok {
			return false
		}
		_, ok = ts.Type.(*ast.StructType)
		return ok
	}).Select(func(spec interface{}) interface{} {
		return spec.(*ast.TypeSpec)
	}).Select(func(ts interface{}) interface{} {
		return ts.(*ast.TypeSpec).Type
	}).SelectMany(func(structType interface{}) linq.Query {
		return linq.From(structType.(*ast.StructType).Fields.List)
	}).Select(func(field interface{}) interface{} {
		return field.(*ast.Field).Type
	}).ToSlice(&tags)
	return tags
}
