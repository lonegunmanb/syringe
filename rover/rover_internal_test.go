package rover

import (
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/syringe/codegen"
	"github.com/lonegunmanb/syringe/ioc"
	. "github.com/lonegunmanb/syringe/util"
	"github.com/lonegunmanb/varys/ast"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"testing"
)

type roverInternalTestSuite struct {
	suite.Suite
	rover            *codeRover
	expectedTypeInfo *codegen.MockTypeInfo
}

func TestRoverInternalTestSuite(t *testing.T) {
	suite.Run(t, &roverInternalTestSuite{})
}

const roverTestStartingPath = "./"

func (suite *roverInternalTestSuite) SetupTest() {
	suite.rover = newCodeRover(roverTestStartingPath)
	ctrl := gomock.NewController(suite.T())
	suite.expectedTypeInfo = codegen.NewMockTypeInfo(ctrl)
	filePath := "filePath"
	fileName := "fileName"
	mockFileInfo := NewMockFileInfo(ctrl)
	mockFileInfo.EXPECT().Dir().AnyTimes().Return(filePath)
	mockFileInfo.EXPECT().Name().AnyTimes().Return(fileName)
	mockFileRetriever := NewMockFileRetriever(ctrl)
	mockFileRetriever.EXPECT().GetFiles(gomock.Eq(roverTestStartingPath)).Return([]ast.FileInfo{mockFileInfo}, nil)
	mockTypeWalker := NewMockTypeWalker(ctrl)
	mockTypeWalker.EXPECT().ParseDir(gomock.Eq(roverTestStartingPath), gomock.Eq("")).Times(1).Return(nil)
	mockTypeWalker.EXPECT().GetTypes().Times(1).Return([]ast.TypeInfo{suite.expectedTypeInfo})
	roverContainer.RegisterFactory((*ast.TypeWalker)(nil), func(ioc ioc.Container) interface{} {
		return mockTypeWalker
	})
}

func (*roverInternalTestSuite) TearDownTest() {
	roverContainer.Clear()
}

func (suite *roverInternalTestSuite) TestGetTypeInfos() {
	Given("a rover starting scan from ./", suite.T(), func() {
		When("get type infos the rover scanned", func() {
			typeInfos, err := suite.rover.getTypeInfos()
			Then("rover should return expected type info", func() {
				So(err, ShouldBeNil)
				And(typeInfos, ShouldHaveLength, 1)
				And(typeInfos[0], ShouldEqual, suite.expectedTypeInfo)
			})
		})
	})
}
