package rover

import (
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/syringe/ioc"
	"github.com/lonegunmanb/syringe/util"
	"github.com/lonegunmanb/varys/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:generate mockgen -package=rover -destination=./mock_gopathenv.go github.com/lonegunmanb/varys/ast GoPathEnv
//go:generate mockgen -package=rover -destination=./mock_file_retriever.go github.com/lonegunmanb/varys/ast FileRetriever
//go:generate mockgen -package=rover -destination=./mock_file_info.go github.com/lonegunmanb/varys/ast FileInfo
func TestCleanGeneratedCodeFiles(t *testing.T) {
	startingPath := "path"
	fileName := "gen_src.go"
	filePath := "path/gen_src.go"
	defer roverContainer.Clear()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFileRetriever := NewMockFileRetriever(ctrl)
	mockFileInfo := NewMockFileInfo(ctrl)
	mockFileInfo.EXPECT().Dir().AnyTimes().Return(startingPath)
	mockFileInfo.EXPECT().Name().AnyTimes().Return(fileName)
	fileInfos := []ast.FileInfo{mockFileInfo}
	mockFileRetriever.EXPECT().GetFiles(startingPath).Times(1).Return(fileInfos, nil)
	roverContainer.RegisterFactory((*ast.FileRetriever)(nil), func(ioc ioc.Container) interface{} {
		return mockFileRetriever
	})
	mockOsEnv := NewMockGoPathEnv(ctrl)
	mockOsEnv.EXPECT().ConcatFileNameWithPath(startingPath, fileName).Times(1).Return(filePath)
	roverContainer.RegisterFactory((*ast.GoPathEnv)(nil), func(ioc ioc.Container) interface{} {
		return mockOsEnv
	})
	mockFileOperator := util.NewMockFileOperator(ctrl)
	mockFileOperator.EXPECT().FirstLine(filePath).Times(1).Return(commentHead, nil)
	mockFileOperator.EXPECT().Del(filePath).Times(1).Return(nil)
	roverContainer.RegisterFactory((*util.FileOperator)(nil), func(ioc ioc.Container) interface{} {
		return mockFileOperator
	})
	err := CleanGeneratedCodeFiles(startingPath)
	assert.Nil(t, err)
}

func TestCleanGeneratedCodeFilesWillNotTouchNonGeneratedSrcFile(t *testing.T) {
	testNotTouchNonGeneratedFile(t, "path", "gen_src.go", "path/gen_src.go", "package ast")
	testNotTouchNonGeneratedFile(t, "path", "gen_src.cpp", "path/gen_src.cpp", commentHead)
}

func testNotTouchNonGeneratedFile(t *testing.T, startingPath string, fileName string, filePath string, firstLine string) {
	defer roverContainer.Clear()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFileRetriever := NewMockFileRetriever(ctrl)
	mockFileInfo := NewMockFileInfo(ctrl)
	mockFileInfo.EXPECT().Dir().AnyTimes().Return(startingPath)
	mockFileInfo.EXPECT().Name().AnyTimes().Return(fileName)
	fileInfos := []ast.FileInfo{mockFileInfo}
	mockFileRetriever.EXPECT().GetFiles(startingPath).Times(1).Return(fileInfos, nil)
	roverContainer.RegisterFactory((*ast.FileRetriever)(nil), func(ioc ioc.Container) interface{} {
		return mockFileRetriever
	})
	mockOsEnv := NewMockGoPathEnv(ctrl)
	mockOsEnv.EXPECT().ConcatFileNameWithPath(startingPath, fileName).AnyTimes().Return(filePath)
	roverContainer.RegisterFactory((*ast.GoPathEnv)(nil), func(ioc ioc.Container) interface{} {
		return mockOsEnv
	})
	mockFileOperator := util.NewMockFileOperator(ctrl)
	mockFileOperator.EXPECT().FirstLine(filePath).AnyTimes().Return(firstLine, nil)
	mockFileOperator.EXPECT().Del(filePath).Times(0).Return(nil)
	roverContainer.RegisterFactory((*util.FileOperator)(nil), func(ioc ioc.Container) interface{} {
		return mockFileOperator
	})
	err := CleanGeneratedCodeFiles(startingPath)
	assert.Nil(t, err)
}
