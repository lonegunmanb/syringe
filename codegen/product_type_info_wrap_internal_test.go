package codegen

import (
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/syrinx/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductTypeInfoWrap_GetFields(t *testing.T) {
	ctrl, typeInfo := prepareTypeInfoMock(t)
	defer ctrl.Finish()
	mockFieldInfo := NewMockFieldInfo(ctrl)
	fieldInfos := []ast.FieldInfo{mockFieldInfo}
	typeInfo.EXPECT().GetFields().Times(1).Return(fieldInfos)
	sut := &typeInfoWrap{TypeInfo: typeInfo}
	fields := sut.GetFieldAssigns()
	assert.Equal(t, 1, len(fields))
	fieldWrap, ok := fields[0].(*productFieldInfoWrap)
	assert.True(t, ok)
	assert.Equal(t, mockFieldInfo, fieldWrap.FieldInfo)
}

func TestProductTypeInfoWrap_GetEmbeddedTypes(t *testing.T) {
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
}

func prepareTypeInfoMock(t *testing.T) (*gomock.Controller, *MockTypeInfo) {
	ctrl := gomock.NewController(t)
	mockTypeInfo := NewMockTypeInfo(ctrl)
	return ctrl, mockTypeInfo
}
