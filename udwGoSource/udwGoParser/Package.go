package udwGoParser

import (
	"fmt"
	"sort"
)

func MustParsePackage(gopath string, pkgPath string) *Package {
	return NewProgram([]string{gopath}).GetPackageByPkgPath(pkgPath)
}

func MustParsePackegeFromDefaultEnv(pkgPath string) *Package {
	return NewProgramFromDefault().GetPackageByPkgPath(pkgPath)
}

type Package struct {
	program       *Program
	pkgPath       string
	pkgName       string
	pkgImportPath string
	dirPath       string
	importMap     map[string]bool
	funcList      []*FuncOrMethodDeclaration
	funcNameMap   map[string]bool
	methodList    []*FuncOrMethodDeclaration
	namedTypeList []*NamedType
	isInGoRoot    bool
}

func (pkg *Package) GetImportList() []string {
	output := make([]string, 0, len(pkg.importMap))
	for imp := range pkg.importMap {
		output = append(output, imp)
	}
	sort.Strings(output)
	return output
}

func (pkg *Package) AddImport(pkgPath string) {
	if pkgPath != "" {
		pkg.importMap[pkgPath] = true
	}
}

func (pkg *Package) GetNamedTypeMethodSet(typ *NamedType) (output []*FuncOrMethodDeclaration) {
	if !(typ.PkgImportPath == pkg.pkgImportPath) {
		panic(fmt.Errorf("can not get MethodSet on diff pacakge typ[%s] pkg[%s]", typ.PkgImportPath, pkg.pkgImportPath))
	}
	for _, decl := range pkg.methodList {
		recvier := decl.ReceiverType
		if recvier.GetKind() == Ptr {
			recvier = recvier.(*PointerType).Elem
		}
		if recvier.GetKind() != Named {
			panic(fmt.Errorf("[GetNamedTypeMethodSet] reciver is not a named type %T %s", recvier, recvier.GetKind()))
		}
		if recvier.(*NamedType).Name == typ.Name {
			output = append(output, decl)
		}
	}
	return output
}

func (pkg *Package) IsMain() bool {
	return pkg.pkgName == "main"
}

func (pkg *Package) IsInGoRoot() bool {
	return pkg.isInGoRoot
}

func (pkg *Package) GetPkgImportPath() string {
	if pkg.IsMain() {
		return "main"
	}
	return pkg.pkgImportPath
}

func (pkg *Package) GetDirPath() string {
	return pkg.dirPath
}

func (pkg *Package) LookupNamedType(name string) *NamedType {
	for i := range pkg.namedTypeList {
		if pkg.namedTypeList[i].Name == name {
			return pkg.namedTypeList[i]
		}
	}
	return nil
}

func (pkg *Package) GetFuncList() []*FuncOrMethodDeclaration {
	return pkg.funcList
}

func (pkg *Package) LookupFunc(name string) *FuncOrMethodDeclaration {
	for i := range pkg.funcList {
		if pkg.funcList[i].GetName() == name {
			return pkg.funcList[i]
		}
	}
	return nil
}

func (pkg *Package) mustAddFile(path string) {
	file := parseFile(pkg.pkgPath, path, pkg)
	for imp := range file.importMap {
		pkg.AddImport(imp)
	}
	for _, funcDecl := range file.funcList {
		name := funcDecl.GetName()
		if pkg.funcNameMap[name] {

			continue
		}
		pkg.funcNameMap[name] = true
		pkg.funcList = append(pkg.funcList, funcDecl)
	}
	for _, funcDecl := range file.methodList {
		pkg.methodList = append(pkg.methodList, funcDecl)
	}
	for _, namedType := range file.namedTypeList {
		pkg.namedTypeList = append(pkg.namedTypeList, namedType)
	}
}

func (pkg *Package) GetAllMethodOnNamedType(name string) (output []*FuncOrMethodDeclaration) {
	nameType := pkg.LookupNamedType(name)
	if nameType == nil {
		panic("[Package.GetAllMethodOnNamedType] can not found [" + name + "]")
	}
	output = pkg.GetNamedTypeMethodSet(nameType)
	return
}
