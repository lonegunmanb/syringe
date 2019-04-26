package rover

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/syringe/ioc"
	. "github.com/lonegunmanb/syringe/util"
	"github.com/lonegunmanb/varys/ast"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"go/types"
	"io"
	"reflect"
	"testing"
)

const testCleanGenFileStartingPath = "path"

type cleanGenFilesTestSuite struct {
	suite.Suite
	mockFileInfo     *MockFileInfo
	mockFileOperator *mockFileOperator
	fileName         string
	filePath         string
	firstLine        string
}

func TestRoverCleanGenFilesSuite(t *testing.T) {
	suite.Run(t, &cleanGenFilesTestSuite{})
}

func (suite *cleanGenFilesTestSuite) setupTest() {
	fileName := suite.fileName
	filePath := suite.filePath
	ctrl := gomock.NewController(suite.T())
	mockFileRetriever := NewMockFileRetriever(ctrl)
	suite.mockFileInfo = NewMockFileInfo(ctrl)
	suite.mockFileInfo.EXPECT().Dir().AnyTimes().Return(testCleanGenFileStartingPath)
	suite.mockFileInfo.EXPECT().Name().AnyTimes().Return(fileName)
	fileInfos := []ast.FileInfo{suite.mockFileInfo}
	mockFileRetriever.EXPECT().GetFiles(gomock.Eq(testCleanGenFileStartingPath)).Times(1).Return(fileInfos, nil)
	roverContainer.RegisterFactory((*ast.FileRetriever)(nil), func(ioc ioc.Container) interface{} {
		return mockFileRetriever
	})
	mockOsEnv := NewMockGoPathEnv(ctrl)
	mockOsEnv.EXPECT().ConcatFileNameWithPath(gomock.Eq(testCleanGenFileStartingPath),
		gomock.Eq(fileName)).AnyTimes().Return(filePath)
	roverContainer.RegisterFactory((*ast.GoPathEnv)(nil), func(ioc ioc.Container) interface{} {
		return mockOsEnv
	})
	suite.mockFileOperator = &mockFileOperator{}
	suite.mockFileOperator.On("FirstLine", filePath).Maybe().Return(suite.firstLine, nil)
	suite.mockFileOperator.On("Del", filePath).Maybe().Return(nil)
	roverContainer.RegisterFactory((*FileOperator)(nil), func(ioc ioc.Container) interface{} {
		return suite.mockFileOperator
	})
}

func (*cleanGenFilesTestSuite) TearDownTest() {
	defer roverContainer.Clear()
}

func (suite *cleanGenFilesTestSuite) TestCleanGeneratedCodeFiles() {
	Given("a generated go file with first line comment mark", suite.T(), func() {
		filePath := "path/gen_src.go"
		suite.filePath = filePath
		suite.fileName = "gen_src.go"
		suite.firstLine = commentHead
		suite.setupTest()
		When("invoke CleanGeneratedCodeFiles", func() {
			err := CleanGeneratedCodeFiles(testCleanGenFileStartingPath)
			Then("file should be deleted", func() {
				So(err, ShouldBeNil)
				And(suite, shouldDeleted, filePath)
			})
		})
	})
}

func (suite *cleanGenFilesTestSuite) TestCleanGeneratedCodeFilesWillNotDeleteNonGeneratedSrcFile() {
	Given("a non gen go file which name is like gen file but no comment mark", suite.T(), func() {
		filePath := "path/gen_src.go"
		suite.filePath = filePath
		suite.fileName = "gen_src.go"
		suite.firstLine = "package ast"
		suite.setupTest()
		When("invoke CleanGeneratedCodeFiles", func() {
			err := CleanGeneratedCodeFiles(testCleanGenFileStartingPath)
			Then("file should not be deleted", func() {
				So(err, ShouldBeNil)
				And(suite, shouldNotDeleted, filePath)
			})
		})
	})
	Given("a non gen cpp file which first line equal to comment head", suite.T(), func() {
		filePath := "path/gen_src.go"
		suite.filePath = filePath
		suite.fileName = "gen_src.cpp"
		suite.firstLine = commentHead
		suite.setupTest()
		When("invoke CleanGeneratedCodeFiles", func() {
			err := CleanGeneratedCodeFiles(testCleanGenFileStartingPath)
			Then("only go file with comment head will be deleted, not other src file", func() {
				So(err, ShouldBeNil)
				And(suite, shouldNotDeleted, filePath)
			})
		})
	})
}

type stubTypeInfo struct {
	pkgPath      string
	pkgName      string
	physicalPath string
}

func (*stubTypeInfo) GetName() string {
	return "StubType"
}

func (s *stubTypeInfo) GetPkgPath() string {
	return s.pkgPath
}

func (s *stubTypeInfo) GetPkgName() string {
	return s.pkgName
}

func (s *stubTypeInfo) GetPhysicalPath() string {
	return s.physicalPath
}

func (*stubTypeInfo) GetFields() []ast.FieldInfo {
	return []ast.FieldInfo{}
}

func (*stubTypeInfo) GetKind() reflect.Kind {
	panic("implement me")
}

func (*stubTypeInfo) GetType() types.Type {
	panic("implement me")
}

func (*stubTypeInfo) GetEmbeddedTypes() []ast.EmbeddedType {
	return []ast.EmbeddedType{}
}

func (*stubTypeInfo) GetFullName() string {
	return "StubType"
}

func (*stubTypeInfo) GetDepPkgPaths(fieldTagFilter string) []string {
	return []string{}
}

type stubFileOperator struct {
	writer *bytes.Buffer
}

func (s *stubFileOperator) Open(path string) (io.Writer, error) {
	return s.writer, nil
}

func (*stubFileOperator) Del(path string) error {
	panic("implement me")
}

func (*stubFileOperator) FirstLine(path string) (string, error) {
	panic("implement me")
}

func TestGenerateRegisterFileWithPkgNameInCodeDifferentWithPkgNameFromPath(t *testing.T) {
	Convey("given one existing type in a path with different pkg name with name resolved from path", t, func() {
		pkgPath := "github.com/lonegunmanb/syringe"
		startingPath := fmt.Sprintf("/go/src/%s", pkgPath)
		expectedName := "notsyringe"
		typeInfo := &stubTypeInfo{
			pkgPath:      pkgPath,
			pkgName:      expectedName,
			physicalPath: startingPath,
		}
		stubFo := &stubFileOperator{
			writer: new(bytes.Buffer),
		}
		Convey("when generate register file code", func() {
			err := generateRegisterFile(startingPath, []ast.TypeInfo{typeInfo}, pkgPath, stubFo, nil)
			Convey("then generated code should has same pkg name with existing type", func() {
				code := stubFo.writer.String()
				So(err, ShouldBeNil)
				So(code, ShouldContainSubstring, fmt.Sprintf("package %s", expectedName))
			})
		})
	})
}

func TestGenerateRegisterFileWithPreferedPkgName(t *testing.T) {
	Convey("given one prefered package name that different than name resolved from path", t, func() {
		pkgPath := "github.com/lonegunmanb/syringe"
		startingPath := fmt.Sprintf("/go/src/%s", pkgPath)
		preferredPkgName := "notsyringe"
		stubFo := &stubFileOperator{
			writer: new(bytes.Buffer),
		}
		Convey("when generate register file code", func() {
			err := generateRegisterFile(startingPath, []ast.TypeInfo{}, pkgPath, stubFo, &preferredPkgName)
			Convey("then generated code should has same pkg name with existing type", func() {
				code := stubFo.writer.String()
				So(err, ShouldBeNil)
				So(code, ShouldContainSubstring, fmt.Sprintf("package %s", preferredPkgName))
			})
		})
	})
}

func shouldDeleted(actual interface{}, expected ...interface{}) string {
	testSuite := actual.(*cleanGenFilesTestSuite)
	mockFileOperator := testSuite.mockFileOperator
	mockFileOperator.AssertCalled(testSuite.T(), "Del", expected[0].(string))
	return ""
}

func shouldNotDeleted(actual interface{}, expected ...interface{}) string {
	testSuite := actual.(*cleanGenFilesTestSuite)
	mockFileOperator := testSuite.mockFileOperator
	mockFileOperator.AssertNotCalled(testSuite.T(), "Del", expected[0].(string))
	return ""
}
