package codegen

import (
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/syringe/ast"
	"github.com/stretchr/testify/assert"
	"go-funk"
	"testing"
)

//go:generate mockgen -package=codegen -destination=./mock_dep_pkg_path_info.go github.com/lonegunmanb/syringe/codegen DepPkgPathInfo
func TestGetPkgNameFromPkgPath(t *testing.T) {
	cases := []*funk.Tuple{
		{"testing", "testing"},
		{"go/ast", "ast"},
		{"github.com/lonegunman/syringe/codegen", "codegen"},
	}
	for _, tuple := range cases {
		assert.Equal(t, tuple.Element2.(string), retrievePkgNameFromPkgPath(tuple.Element1.(string)))
	}
}

func TestGetDepPkgPathsWithPkgNameDuplicate(t *testing.T) {
	paths := []string{
		"a",
		"b/b",
		"c/b",
	}
	var typeInfos []ast.TypeInfo
	ctrl := gomock.NewController(t)
	for _, path := range paths {
		mockTypeInfo := NewMockTypeInfo(ctrl)
		mockTypeInfo.EXPECT().GetPkgPath().Times(1).Return("ast")
		mockTypeInfo.EXPECT().GetDepPkgPaths("inject").Times(1).Return([]string{path})
		typeInfos = append(typeInfos, mockTypeInfo)
	}
	sut := &depPkgPathInfo{
		typeInfos: typeInfos,
		mode:      ProductCodegenMode,
	}
	expectedPaths := []string{
		"ast",
		"a",
		"b/b",
		"c/b",
	}
	pathsReceived := sut.GetDepPkgPaths()
	assert.Equal(t, expectedPaths, pathsReceived)
	assert.Equal(t, "a", sut.GetPkgNameFromPkgPath("a"))
	assert.Equal(t, "p0", sut.GetPkgNameFromPkgPath("b/b"))
	assert.Equal(t, "p1", sut.GetPkgNameFromPkgPath("c/b"))
	imports := sut.GenImportDecls()
	expected := []string{
		`"ast"`,
		`"a"`,
		`p0 "b/b"`,
		`p1 "c/b"`,
	}
	assert.Equal(t, expected, imports)
}

func TestGetDepPkgPathsForRegister(t *testing.T) {
	paths := []string{
		"a",
		"b/b",
		"c/b",
	}
	var typeInfos []ast.TypeInfo
	ctrl := gomock.NewController(t)
	for _, path := range paths {
		mockTypeInfo := NewMockTypeInfo(ctrl)
		mockTypeInfo.EXPECT().GetPkgPath().Times(1).Return("ast")
		mockTypeInfo.EXPECT().GetDepPkgPaths("inject").Times(1).Return([]string{path})
		typeInfos = append(typeInfos, mockTypeInfo)
	}
	sut := &depPkgPathInfo{
		typeInfos: typeInfos,
		mode:      RegisterCodegenMode,
	}
	expectedPaths := []string{
		"ast",
	}
	pathsReceived := sut.GetDepPkgPaths()
	assert.Equal(t, expectedPaths, pathsReceived)
}

//TODO:random fail, investigate later
//func TestGetDepPkgPathsWithPkgNameDuplicateAndConflictWithGeneratedPackageName(t *testing.T) {
//	testDuplicateAndConflictPackageName(t, []string{
//		"p0",
//		"a/b",
//		"b/b",
//	}, []string{
//		`"ast"`,
//		`"p0"`,
//		`p1 "a/b"`,
//		`p2 "b/b"`,
//	})
//}

//func testDuplicateAndConflictPackageName(t *testing.T, depPkgPaths []string, expected []string) {
//	typeInfos := []ast.TypeInfo{}
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	for _, path := range depPkgPaths {
//		mockTypeInfo := NewMockTypeInfo(ctrl)
//		mockTypeInfo.EXPECT().GetPkgPath().Times(1).Return("ast")
//		mockTypeInfo.EXPECT().GetDepPkgPaths("inject").Times(1).Return([]string{path})
//		typeInfos = append(typeInfos, mockTypeInfo)
//	}
//	sut := &depPkgPathInfo{
//		typeInfos: typeInfos,
//		mode:      ProductCodegenMode,
//	}
//	imports := sut.GenImportDecls()
//	sort.Strings(expected)
//	sort.Strings(imports)
//	assert.Equal(t, expected, imports)
//	assert.Equal(t, expected, imports)
//}
