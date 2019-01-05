package ast

//go:generate mockgen -package=ast -destination=./mock_gopathenv.go github.com/lonegunmanb/syringe/ast GoPathEnv
import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

const expectedPkgName = "github.com/lonegunmanb/syringe"

func TestGetPkgPathFromSystemPathUsingGoPath(t *testing.T) {

	testGetPkgPathFromSystemPath(t, []string{
		"/Users/user/go",
	},
		"/Users/user/go/src/github.com/lonegunmanb/syringe",
		expectedPkgName)
	testGetPkgPathFromSystemPath(t, []string{
		"/Users/user2/go",
		"/Users/user/go",
	},
		"/Users/user/go/src/github.com/lonegunmanb/syringe",
		expectedPkgName)
}

func TestGetPkgPathInWindows(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockEnv := NewMockGoPathEnv(ctrl)
	mockEnv.EXPECT().IsWindows().Times(2).Return(true)
	mockEnv.EXPECT().GetGoPath().Times(1).Return("c:\\go")
	pkgPath, err := GetPkgPath(mockEnv, "c:\\go\\src\\github.com\\lonegunmanb\\syringe")
	assert.Nil(t, err)
	assert.Equal(t, expectedPkgName, pkgPath)
}

func TestConcatFileNameWithPath(t *testing.T) {
	path := concatFileNameWithPath(false, "/Users/user/go", "file")
	assert.Equal(t, "/Users/user/go/file", path)
	path = concatFileNameWithPath(true, "c:\\go", "file")
	assert.Equal(t, "c:\\go\\file", path)
}

func testGetPkgPathFromSystemPath(t *testing.T, goPaths []string, systemPath string, expected string) {
	pkgPath, err := getPkgPathFromSystemPathUsingGoPath(false, goPaths, systemPath)
	assert.Nil(t, err)
	assert.Equal(t, expected, pkgPath)
}
