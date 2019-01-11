package codegen

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomContainerName(t *testing.T) {
	path := "path"
	name := "name"
	ctrl := gomock.NewController(t)
	mockTypeInfo := NewMockTypeInfoWrap(ctrl)
	mockTypeInfo.EXPECT().GetPkgPath().Times(1).Return(path)
	mockTypeInfo.EXPECT().GetName().Times(1).Return(name)
	sut := &registerCodeWriter{
		workingPkgPath: path,
		typeInfo:       mockTypeInfo,
	}
	originIdent := ContainerIdentName
	ContainerIdentName = "c"
	defer func() {
		ContainerIdentName = originIdent
	}()
	assert.Equal(t, "Register_name(c)", sut.RegisterCode())
}
