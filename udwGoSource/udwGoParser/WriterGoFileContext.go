package udwGoParser

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal"
	"path"
	"strings"
)

type GoFileWriter struct {
	packageName      string
	pkgImportPath    string
	pkgImportPathMap map[string]string

	buildFlagContent string
	Buf              bytes.Buffer
}

func NewGoFileContext(PkgImportPath string) *GoFileWriter {
	return &GoFileWriter{
		pkgImportPathMap: map[string]string{},

		packageName:   path.Base(PkgImportPath),
		pkgImportPath: PkgImportPath,
	}
}

func (gotpl *GoFileWriter) GetPkgImportPath() string {
	return gotpl.pkgImportPath
}

func (gotpl *GoFileWriter) MustWriteGoTypes(objTyp Type) string {
	ObjectTypeStr, importPathList := MustWriteGoTypes(gotpl.pkgImportPath, objTyp)
	gotpl.AddImportPathList(importPathList)
	return ObjectTypeStr
}

func (gotpl *GoFileWriter) MustWriteGoTypePackagePrefix(pkgImportPath string) string {
	if pkgImportPath == gotpl.pkgImportPath {
		return ""
	}
	gotpl.AddImportPath(pkgImportPath)
	return path.Base(pkgImportPath) + "."
}

func (gotpl *GoFileWriter) AddImportPathList(importPathList []string) {
	for _, importPath := range importPathList {
		gotpl.AddImportPath(importPath)
	}
}

func (gotpl *GoFileWriter) AddImportPath(importPath string) {
	gotpl.pkgImportPathMap[importPath] = ""
}
func (gotpl *GoFileWriter) AddImportPathWithAlias(importPath string, alias string) {
	gotpl.pkgImportPathMap[importPath] = alias
}
func (gotpl *GoFileWriter) AddUnderScoreImportPath(importPath string) {
	gotpl.AddImportPathWithAlias(importPath, "_")
}

func (gotpl *GoFileWriter) SetBuildFlagContent(content string) {
	gotpl.buildFlagContent = content
}

func (gotpl *GoFileWriter) MergeFile(gotpl1 *GoFileWriter) *GoFileWriter {
	if gotpl.packageName != gotpl1.packageName {
		panic(fmt.Errorf("[goFileTplRequest.Merge] can not merge packageName different %s %s",
			gotpl.packageName, gotpl1.packageName))
	}
	out := &GoFileWriter{
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

func (gotpl *GoFileWriter) MustWriteFile(_filepath string, body []byte) {
	udwFile.MustCheckContentAndWriteFileWithMkdirWithCorrectFold(_filepath, gotpl.getGoFileContent(body))
}
func (gotpl *GoFileWriter) MustWriteFileWithSelfBuffer(_filepath string) {
	udwFile.MustCheckContentAndWriteFileWithMkdirWithCorrectFold(_filepath, gotpl.getGoFileContent(gotpl.Buf.Bytes()))
}
func (gotpl *GoFileWriter) getGoFileContent(content []byte) []byte {
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
		for importPkg, alias := range gotpl.pkgImportPathMap {
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

func (gotpl *GoFileWriter) MustWriteNamedTypeDefine(namedType *NamedType) string {
	return `type ` + namedType.Name + ` ` + gotpl.MustWriteGoTypes(namedType.GetUnderType()) + "\n"
}

func (gotpl *GoFileWriter) MustWriteFuncOrMethodDecl(decl *FuncOrMethodDeclaration) string {
	_buf := bytes.Buffer{}
	_buf.WriteString("func ")
	if decl.ReceiverType != nil {
		_buf.WriteString("(" + decl.ReceiverVarName + " " + gotpl.MustWriteGoTypes(decl.ReceiverType) + ")")
	}
	_buf.WriteString(decl.Name + "(")
	for _, arg := range decl.InParameter {
		_buf.WriteString(arg.Name + " " + gotpl.MustWriteGoTypes(arg.Type) + ",")
	}
	_buf.WriteString(")")
	if len(decl.OutParameter) > 0 {
		_buf.WriteString("(")
		for _, arg := range decl.OutParameter {
			_buf.WriteString(arg.Name + " " + gotpl.MustWriteGoTypes(arg.Type) + ",")
		}
		_buf.WriteString(")")
	}
	return _buf.String()
}

type MethodCallWriteRequest struct {
	MethodExprString string
	InExprList       []string
	OutExprList      []string
	HasNewVarLeft    bool
}

func MustWriteMethodCall(req MethodCallWriteRequest) string {
	s := ""
	if len(req.OutExprList) > 0 {
		s = strings.Join(req.OutExprList, ",")
		if req.HasNewVarLeft {
			s += ":="
		} else {
			s += "="
		}
	}
	s += req.MethodExprString + "(" + strings.Join(req.InExprList, ",") + ")\n"
	return s
}
