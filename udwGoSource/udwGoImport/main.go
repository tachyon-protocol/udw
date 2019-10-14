package udwGoImport

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoBuildCtx"
	"go/build"
	"path/filepath"
	"strings"
)

type MustMulitGetPackageImportRequest struct {
	AbsImportPathList []string

	IgnoreImportPackageFromGoRoot bool

	SimpleGoPathDirSearchNoError bool

	BuildCtx *udwGoBuildCtx.Ctx

	NotIgnoreImportC bool
}

func MustMulitGetPackageImport(req MustMulitGetPackageImportRequest) (resp MustMulitGetPackageImportResponse) {
	if req.BuildCtx == nil {
		req.BuildCtx = udwGoBuildCtx.NewCtxDefault()
	}
	GoBuildCtx := req.BuildCtx.ToGoBuildCtx()
	firstGoPath := req.BuildCtx.GetFirstGoPathString()
	importPathMapMap := map[string]map[string]struct{}{}
	absImportPathToDirMap := map[string]string{}
	seeImportPath := map[string]struct{}{}
	thisParentList := []string{}
	var readOnePackageFn func(AbsImportPath string)
	readOnePackageFn = func(AbsImportPath string) {
		thisParentList = append(thisParentList, AbsImportPath)
		defer func() {
			thisParentList = thisParentList[:len(thisParentList)-1]
		}()
		_, ok := seeImportPath[AbsImportPath]
		if ok {
			return
		}
		seeImportPath[AbsImportPath] = struct{}{}
		pkg, err := GoBuildCtx.Import(AbsImportPath, "", 0)
		if err != nil {
			if strings.Contains(err.Error(), "no buildable Go source files") {
				return
			}
			panic(err)
		}
		absImportPathToDirMap[AbsImportPath] = pkg.Dir
		importPathMapMap[AbsImportPath] = map[string]struct{}{}
		for _, s := range pkg.Imports {
			if s == "C" {
				if req.NotIgnoreImportC {
					importPathMapMap[AbsImportPath][s] = struct{}{}
				}
				continue
			}
			if req.SimpleGoPathDirSearchNoError {
				_, hasSeen := seeImportPath[s]
				if hasSeen == false && udwFile.MustIsDir(filepath.Join(firstGoPath, "src", s)) == false {
					continue
				}
				thisAbsImportPath := s
				importPathMapMap[AbsImportPath][thisAbsImportPath] = struct{}{}
				readOnePackageFn(thisAbsImportPath)
				continue
			}
			thisPkg, err := GoBuildCtx.Import(s, pkg.Dir, build.FindOnly)
			if err != nil {
				if strings.Contains(err.Error(), "no buildable Go source files") {
					continue
				}
				panic("ctxt.Import fail \n" + strings.Join(thisParentList, "\n") + "\n" + err.Error())
			}
			if req.IgnoreImportPackageFromGoRoot && strings.HasPrefix(thisPkg.Dir, req.BuildCtx.GetGoRoot()) == true {
				continue
			}
			thisAbsImportPath := thisPkg.ImportPath
			importPathMapMap[AbsImportPath][thisAbsImportPath] = struct{}{}
			readOnePackageFn(thisAbsImportPath)
		}
		for p1 := range importPathMapMap[AbsImportPath] {
			if p1 == "C" {
				continue
			}
			if absImportPathToDirMap[p1] == "" {
				panic("import " + p1 + " but not file from " + AbsImportPath)
			}
		}
	}
	readOnePackageTopFn := func(AbsImportPath string) {
		thisParentList = thisParentList[:0]
		readOnePackageFn(AbsImportPath)
	}
	for _, pkgPath := range req.AbsImportPathList {
		readOnePackageTopFn(pkgPath)
	}
	return MustMulitGetPackageImportResponse{
		absImportPathMapMap:   importPathMapMap,
		absImportPathToDirMap: absImportPathToDirMap,
	}
}
