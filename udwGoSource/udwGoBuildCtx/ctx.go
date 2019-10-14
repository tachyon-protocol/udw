package udwGoBuildCtx

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwMap"
	"github.com/tachyon-protocol/udw/udwStrings"
	"path/filepath"
	"runtime"
	"strings"
)

type Ctx struct {
	goPathList []string
	goRoot     string
	goOs       string
	goArch     string
	tagList    []string

	DisableCgo                   bool
	BinOutputPath                string
	BuildTargetPkgPath           string
	BuildVariableMap             map[string]string
	EnableRace                   bool
	EnableSymbolDebug            bool
	BuildWindowsNeedManifest     bool
	BuildWindowsIconPngContent   []byte
	BuildIsRequireRootPermission bool
	BuildWindowsDisableConsole   bool

	BuildCcPath         string
	BuildCxxPath        string
	BuildAndroidHome    string
	BuildGoArm          string
	BuildGo386          string
	BuildSo             bool
	BuildCgoCflags      string
	BuildCgoCppFlags    string
	BuildCgoLdFlags     string
	BuildA              bool
	BuildCgoDebug       bool
	BuildGoPathOnlyBase bool
}

func (ctx *Ctx) Clone() *Ctx {
	ctx2 := *ctx
	ctx2.goPathList = udwStrings.StringSliceClone(ctx2.goPathList)
	ctx2.tagList = udwStrings.StringSliceClone(ctx2.tagList)
	ctx2.BuildVariableMap = udwMap.MapStringStringClone(ctx2.BuildVariableMap)
	ctx2.BuildWindowsIconPngContent = udwBytes.Clone(ctx2.BuildWindowsIconPngContent)
	return &ctx2
}

func (ctx *Ctx) SetGoOs(os string) {
	ctx.goOs = os
}
func (ctx *Ctx) GetGoOs() string {
	return ctx.goOs
}

func (ctx *Ctx) SetGoArch(arch string) {
	ctx.goArch = arch
}
func (ctx *Ctx) GetGoArch() string {
	return ctx.goArch
}

func (ctx *Ctx) SetGoPathString(v string) {
	ctx.goPathList = strings.Split(v, string(filepath.ListSeparator))
}

func (ctx *Ctx) SetGoPathList(v []string) {
	ctx.goPathList = v
}

func (ctx *Ctx) GetGoPathList() []string {
	return ctx.goPathList
}

func (ctx *Ctx) GetGoPathString() string {
	return strings.Join(ctx.goPathList, string(filepath.ListSeparator))
}

func (ctx *Ctx) GetFirstGoPathString() string {
	if len(ctx.goPathList) > 0 {
		return ctx.goPathList[0]
	}
	return ""
}

func (ctx *Ctx) SetGoRoot(v string) {
	ctx.goRoot = v
}

func (ctx *Ctx) GetGoRoot() string {
	return ctx.goRoot
}

func (ctx *Ctx) GetGoSearchPathList() []string {
	out := []string{}
	out = append(out, ctx.goPathList...)
	out = append(out, ctx.goRoot)
	return out
}

func (ctx *Ctx) GetTagList() []string {
	return ctx.tagList
}

func (ctx *Ctx) SetTagList(tagList []string) {
	ctx.tagList = tagList
}

func (ctx *Ctx) IsGoOsWindows() bool {
	return ctx.goOs == "windows"
}

func (ctx *Ctx) IsGoOsLinux() bool {
	return ctx.goOs == "linux"
}

func (ctx *Ctx) IsGoOsDarwin() bool {
	return ctx.goOs == "darwin"
}

func (ctx *Ctx) IsGoOsAndroid() bool {
	return ctx.goOs == "android"
}

func (ctx *Ctx) IsGoArchArm() bool {
	return ctx.goArch == "arm"
}

func (ctx *Ctx) IsGoArchAmd64() bool {
	return ctx.goArch == "amd64"
}

func (ctx *Ctx) GetOutputExeFilePath() string {
	if ctx.BinOutputPath != "" {
		return udwFile.MustGetFullPath(ctx.BinOutputPath)
	}
	return ctx.GetGoInstallOutputExeFilePath()
}

func (ctx *Ctx) GetOutputExeFileExt() string {
	if ctx.IsGoOsWindows() {
		return ".exe"
	} else if ctx.BuildA {
		return ".a"
	} else if ctx.BuildSo {
		return ".so"
	}
	return ""
}

func (ctx *Ctx) GetGoInstallOutputExeFilePath() string {
	baseName := filepath.Base(ctx.BuildTargetPkgPath)
	if ctx.IsGoOsWindows() {
		if !strings.HasSuffix(baseName, ".exe") {
			baseName += ".exe"
		}
	} else if ctx.BuildA {
		baseName += ".a"
	} else if ctx.BuildSo {
		baseName += ".so"
	}
	if ctx.goOs == runtime.GOOS && ctx.goArch == runtime.GOARCH {
		return udwFile.MustGetFullPath(filepath.Join(ctx.GetFirstGoPathString(), "bin", baseName))
	}
	return udwFile.MustGetFullPath(filepath.Join(ctx.GetFirstGoPathString(), "bin", ctx.goOs+"_"+ctx.goArch, baseName))
}
