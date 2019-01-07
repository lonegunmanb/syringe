package rover

//TODO:fix it after proper refactor
//func TestGetTypeInfos(t *testing.T) {
//	roverStartingPath := "./"
//	rover := newCodeRover(roverStartingPath)
//	ctrl := gomock.NewController(t)
//	filePath := "filePath"
//	fileName := "fileName"
//	mockFileInfo := NewMockFileInfo(ctrl)
//	mockFileInfo.EXPECT().Path().Times(1).Return(filePath)
//	mockFileInfo.EXPECT().Name().Times(1).Return(fileName)
//	mockFileRetriever := NewMockFileRetriever(ctrl)
//	mockFileRetriever.EXPECT().GetFiles(gomock.Eq(roverStartingPath), gomock.Any()).Return([]FileInfo{mockFileInfo}, nil)
//	mockTypeInfo := codegen.NewMockTypeInfo(ctrl)
//	mockTypeWalker := ast.NewMockTypeWalker(ctrl)
//	mockTypeWalker.EXPECT().ParseFile(gomock.Eq(filePath), gomock.Eq(fileName)).Times(1).Return(nil)
//	mockTypeWalker.EXPECT().GetTypes().Times(1).Return([]ast.TypeInfo{mockTypeInfo})
//	mockTypeWalkerFactory := func() ast.TypeWalker {
//		return mockTypeWalker
//	}
//	rover.walkerFactory = mockTypeWalkerFactory
//	rover.fileRetriever = mockFileRetriever
//	typeInfos, err := rover.getTypeInfos()
//	assert.Nil(t, err)
//	assert.Equal(t, 1, len(typeInfos))
//	assert.Equal(t, mockTypeInfo, typeInfos[0])
//}
