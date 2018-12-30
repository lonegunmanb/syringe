package codegen

//go:generate mockgen -package=codegen -destination=./mock_gen_task.go github.com/lonegunmanb/syrinx/codegen GenTask
import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenPackageDecl(t *testing.T) {
	testGen(t, func(mockTask *MockGenTask) {
		mockTask.EXPECT().GetPkgName().Times(1).Return("ast")
	}, func(gen *codegen) error {
		return gen.genPkgDecl()
	}, "package ast")
}

func testGen(t *testing.T, setupMockFunc func(info *MockGenTask),
	testMethod func(gen *codegen) error, expected string) {
	writer := &bytes.Buffer{}
	ctrl, task := prepareCodegenTaskMock(t)
	defer ctrl.Finish()
	//
	setupMockFunc(task)
	codegen := &codegen{writer: writer, genTask: task}
	err := testMethod(codegen)
	assert.Nil(t, err)
	code := writer.String()
	assert.Equal(t, expected, code)
}

func prepareCodegenTaskMock(t *testing.T) (*gomock.Controller, *MockGenTask) {
	ctrl := gomock.NewController(t)
	task := NewMockGenTask(ctrl)
	return ctrl, task
}
