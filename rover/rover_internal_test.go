package rover

import (
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/syrinx/ast"
	"github.com/lonegunmanb/syrinx/codegen"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:generate mockgen -package=rover -destination=./mock_file_retriever.go github.com/lonegunmanb/syrinx/rover FileRetriever
//go:generate mockgen -package=rover -destination=./mock_file_info.go github.com/lonegunmanb/syrinx/rover FileInfo

func TestGetTypeInfos(t *testing.T) {
	roverStartingPath := "./"
	rover := newCodeRover(roverStartingPath)
	ctrl := gomock.NewController(t)
	filePath := "filePath"
	fileName := "fileName"
	mockFileInfo := NewMockFileInfo(ctrl)
	mockFileInfo.EXPECT().Path().Times(1).Return(filePath)
	mockFileInfo.EXPECT().Name().Times(1).Return(fileName)
	mockFileRetriever := NewMockFileRetriever(ctrl)
	mockFileRetriever.EXPECT().GetFiles(gomock.Eq(roverStartingPath), gomock.Any()).Return([]FileInfo{mockFileInfo}, nil)
	mockTypeInfo := codegen.NewMockTypeInfo(ctrl)
	mockTypeWalker := ast.NewMockTypeWalker(ctrl)
	mockTypeWalker.EXPECT().ParseFile(gomock.Eq(filePath), gomock.Eq(fileName)).Times(1).Return(nil)
	mockTypeWalker.EXPECT().GetTypes().Times(1).Return([]ast.TypeInfo{mockTypeInfo})
	mockTypeWalkerFactory := func() ast.TypeWalker {
		return mockTypeWalker
	}
	rover.walkerFactory = mockTypeWalkerFactory
	rover.fileRetriever = mockFileRetriever
	typeInfos, err := rover.getTypeInfos()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(typeInfos))
	assert.Equal(t, mockTypeInfo, typeInfos[0])
}

func TestIsNotTestFile(t *testing.T) {
	assert.True(t, isTestFile("rover_test.go"))
	assert.True(t, isTestFile("test.go"))
	assert.False(t, isTestFile("rover.go"))
}

func TestIsGoFile(t *testing.T) {
	assert.True(t, isGoSrcFile("src.go"))
	assert.False(t, isGoSrcFile("src.cpp"))
	assert.False(t, isGoSrcFile("go"))
}
