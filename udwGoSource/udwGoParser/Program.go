package udwGoParser

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoBuildCtx"
	"github.com/tachyon-protocol/udw/udwSingleFlight"
	"path/filepath"
	"strings"
	"sync"
)

type Program struct {
	packageLookupPathList []string
	goroot                string
	cachedPackageMap      map[string]*Package
	groupCache            udwSingleFlight.Group
	locker                sync.Mutex
	mainPackage           *Package
}

func NewProgramFromDefault() *Program {
	return &Program{
		packageLookupPathList: append(udwGoBuildCtx.GetDefaultGoPathList(), udwGoBuildCtx.GetDefaultGoRoot()),
		goroot:                udwGoBuildCtx.GetDefaultGoRoot(),
		cachedPackageMap:      map[string]*Package{},
	}
}

func NewProgram(lookupPathList []string) *Program {
	return &Program{
		packageLookupPathList: lookupPathList,
		cachedPackageMap:      map[string]*Package{},
	}
}

func (prog *Program) GetPackageByPkgPath(pkgPath string) *Package {
	pkg := prog.getCachedPackage(pkgPath)
	if pkg != nil {
		return pkg
	}
	_, err := prog.groupCache.Do(pkgPath, func() (interface{}, error) {
		pkg = prog.mustParsePackage(pkgPath)
		prog.setCachedPackage(pkg)
		return nil, nil
	})
	if err != nil {
		panic(err)
	}
	return prog.getCachedPackage(pkgPath)
}

func (prog *Program) GetPackageByPkgImportPath(PkgImportPath string) *Package {
	if PkgImportPath == "main" {
		return prog.mainPackage
	}
	return prog.GetPackageByPkgPath(PkgImportPath)
}

func (prog *Program) GetNamedType(PkgImportPath string, name string) *NamedType {
	pkg := prog.GetPackageByPkgImportPath(PkgImportPath)
	return pkg.LookupNamedType(name)
}

func (prog *Program) mustParsePackage(pkgPath string) *Package {
	var dirPath string
	var found bool
	isInGoroot := false
	for _, lookupPath := range prog.packageLookupPathList {
		dirPath = filepath.Join(lookupPath, "src", pkgPath)
		if udwFile.MustFileExist(dirPath) {
			if lookupPath == prog.goroot {
				isInGoroot = true
			}
			found = true
			break
		}
	}
	if !found {
		panic(fmt.Errorf("can not found pkgPath %s %#v", pkgPath, prog.packageLookupPathList))
	}

	pkg := &Package{
		importMap:   map[string]bool{},
		funcNameMap: map[string]bool{},
		pkgPath:     pkgPath,
		dirPath:     dirPath,
		isInGoRoot:  isInGoroot,
		program:     prog,
	}

	for _, path := range udwFile.MustReadDirFileOneLevel(dirPath) {
		if strings.HasSuffix(path, ".go") {
			pkg.mustAddFile(filepath.Join(dirPath, path))
		}
	}

	if pkg.IsMain() {
		for i := range pkg.namedTypeList {
			pkg.namedTypeList[i].PkgImportPath = "main"
		}
		prog.mainPackage = pkg
	}
	return pkg
}

func (prog *Program) getCachedPackage(pkgPath string) *Package {
	prog.locker.Lock()
	defer prog.locker.Unlock()
	return prog.cachedPackageMap[pkgPath]
}

func (prog *Program) setCachedPackage(pkg *Package) {
	prog.locker.Lock()
	defer prog.locker.Unlock()
	prog.cachedPackageMap[pkg.pkgPath] = pkg
}
