package udwGoImport

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoBuildCtx"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal"
	"github.com/tachyon-protocol/udw/udwMap"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

type MustMulitGetPackageImportAllFilesRequest struct {
	AbsImportPathList []string

	IgnoreImportPackageFromGoRoot bool
	GoBuildCtx                    *udwGoBuildCtx.Ctx
}

func MustMulitGetPackageImportAllFiles(req MustMulitGetPackageImportAllFilesRequest) (resp MustMulitGetPackageImportResponse) {
	ctx := &ctxGoImport{
		req: req,
	}
	if ctx.req.GoBuildCtx == nil {
		ctx.req.GoBuildCtx = udwGoBuildCtx.NewCtxDefault()
	}
	ctx.resp.absImportPathMapMap = map[string]map[string]struct{}{}
	ctx.resp.absImportPathToDirMap = map[string]string{}
	ctx.seeImportPath = map[string]struct{}{}
	ctx.isDirMapCache = map[string]bool{}
	ctx.thisParentList = []string{}
	ctx.fset = token.NewFileSet()
	ctx.pathSearchList = ctx.req.GoBuildCtx.GetGoSearchPathList()
	for _, pkgPath := range req.AbsImportPathList {
		ctx.readOnePackageTopFn(pkgPath)
	}
	return ctx.resp
}

type ctxGoImport struct {
	req            MustMulitGetPackageImportAllFilesRequest
	resp           MustMulitGetPackageImportResponse
	thisParentList []string
	seeImportPath  map[string]struct{}
	fset           *token.FileSet
	pathSearchList []string
	isDirMapCache  map[string]bool
}

func (ctx *ctxGoImport) readOnePackageTopFn(AbsImportPath string) {
	ctx.thisParentList = ctx.thisParentList[:0]
	ctx.readOnePackageFn(AbsImportPath)
}

func (ctx *ctxGoImport) readOnePackageFn(AbsImportPath string) {
	ctx.thisParentList = append(ctx.thisParentList, AbsImportPath)
	defer func() {
		ctx.thisParentList = ctx.thisParentList[:len(ctx.thisParentList)-1]
	}()
	_, ok := ctx.seeImportPath[AbsImportPath]
	if ok {
		return
	}
	ctx.seeImportPath[AbsImportPath] = struct{}{}
	dir := ctx.searchAbsImportPathDir(AbsImportPath)
	if dir == "" {

		ctx.addErrMsg("vmhmrhg89y " + AbsImportPath)
		return
	}
	ctx.resp.absImportPathToDirMap[AbsImportPath] = dir
	ctx.resp.absImportPathMapMap[AbsImportPath] = map[string]struct{}{}
	importList := ctx.getImportListFromPkgDir(dir)
	for _, s := range importList {
		if s == "C" {
			continue
		}
		thisAbsImportPath := ctx.searchRealAbsImportPath(s, dir)
		if thisAbsImportPath == "" {
			ctx.addErrMsg("puyj7zhcgg " + s + " " + dir)
			continue
		}
		if thisAbsImportPath == AbsImportPath {

			continue
		}
		if ctx.req.IgnoreImportPackageFromGoRoot {
			thisDir := ctx.searchAbsImportPathDir(thisAbsImportPath)
			if strings.HasPrefix(thisDir, ctx.req.GoBuildCtx.GetGoRoot()) {
				continue
			}
		}
		ctx.resp.absImportPathMapMap[AbsImportPath][thisAbsImportPath] = struct{}{}
		ctx.readOnePackageFn(thisAbsImportPath)
	}
}

func (ctx *ctxGoImport) addErrMsg(errMsg string) {

}

func (ctx *ctxGoImport) searchAbsImportPathDir(AbsImportPath string) string {
	for _, gopath := range ctx.pathSearchList {
		thisPath := filepath.Join(gopath, "src", AbsImportPath)
		if ctx.isDir(thisPath) == true {
			return thisPath
		}
	}
	return ""
}

func (ctx *ctxGoImport) isDir(path string) bool {

	result, ok := ctx.isDirMapCache[path]
	if ok {
		return result
	}
	result = udwFile.MustIsDir(path)
	ctx.isDirMapCache[path] = result
	return result
}

func (ctx *ctxGoImport) searchRealAbsImportPath(ImportPath string, fromDir string) string {
	trySearchFn := func(root string) string {
		thisPath := filepath.Join(root, "src", ImportPath)
		if ctx.isDir(thisPath) == true {
			return ImportPath
		}
		thisFromDir := fromDir
		thisRootDir := filepath.Join(root, "src")
		for {
			if strings.HasPrefix(thisFromDir, thisRootDir) == false {
				return ""
			}
			thisVendorPath := filepath.Join(thisFromDir, "vendor")
			if ctx.isDir(thisVendorPath) {
				thisPath := filepath.Join(thisVendorPath, ImportPath)
				if ctx.isDir(thisPath) {
					outImportPath := strings.TrimPrefix(strings.TrimPrefix(thisPath, thisRootDir), "/")
					return outImportPath
				}
			}
			thisFromDir = filepath.Dir(thisFromDir)
		}
	}
	for _, gopath := range ctx.pathSearchList {
		out := trySearchFn(gopath)
		if out != "" {
			return out
		}
	}
	return ""
}

func (ctx *ctxGoImport) getImportListFromPkgDir(dir string) []string {
	importSet := map[string]struct{}{}
	allFileList := udwFile.MustGetAllFileOneLevel(dir)
	for _, file := range allFileList {
		ext := udwFile.GetExt(file)
		if ext != ".go" {
			continue
		}
		pf, err := parser.ParseFile(ctx.fset, file, nil, parser.ImportsOnly)
		if err != nil {

			continue
		}
		for _, spec := range pf.Imports {
			thisImportPath := importSpecGetImportPath(spec)
			importSet[thisImportPath] = struct{}{}
		}
	}
	return udwMap.SetStringToStringListAes(importSet)
}

func (ctx *ctxGoImport) mustAstDebugString(node interface{}) string {
	_buf := &bytes.Buffer{}
	err := ast.Fprint(_buf, ctx.fset, node, nil)
	if err != nil {
		panic(err)
	}
	return _buf.String()
}

func importSpecGetImportPath(spec *ast.ImportSpec) string {
	return udwGoTypeMarshal.MustReadGoStringFromString(spec.Path.Value)
}
