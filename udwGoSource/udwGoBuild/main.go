package udwGoBuild

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoBuildCtx"
	"github.com/tachyon-protocol/udw/udwProjectPath"
	"github.com/tachyon-protocol/udw/udwWindowsUI/rsrc"
	"path/filepath"
	"runtime"
)

const (
	TargetOsLinux      = "linux"
	TargetOsDarwin     = "darwin"
	TargetOsWindows    = "windows"
	TargetCpuArchAmd64 = "amd64"
	TargetCpuArchArm   = "arm"
	TargetCpuArch386   = "386"
	TargetLinuxAmd64   = "LinuxAmd64"
	TargetWindows386   = "Windows386"
	TargetWindowsAmd64 = "WindowsAmd64"
	TargetDarwinAmd64  = "DarwinAmd64"
)

type BuildRequest struct {
	PkgPath           string
	OutputPath        string
	TargetOs          string
	TargetCpuArch     string
	TargetOsCpuArch   string
	BuildTagList      []string
	DisableCgo        bool
	EnableRace        bool
	EnableSymbolDebug bool
	VariableMap       map[string]string
	GopathString      string
	GorootString      string

	WindowsDisableConsole bool

	WindowsNeedManifest bool

	WindowsManifestConfig rsrc.MustBuildToWin32WithCacheRequest

	GoPathOnlyBase bool
}

func (req *BuildRequest) ToCtx() *udwGoBuildCtx.Ctx {
	req.Init()
	ctx := udwGoBuildCtx.NewCtx(udwGoBuildCtx.CtxReq{
		GoPathString: req.GopathString,
		GoRoot:       req.GorootString,
		TagList:      req.BuildTagList,
		TargetOs:     req.TargetOs,
		TargetArch:   req.TargetCpuArch,
	})
	ctx.DisableCgo = req.DisableCgo
	ctx.BinOutputPath = req.OutputPath
	ctx.BuildTargetPkgPath = req.PkgPath
	ctx.BuildVariableMap = req.VariableMap
	ctx.EnableRace = req.EnableRace
	ctx.BuildWindowsNeedManifest = req.WindowsNeedManifest
	ctx.BuildWindowsIconPngContent = req.WindowsManifestConfig.IconPngContent
	ctx.BuildIsRequireRootPermission = req.WindowsManifestConfig.IsRequireRootPermission
	ctx.BuildWindowsDisableConsole = req.WindowsDisableConsole
	ctx.BuildGoPathOnlyBase = req.GoPathOnlyBase
	return ctx
}

func (req *BuildRequest) Init() {
	switch req.TargetOsCpuArch {
	case "":
	case TargetLinuxAmd64:
		req.TargetOs = TargetOsLinux
		req.TargetCpuArch = TargetCpuArchAmd64
	case TargetWindows386:
		req.TargetOs = TargetOsWindows
		req.TargetCpuArch = TargetCpuArch386
	case TargetWindowsAmd64:
		req.TargetOs = TargetOsWindows
		req.TargetCpuArch = TargetCpuArchAmd64
	case TargetDarwinAmd64:
		req.TargetOs = TargetOsDarwin
		req.TargetCpuArch = TargetCpuArchAmd64
	default:
		panic("unknow TargetOsCpuArch " + req.TargetOsCpuArch)
	}
	if (req.TargetOs != "" && req.TargetCpuArch == "") ||
		(req.TargetOs == "" && req.TargetCpuArch != "") {
		panic(`TargetOs and TargetCpuArch must both valid or ""`)
	}
	if req.TargetOs == "" && req.TargetCpuArch == "" && runtime.GOOS == "windows" {
		req.TargetOs = "windows"
		req.TargetCpuArch = "amd64"
	}
	if req.TargetOs == "" {
		req.TargetOs = runtime.GOOS
	}
	if req.TargetCpuArch == "" {
		req.TargetCpuArch = runtime.GOARCH
	}
	if req.GopathString == "" {
		req.GopathString = udwProjectPath.MustGetProjectPath()
	}
}

func MustBuild(req BuildRequest) *udwGoBuildCtx.Ctx {
	ctx := req.ToCtx()
	MustBuildCtx(ctx)
	return ctx

}

func MustBuildCtx(ctx *udwGoBuildCtx.Ctx) {
	envMap := map[string]string{}
	envMap["GOOS"] = ctx.GetGoOs()
	envMap["GOARCH"] = ctx.GetGoArch()
	envMap["GOPATH"] = ctx.GetGoPathString()
	envMap["GOROOT"] = ctx.GetGoRoot()
	if ctx.BuildGoPathOnlyBase {
		envMap["GoPathOnlyBase"] = "true"
	}
	if ctx.DisableCgo {
		envMap["CGO_ENABLED"] = "0"
	} else {
		envMap["CGO_ENABLED"] = "1"
	}
	if runtime.GOOS == "darwin" &&
		(ctx.IsGoOsLinux() || ctx.IsGoOsWindows()) {
		envMap["CGO_LDFLAGS"] = " -Wl,--unresolved-symbols=ignore-all"
	}
	if ctx.BuildCcPath != "" {
		envMap["CC"] = ctx.BuildCcPath
	}
	if ctx.BuildCxxPath != "" {
		envMap["CXX"] = ctx.BuildCxxPath
	}
	var staticLink = false
	if envMap["CC"] == "" && envMap["CXX"] == "" &&
		runtime.GOOS == "darwin" && ctx.DisableCgo == false {
		checkInstall := func(gcc string, cmd string) {
			exist := udwFile.MustFileExist(gcc)
			if !exist {
				fmt.Println(filepath.Base(gcc), "does not found, installing it")
				udwCmd.MustRun(cmd)
				exist = udwFile.MustFileExist(gcc)
				if exist {
					fmt.Println("install completed")
				} else {
					fmt.Println("install failed")
				}
			}
		}
		if ctx.IsGoOsLinux() && ctx.IsGoArchArm() {
			gcc := "/usr/local/arm-mac-linux-gnueabihf/bin/arm-mac-linux-gnueabihf-gcc"
			gxx := "/usr/local/arm-mac-linux-gnueabihf/bin/arm-mac-linux-gnueabihf-g++"
			checkInstall(gcc, "udw install linuxArmCrossGcc")
			envMap["CC"] = gcc
			envMap["CXX"] = gxx
		} else if ctx.IsGoOsLinux() && ctx.IsGoArchAmd64() {
			gcc := "/usr/local/x86_64-ubuntu14.04-linux-gnu/bin/x86_64-ubuntu14.04-linux-gnu-gcc"
			gxx := "/usr/local/x86_64-ubuntu14.04-linux-gnu/bin/x86_64-ubuntu14.04-linux-gnu-g++"
			checkInstall(gcc, "udw install linuxCrossGcc")
			envMap["CC"] = gcc
			envMap["CXX"] = gxx
		} else if ctx.IsGoOsWindows() {
			if ctx.IsGoArchAmd64() {
				gcc := `/usr/local/x86_64-mingw-w64/bin/x86_64-w64-mingw32-gcc`
				gxx := `/usr/local/x86_64-mingw-w64/bin/x86_64-w64-mingw32-g++`
				checkInstall(gcc, `udw install windowsCrossGccAmd64`)
				envMap[`CC`] = gcc
				envMap[`CXX`] = gxx
				staticLink = true
			} else {
				gcc := "/usr/local/i686-w64-mingw32/bin/i686-w64-mingw32-gcc"
				gxx := "/usr/local/i686-w64-mingw32/bin/i686-w64-mingw32-g++"
				checkInstall(gcc, "udw install windowsCrossGcc386")
				envMap["CC"] = gcc
				envMap["CXX"] = gxx
			}
		}
	}
	if ctx.BuildAndroidHome != "" {
		envMap["ANDROID_HOME"] = ctx.BuildAndroidHome
	}
	if ctx.BuildGoArm != "" {
		envMap["GOARM"] = ctx.BuildGoArm
	}
	if ctx.BuildGo386 != "" {
		envMap["GO386"] = ctx.BuildGo386
	}

	udwFile.MustMkdirForFile777(ctx.BinOutputPath)
	p := mustNewProgramV2MergeDefaultEnv(envMap, ctx)
	if ctx.EnableSymbolDebug {
		p.SetLdflags("")
	} else if staticLink {
		p.SetLdflags(`-linkmode external -extldflags -static -s -w`)
	}
	if len(ctx.GetTagList()) > 0 {
		p.SetBuildTagList(ctx.GetTagList())
	}
	outputExeFilePath := ctx.GetOutputExeFilePath()
	if ctx.IsGoOsWindows() {
		if ctx.EnableRace {
			panic("windows does not support -race")
		}
		if ctx.BuildWindowsNeedManifest {
			rsrc.MustBuildToWin32WithCache(rsrc.MustBuildToWin32WithCacheRequest{
				IconPngContent:          ctx.BuildWindowsIconPngContent,
				IsRequireRootPermission: ctx.BuildIsRequireRootPermission,
				TargetPackagePath:       ctx.BuildTargetPkgPath,
				GopathString:            ctx.GetGoPathString(),
			})
		}
	}
	if ctx.BuildWindowsDisableConsole {
		p.BuildWindowsNoConsoleExe()
	} else if ctx.BuildSo {
		p.BuildSoFile()
	} else if ctx.BuildA {
		p.BuildIosAFile()
	} else {
		p.MustGoInstall()
	}
	installOutputPath := ctx.GetGoInstallOutputExeFilePath()
	if installOutputPath != outputExeFilePath {
		udwFile.MustCopy(installOutputPath, outputExeFilePath)
	}
	return
}
