package rover

import (
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/syringe/codegen"
	"github.com/lonegunmanb/syringe/ioc"
	"github.com/lonegunmanb/varys/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTypeInfos(t *testing.T) {
	roverStartingPath := "./"
	rover := newCodeRover(roverStartingPath)
	ctrl := gomock.NewController(t)
	filePath := "filePath"
	fileName := "fileName"
	mockFileInfo := NewMockFileInfo(ctrl)
	mockFileInfo.EXPECT().Path().AnyTimes().Return(filePath)
	mockFileInfo.EXPECT().Name().AnyTimes().Return(fileName)
	mockFileRetriever := NewMockFileRetriever(ctrl)
	mockFileRetriever.EXPECT().GetFiles(gomock.Eq(roverStartingPath)).Return([]ast.FileInfo{mockFileInfo}, nil)
	mockTypeInfo := codegen.NewMockTypeInfo(ctrl)
	mockTypeWalker := ast.NewMockTypeWalker(ctrl)
	mockTypeWalker.EXPECT().ParseDir(gomock.Eq(roverStartingPath), gomock.Eq("")).Times(1).Return(nil)
	mockTypeWalker.EXPECT().GetTypes().Times(1).Return([]ast.TypeInfo{mockTypeInfo})
	defer roverContainer.Clear()
	roverContainer.RegisterFactory((*ast.TypeWalker)(nil), func(ioc ioc.Container) interface{} {
		return mockTypeWalker
	})
	typeInfos, err := rover.getTypeInfos()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(typeInfos))
	assert.Equal(t, mockTypeInfo, typeInfos[0])
}
