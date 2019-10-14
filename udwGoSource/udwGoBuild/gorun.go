package udwGoBuild

import (
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoBuildCtx"
)

func MustGoRunSimple(pkgPath string, remainArgs []string) {
	resp := MustBuild(BuildRequest{
		PkgPath: pkgPath,
	})
	firstGoPath := udwGoBuildCtx.NewCtxDefault().GetFirstGoPathString()
	udwCmd.CmdSlice(append([]string{resp.GetOutputExeFilePath()}, remainArgs...)).
		SetDir(firstGoPath).
		MustSetEnv("GOPATH", firstGoPath).
		MustRun()

}

func MustGoRunRace(pkgPath string, args []string) {
	resp := MustBuild(BuildRequest{
		PkgPath:    pkgPath,
		EnableRace: true,
	})
	firstGoPath := udwGoBuildCtx.NewCtxDefault().GetFirstGoPathString()
	udwCmd.CmdSlice(append([]string{resp.GetOutputExeFilePath()}, args...)).
		SetDir(firstGoPath).
		MustSetEnv("GOPATH", firstGoPath).
		MustRun()
}
