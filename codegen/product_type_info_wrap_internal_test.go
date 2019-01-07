package codegen

import (
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/syringe/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductTypeInfoWrap_GetFieldAssigns(t *testing.T) {
	ctrl, typeInfo := prepareTypeInfoMock(t)
	defer ctrl.Finish()
	mockFieldInfo := NewMockFieldInfo(ctrl)
	mockFieldInfo.EXPECT().GetTag().Times(1).Return("inject:\"\"")
	fieldInfos := []ast.FieldInfo{mockFieldInfo}
	typeInfo.EXPECT().GetFields().Times(1).Return(fieldInfos)
	sut := &typeInfoWrap{TypeInfo: typeInfo}
	fields := sut.GetFieldAssigns()
	assert.Equal(t, 1, len(fields))
	fieldWrap, ok := fields[0].(*productFieldInfoWrap)
	assert.True(t, ok)
	assert.Equal(t, mockFieldInfo, fieldWrap.FieldInfo)
	assert.Equal(t, sut, fieldWrap.typeInfo)
}

func TestProductTypeInfoWrap_FieldWithoutInjectTagShouldNotExport(t *testing.T) {
	ctrl, typeInfo := prepareTypeInfoMock(t)
	defer ctrl.Finish()
	mockFieldInfo := NewMockFieldInfo(ctrl)
	mockFieldInfo.EXPECT().GetTag().Times(1).Return("")
	fieldInfos := []ast.FieldInfo{mockFieldInfo}
	typeInfo.EXPECT().GetFields().Times(1).Return(fieldInfos)
	sut := &typeInfoWrap{TypeInfo: typeInfo}
	fields := sut.GetFieldAssigns()
	assert.Equal(t, 0, len(fields))
}

func TestProductTypeInfoWrap_GetEmbeddedTypeAssigns(t *testing.T) {
	ctrl, typeInfo := prepareTypeInfoMock(t)
	defer ctrl.Finish()
	mockEmbeddedTypes := NewMockEmbeddedType(ctrl)
	embeddedTypes := []ast.EmbeddedType{mockEmbeddedTypes}
	typeInfo.EXPECT().GetEmbeddedTypes().Times(1).Return(embeddedTypes)
	sut := &typeInfoWrap{TypeInfo: typeInfo}
	embeddedTypesGot := sut.GetEmbeddedTypeAssigns()
	assert.Equal(t, 1, len(embeddedTypesGot))
	embeddedTypeWrap, ok := embeddedTypesGot[0].(*productEmbeddedTypeWrap)
	assert.True(t, ok)
	assert.Equal(t, mockEmbeddedTypes, embeddedTypeWrap.EmbeddedType)
	assert.Equal(t, sut, embeddedTypeWrap.typeInfoWrap)
}

func TestProductTypeInfoWrap_GetPkgName(t *testing.T) {
	ctrl, mockTypeInfo := prepareTypeInfoMock(t)
	defer ctrl.Finish()
	expected := "expected"
	mockTypeInfo.EXPECT().GetPkgName().Times(1).Return(expected)
	sut := &typeInfoWrap{TypeInfo: mockTypeInfo}
	actual := sut.GetPkgName()
	assert.Equal(t, expected, actual)
}

func TestProductTypeInfoWrap_GetImportDecls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expected := []string{"a"}
	mockDepPkgPathInfo := NewMockDepPkgPathInfo(ctrl)
	mockDepPkgPathInfo.EXPECT().GenImportDecls().Times(1).Return(expected)
	sut := &typeInfoWrap{depPkgPathInfo: mockDepPkgPathInfo}
	imports := sut.GenImportDecls()
	assert.Equal(t, expected, imports)
}

func TestProductTypeInfoWrap_GetPkgNameFromPkgPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expected := "pkgName"
	mockDepPkgPathInfo := NewMockDepPkgPathInfo(ctrl)
	mockDepPkgPathInfo.EXPECT().GetPkgNameFromPkgPath("input").Times(1).Return(expected)
	sut := &typeInfoWrap{depPkgPathInfo: mockDepPkgPathInfo}
	imports := sut.GetPkgNameFromPkgPath("input")
	assert.Equal(t, expected, imports)
}

func TestProductTypeInfoWrap_GenRegisterCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDepPkgPathInfo := NewMockDepPkgPathInfo(ctrl)
	const pkgPath = "github.com/lonegunmanb/test_code/check_package_name_duplicate_a/model"
	mockDepPkgPathInfo.EXPECT().GetPkgNameFromPkgPath(pkgPath).Times(1).Return("p0")
	mockTypeInfo := NewMockTypeInfo(ctrl)
	mockTypeInfo.EXPECT().GetPkgPath().Times(2).Return(pkgPath)
	mockTypeInfo.EXPECT().GetName().Times(1).Return("Request")
	sut := NewTypeInfoWrapWithDepPkgPath(mockTypeInfo, mockDepPkgPathInfo)
	actual := sut.RegisterCode()
	const expected = "p0.Register_Request(container)"
	assert.Equal(t, expected, actual)
}

func prepareTypeInfoMock(t *testing.T) (*gomock.Controller, *MockTypeInfo) {
	ctrl := gomock.NewController(t)
	mockTypeInfo := NewMockTypeInfo(ctrl)
	return ctrl, mockTypeInfo
}
