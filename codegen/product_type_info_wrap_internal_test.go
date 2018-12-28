package codegen

import (
	"github.com/lonegunmanb/syrinx/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductTypeInfoWrap_GetFields(t *testing.T) {
	ctrl, typeInfo := prepareMock(t)
	defer ctrl.Finish()
	mockFieldInfo := NewMockFieldInfo(ctrl)
	fieldInfos := []ast.FieldInfo{mockFieldInfo}
	typeInfo.EXPECT().GetFields().Times(1).Return(fieldInfos)
	sut := &productTypeInfoWrap{TypeInfo: typeInfo}
	fields := sut.GetFields()
	assert.Equal(t, 1, len(fields))
	fieldWrap, ok := fields[0].(*productFieldInfoWrap)
	assert.True(t, ok)
	assert.Equal(t, mockFieldInfo, fieldWrap.FieldInfo)
}
