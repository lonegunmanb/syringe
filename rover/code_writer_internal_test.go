package rover

import (
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/syringe/ioc"
	"github.com/lonegunmanb/syringe/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCleanGeneratedCodeFiles(t *testing.T) {
	startingPath := "path"
	fileName := "gen_src.go"
	filePath := "path/gen_src.go"
	defer roverContainer.Clear()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFileRetriever := util.NewMockFileRetriever(ctrl)
	mockFileInfo := util.NewMockFileInfo(ctrl)
	mockFileInfo.EXPECT().Path().AnyTimes().Return(startingPath)
	mockFileInfo.EXPECT().Name().AnyTimes().Return(fileName)
	fileInfos := []util.FileInfo{mockFileInfo}
	mockFileRetriever.EXPECT().GetFiles(startingPath).Times(1).Return(fileInfos, nil)
	roverContainer.RegisterFactory((*util.FileRetriever)(nil), func(ioc ioc.Container) interface{} {
		return mockFileRetriever
	})
	mockOsEnv := util.NewMockGoPathEnv(ctrl)
	mockOsEnv.EXPECT().ConcatFileNameWithPath(startingPath, fileName).Times(1).Return(filePath)
	roverContainer.RegisterFactory((*util.GoPathEnv)(nil), func(ioc ioc.Container) interface{} {
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
	mockFileRetriever := util.NewMockFileRetriever(ctrl)
	mockFileInfo := util.NewMockFileInfo(ctrl)
	mockFileInfo.EXPECT().Path().AnyTimes().Return(startingPath)
	mockFileInfo.EXPECT().Name().AnyTimes().Return(fileName)
	fileInfos := []util.FileInfo{mockFileInfo}
	mockFileRetriever.EXPECT().GetFiles(startingPath).Times(1).Return(fileInfos, nil)
	roverContainer.RegisterFactory((*util.FileRetriever)(nil), func(ioc ioc.Container) interface{} {
		return mockFileRetriever
	})
	mockOsEnv := util.NewMockGoPathEnv(ctrl)
	mockOsEnv.EXPECT().ConcatFileNameWithPath(startingPath, fileName).AnyTimes().Return(filePath)
	roverContainer.RegisterFactory((*util.GoPathEnv)(nil), func(ioc ioc.Container) interface{} {
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
