package udwGoParser

import (
	"fmt"
	"path"
)

type file struct {
	packageName string

	importMap      map[string]bool
	aliasImportMap map[string]string
	funcList       []*FuncOrMethodDeclaration
	methodList     []*FuncOrMethodDeclaration
	namedTypeList  []*NamedType
	pkg            *Package
}

func (pkg *file) addImport(pkgPath string, aliasPath string) {
	if pkgPath == "" {
		return
	}
	pkg.importMap[pkgPath] = true
	if aliasPath == "" {
		aliasPath = path.Base(pkgPath)
	}
	pkg.aliasImportMap[aliasPath] = pkgPath
}

func (gofile *file) lookupFullPackagePath(pkgAliasPath string) (string, error) {
	pkgPath := gofile.aliasImportMap[pkgAliasPath]
	if pkgPath == "" {
		return pkgAliasPath, fmt.Errorf("unable to find import alias %s", pkgAliasPath)
	}
	return pkgPath, nil
}
