package udwGoWriter

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoParser"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal"
	"github.com/tachyon-protocol/udw/udwMap"
	"go/types"
	"path"
	"strings"
)

type GoFileContext struct {
	packageName      string
	pkgImportPath    string
	pkgImportPathMap map[string]string

	buildFlagContent string
	Buf              bytes.Buffer
}

func NewGoFileContext(PkgImportPath string) *GoFileContext {
	return &GoFileContext{
		pkgImportPathMap: map[string]string{},

		packageName:   path.Base(PkgImportPath),
		pkgImportPath: PkgImportPath,
	}
}

func (gotpl *GoFileContext) GetPkgImportPath() string {
	return gotpl.pkgImportPath
}

func (gotpl *GoFileContext) MustWriteGoTypes(objTyp udwGoParser.Type) string {
	ObjectTypeStr, importPathList := udwGoParser.MustWriteGoTypes(gotpl.pkgImportPath, objTyp)

	gotpl.AddImportPathList(importPathList)
	return ObjectTypeStr
}

func (gotpl *GoFileContext) MustWriteGoTypePackagePrefix(pkgImportPath string) string {
	if pkgImportPath == gotpl.pkgImportPath {
		return ""
	}
	gotpl.AddImportPath(pkgImportPath)
	return path.Base(pkgImportPath) + "."
}

func (gotpl *GoFileContext) AddImportPathList(importPathList []string) {
	for _, importPath := range importPathList {
		gotpl.AddImportPath(importPath)
	}
}

func (gotpl *GoFileContext) AddImportPath(importPath string) {
	gotpl.pkgImportPathMap[importPath] = ""
}
func (gotpl *GoFileContext) AddImportPathWithAlias(importPath string, alias string) {
	gotpl.pkgImportPathMap[importPath] = alias
}
func (gotpl *GoFileContext) AddUnderScoreImportPath(importPath string) {
	gotpl.AddImportPathWithAlias(importPath, "_")
}

func (gotpl *GoFileContext) SetBuildFlagContent(content string) {
	gotpl.buildFlagContent = content
}
func (gotpl *GoFileContext) SetBuildFlagList(flagList []string) {
	gotpl.buildFlagContent = strings.Join(flagList, ",")
}

func (gotpl *GoFileContext) MergeFile(gotpl1 *GoFileContext) *GoFileContext {
	if gotpl.packageName != gotpl1.packageName {
		panic(fmt.Errorf("[goFileTplRequest.Merge] can not merge packageName different %s %s",
			gotpl.packageName, gotpl1.packageName))
	}
	out := &GoFileContext{
		packageName:      gotpl.packageName,
		pkgImportPathMap: map[string]string{},
	}
	for pkg, alias := range gotpl.pkgImportPathMap {
		out.pkgImportPathMap[pkg] = alias
	}
	for pkg, alias := range gotpl1.pkgImportPathMap {
		out.pkgImportPathMap[pkg] = alias
	}

	return out
}

func (gotpl *GoFileContext) MustWriteFile(_filepath string, body []byte) {
	udwFile.MustCheckContentAndWriteFileWithMkdirWithCorrectFold(_filepath, gotpl.getGoFileContent(body))
}
func (gotpl *GoFileContext) MustWriteFileWithSelfBuffer(_filepath string) {
	udwFile.MustCheckContentAndWriteFileWithMkdirWithCorrectFold(_filepath, gotpl.getGoFileContent(gotpl.Buf.Bytes()))
}
func (gotpl *GoFileContext) getGoFileContent(content []byte) []byte {
	var _buf bytes.Buffer
	if gotpl.buildFlagContent != "" {
		_buf.WriteString("// +build ")
		_buf.WriteString(gotpl.buildFlagContent)
		_buf.WriteString("\n\n")
	}
	_buf.WriteString(`package `)
	_buf.WriteString(gotpl.packageName)
	_buf.WriteByte('\n')
	if len(gotpl.pkgImportPathMap) > 0 {
		_buf.WriteString(`
import (
`)
		pairList := udwMap.MapStringStringToKeyValuePairListAes(gotpl.pkgImportPathMap)
		for i := range pairList {
			alias := pairList[i].Value
			importPkg := pairList[i].Key
			if alias == "" {
				_buf.WriteString("	" + udwGoTypeMarshal.WriteStringToGolangDoubleQuotation(importPkg) + "\n")
			} else {
				_buf.WriteString(alias + "	" + udwGoTypeMarshal.WriteStringToGolangDoubleQuotation(importPkg) + "\n")
			}
		}
		_buf.WriteString(`)
`)
	}

	_buf.Write(content)
	outB := _buf.Bytes()

	return outB
}

func (gotpl *GoFileContext) MustWriteGoTypesV2(typ types.Type) string {
	buf := &bytes.Buffer{}
	types.WriteType(buf, typ, func(typesPkg *types.Package) string {
		importPath := typesPkg.Path()
		if importPath == gotpl.pkgImportPath {
			return ""
		} else {
			gotpl.AddImportPath(importPath)
			return importPath
		}
	})
	return buf.String()
}
