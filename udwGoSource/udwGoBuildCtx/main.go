package udwGoBuildCtx

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwProjectPath"
	"os"
	"path/filepath"
)

func GetDefaultGoPathString() string {
	if udwProjectPath.HasProjectPath() {
		return udwProjectPath.MustGetProjectPath()
	}
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return ""
	}
	goPathList := filepath.SplitList(gopath)
	return goPathList[0]
}

func GetDefaultGoPathList() []string {
	return []string{GetDefaultGoPathString()}
}

func GetDefaultGoRoot() string {
	goroot := os.Getenv("GOROOT")
	if goroot != "" {
		return goroot
	}
	if udwFile.MustFileExist("/usr/local/go") {
		return "/usr/local/go"
	}
	if udwFile.MustFileExist(`c:\go`) {
		return `c:\go`
	}
	return ""
}
