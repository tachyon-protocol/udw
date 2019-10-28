package udwGoWriter

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal"
	"path/filepath"
	"reflect"
)

type WriteGoFileGetterRequest struct {
	PackageName  string
	FunctionName string
	Obj          interface{}

	GoFilePath                  string
	BuildFlag                   string
	ObjTypeByteSliceHexEncoding bool
}

func WriteGoFileGetterWithGlobalVariable(req WriteGoFileGetterRequest) {
	if req.Obj == nil {
		panic("[WriteGoFileGetterWithGlobalVariable] req.Obj==nil")
	}
	if req.PackageName == "" {
		panic(`[WriteGoFileGetterWithGlobalVariable] req.PackageName==""`)
	}
	if req.FunctionName == "" {
		panic(`[WriteGoFileGetterWithGlobalVariable] req.FunctionName==""`)
	}
	if req.GoFilePath == "" {
		req.GoFilePath = filepath.Join("src", req.PackageName, req.FunctionName+".go")
	}
	goFile := NewGoFileContext(req.PackageName)
	goFile.SetBuildFlagContent(req.BuildFlag)
	typ := reflect.TypeOf(req.Obj)
	typStr := goFile.MustWriteTypeByReflect(typ)
	if udwGoTypeMarshal.IsObjHasInitAlloc(reflect.ValueOf(req.Obj), typ) {
		goFile.AddImportPath("sync")
		goFile.Buf.WriteString(`
var g` + req.FunctionName + ` ` + typStr + `
var g` + req.FunctionName + `Once sync.Once

func ` + req.FunctionName + `() ` + typStr + `{
	g` + req.FunctionName + `Once.Do(func(){
		g` + req.FunctionName + ` = `)
		goFile.MustWriteObjectToBuf(req.Obj, &goFile.Buf)
		goFile.Buf.WriteString(`
	})
	return g` + req.FunctionName + `
}
`)
	} else {
		goFile.Buf.WriteString(`var g` + req.FunctionName + ` ` + typStr + ` = `)
		hasWrite := false
		if req.ObjTypeByteSliceHexEncoding {
			v, ok := req.Obj.([]byte)
			if ok {
				goFile.Buf.WriteString(udwGoTypeMarshal.WriteByteSlice(v))
				hasWrite = true
			}
		}
		if hasWrite == false {
			goFile.MustWriteObjectToBuf(req.Obj, &goFile.Buf)
		}
		goFile.Buf.WriteString(`
func ` + req.FunctionName + `() ` + typStr + `{
	return g` + req.FunctionName + `
}
`)
	}
	goFile.MustWriteFileWithSelfBuffer(req.GoFilePath)
}

type WriteGoFileGetterWithDirFileListRequest struct {
	PackageName  string
	FunctionName string
	RootFileDir  string

	GoFilePath string
	BuildFlag  string
}

func WriteGoFileGetterWithDirFileList(req WriteGoFileGetterWithDirFileListRequest) {
	fileMap := map[string]string{}
	rootPath := udwFile.MustGetFullPath(req.RootFileDir)
	for _, fullFilePath := range udwFile.MustGetAllFiles(rootPath) {
		relFilePath := udwFile.MustGetRelativePath(rootPath, fullFilePath)
		fileMap[relFilePath] = string(udwFile.MustReadFile(fullFilePath))
	}
	WriteGoFileGetterWithGlobalVariable(WriteGoFileGetterRequest{
		PackageName:  req.PackageName,
		FunctionName: req.FunctionName,
		Obj:          fileMap,
		GoFilePath:   req.GoFilePath,
		BuildFlag:    req.BuildFlag,
	})
}
