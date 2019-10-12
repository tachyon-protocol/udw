package udwGoBuild

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwProfile/udwProfileDelay"
	"github.com/tachyon-protocol/udw/udwStrings"
	"path/filepath"
	"strings"
)

func (p *programV2) BuildSoFile() {
	if !p.MustIsTargetExist() {
		panic(fmt.Errorf("[BuildSoFile] target not exist target:[%s]", p.targetPackagePathOrFilePath))
	}

	pkgName := p.GetGoos() + "_" + p.GetGoArch() + "_shared"
	pkgPath := filepath.Join(p.ctx.GetFirstGoPathString(), "pkg", pkgName)

	installBinPath := p.ctx.GetGoInstallOutputExeFilePath()
	slice := udwStrings.StringSliceMerge("go", "build", "-buildmode=c-shared", "-i", "-pkgdir", pkgPath,
		p.getBuildFlagCmdSlice(), "-o="+installBinPath, p.targetPackagePathOrFilePath)
	p.mustUdwGoInstall(slice)
}

func (p *programV2) BuildIosAFile() {
	if !p.MustIsTargetExist() {
		panic(fmt.Errorf("[BuildIosAFile] target not exist target:[%s]", p.targetPackagePathOrFilePath))
	}

	pkgName := p.GetGoos() + "_" + p.GetGoArch()
	if len(p.buildTagList) > 0 {
		pkgName += "_" + strings.Join(p.buildTagList, "_")
	}
	pkgPath := filepath.Join(p.gopathList[0], "pkg", pkgName)

	installBinPath := p.ctx.GetGoInstallOutputExeFilePath()
	udwProfileDelay.P()
	p.mustUdwGoInstall(udwStrings.StringSliceMerge("go", "build", "-buildmode=c-archive",
		p.getBuildFlagCmdSlice(), "-i", "-pkgdir", pkgPath,
		"-o="+installBinPath, p.targetPackagePathOrFilePath))
	udwProfileDelay.P()
}

func (p *programV2) BuildWindowsNoConsoleExe() {
	p.SetLdflags("-H=windowsgui -s -w")
	if !p.MustIsTargetExist() {
		panic("[BuildWindowsNoConsoleExe] target not exist target:[" + p.targetPackagePathOrFilePath + "]")
	}

	pkgPath := filepath.Join(p.gopathList[0], "pkg", p.GetGoos()+"_"+p.GetGoArch()+"_windowsgui")

	installBinPath := p.ctx.GetGoInstallOutputExeFilePath()

	p.mustUdwGoInstall(udwStrings.StringSliceMerge("go", "build",
		p.getBuildFlagCmdSlice(), "-pkgdir", pkgPath, "-o="+installBinPath, p.targetPackagePathOrFilePath))
}
