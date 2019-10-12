package udwGoParser

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwProjectPath"
	"github.com/tachyon-protocol/udw/udwTest"
	"path/filepath"
	"strings"
	"testing"
)

func TestMustParsePackage(ot *testing.T) {
	pkg := MustParsePackage(udwProjectPath.MustGetProjectPath(), "github.com/tachyon-protocol/udw/udwGoSource/udwGoParser/testPackage")
	udwTest.Equal(pkg.GetImportList(), []string{"bytes", "errors"})
}

func TestMustParsePackageFunc(ot *testing.T) {
	pkg := MustParsePackage(udwProjectPath.MustGetProjectPath(), "github.com/tachyon-protocol/udw/udwGoSource/udwGoParser/testPackage/testFunc")
	udwTest.Equal(len(pkg.funcList), 7)
}

func TestParseGoSrc(ot *testing.T) {
	gopath := "/usr/local/go"
	goSourcePath := filepath.Join(gopath, "src")
	dirList := udwFile.MustGetAllDir(goSourcePath)
	for _, dir := range dirList {
		if strings.Contains(dir, "testdata") {
			continue
		}
		dir, err := filepath.Rel(goSourcePath, dir)
		if err != nil {
			panic(err)
		}
		MustParsePackage(gopath, dir)
	}
}

func TestParseCurrentProject(ot *testing.T) {
	gopath := udwProjectPath.MustGetProjectPath()
	goSourcePath := filepath.Join(gopath, "src")
	dirList := udwFile.MustGetAllDir(goSourcePath)
	for _, dir := range dirList {
		if strings.Contains(dir, "go/loader/testdata") {
			continue
		}
		dir, err := filepath.Rel(goSourcePath, dir)
		if err != nil {
			panic(err)
		}
		MustParsePackage(gopath, dir)
	}
}
