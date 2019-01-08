package codegen

//go:generate mockgen -package=codegen -destination=./mock_register.go github.com/lonegunmanb/syringe/codegen Register
import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenRegisterPackageDecl(t *testing.T) {
	writer := &bytes.Buffer{}
	sut := &registerCodegen{pkgName: "ast", writer: writer}
	err := sut.genPkgDecl()
	assert.Nil(t, err)
	assert.Equal(t, expectedPackageDecl, writer.String())
}

func TestGenRegisterImportDecl(t *testing.T) {
	writer := &bytes.Buffer{}
	ctrl := gomock.NewController(t)
	mockDepPkgPathInfo := NewMockDepPkgPathInfo(ctrl)
	mockDepPkgPathInfo.EXPECT().GenImportDecls().Times(1).Return([]string{
		`"go/ast"`,
		`"go/token"`,
		`"go/types"`,
	})
	sut := &registerCodegen{writer: writer, depPkgPathInfo: mockDepPkgPathInfo}
	err := sut.genImportDecls()
	assert.Nil(t, err)
	assert.Equal(t, expectedImportDecl, writer.String())
}

func TestGenRegisterCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRegister := NewMockRegister(ctrl)
	mockRegister.EXPECT().RegisterCode().Times(1).Return("a.Register_a(container)")
	writer := &bytes.Buffer{}
	sut := &registerCodegen{writer: writer, typeInfos: []Register{mockRegister}}
	err := sut.genRegister()
	assert.Nil(t, err)
	const expected = `
func CreateIoc() ioc.Container {
    container := ioc.NewContainer()
    Register(container)
    return container
}
func Register(container ioc.Container) {
    a.Register_a(container)
}`
	assert.Equal(t, expected, writer.String())
}
