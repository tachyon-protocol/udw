package udwGoImport

import (
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoBuildCtx"
	"github.com/tachyon-protocol/udw/udwProjectPath"
	"github.com/tachyon-protocol/udw/udwTest"
	"path/filepath"
	"testing"
)

const selfPkgPath = "github.com/tachyon-protocol/udw/udwGoSource/udwGoImport"

func TestMustMulitGetPackageImportResponse_GetAllIncludeDirList(t *testing.T) {
	return
	projectPath := udwProjectPath.MustGetProjectPath()
	resp := MustMulitGetPackageImport(MustMulitGetPackageImportRequest{
		AbsImportPathList: []string{
			selfPkgPath + "/kgiTest/kgiT1",
		},
		IgnoreImportPackageFromGoRoot: true,
	})
	dirList := resp.GetAllIncludeDirList()
	udwTest.Equal(len(dirList), 2)
	udwTest.Equal(dirList[0], filepath.Join(projectPath, "src/"+selfPkgPath+"/kgiTest/kgiT1"))
	udwTest.Equal(dirList[1], filepath.Join(projectPath, "src/"+selfPkgPath+"/kgiTest/kgiT1/kgiT1_1"))

	resp = MustMulitGetPackageImport(MustMulitGetPackageImportRequest{
		AbsImportPathList: []string{
			"github.com/tachyon-protocol/udw/udwGoSource/udwGoImport/kgiTest/kgiT1",
		},
		IgnoreImportPackageFromGoRoot: true,
		BuildCtx: udwGoBuildCtx.NewCtx(udwGoBuildCtx.CtxReq{
			TagList: []string{"ios"},
		}),
	})
	dirList = resp.GetAllIncludeDirList()
	udwTest.Equal(len(dirList), 3)
	udwTest.Equal(dirList[0], filepath.Join(projectPath, "src/"+selfPkgPath+"/kgiTest/kgiT1"))
	udwTest.Equal(dirList[1], filepath.Join(projectPath, "src/"+selfPkgPath+"/kgiTest/kgiT1/kgiT1_1"))
	udwTest.Equal(dirList[2], filepath.Join(projectPath, "src/"+selfPkgPath+"/kgiTest/kgiT1/kgiT1_2"))

	resp = MustMulitGetPackageImportAllFiles(MustMulitGetPackageImportAllFilesRequest{
		AbsImportPathList: []string{
			"github.com/tachyon-protocol/udw/udwGoSource/udwGoImport/kgiTest/kgiT1",
		},
		IgnoreImportPackageFromGoRoot: true,
	})
	importPathList := resp.GetAllIncludeImportPathList()
	udwTest.Equal(len(importPathList), 4)
	udwTest.Equal(importPathList[0], selfPkgPath+"/kgiTest/kgiT1")
	udwTest.Equal(importPathList[1], selfPkgPath+"/kgiTest/kgiT1/kgiT1_1")
	udwTest.Equal(importPathList[2], selfPkgPath+"/kgiTest/kgiT1/kgiT1_2")
	udwTest.Equal(importPathList[3], selfPkgPath+"/kgiTest/kgiT1/kgiT1_3")
}
