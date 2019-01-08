package rover

import (
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/syringe/ast"
	"github.com/lonegunmanb/syringe/codegen"
	"github.com/lonegunmanb/syringe/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTypeInfos(t *testing.T) {
	roverStartingPath := "./"
	rover := newCodeRover(roverStartingPath)
	ctrl := gomock.NewController(t)
	filePath := "filePath"
	fileName := "fileName"
	mockFileInfo := util.NewMockFileInfo(ctrl)
	mockFileInfo.EXPECT().Path().AnyTimes().Return(filePath)
	mockFileInfo.EXPECT().Name().AnyTimes().Return(fileName)
	mockFileRetriever := util.NewMockFileRetriever(ctrl)
	mockFileRetriever.EXPECT().GetFiles(gomock.Eq(roverStartingPath)).Return([]util.FileInfo{mockFileInfo}, nil)
	mockTypeInfo := codegen.NewMockTypeInfo(ctrl)
	mockTypeWalker := ast.NewMockTypeWalker(ctrl)
	mockTypeWalker.EXPECT().ParseDir(gomock.Eq(roverStartingPath), gomock.Eq("")).Times(1).Return(nil)
	mockTypeWalker.EXPECT().GetTypes().Times(1).Return([]ast.TypeInfo{mockTypeInfo})
	mockTypeWalkerFactory := func() ast.TypeWalker {
		return mockTypeWalker
	}
	rover.walkerFactory = mockTypeWalkerFactory
	typeInfos, err := rover.getTypeInfos()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(typeInfos))
	assert.Equal(t, mockTypeInfo, typeInfos[0])
}
