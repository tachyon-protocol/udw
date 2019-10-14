package udwGoBuildCtx

import (
	"github.com/tachyon-protocol/udw/udwMap"
	"github.com/tachyon-protocol/udw/udwProjectPath"
	"github.com/tachyon-protocol/udw/udwStrings"
	"go/build"
	"path/filepath"
	"runtime"
	"strings"
)

type CtxReq struct {
	GoPathList   []string
	GoPathString string
	GoRoot       string

	TagList []string

	TargetOsArch string
	TargetOs     string
	TargetArch   string

	BuildTargetPkgPath string
}

func NewCtx(req CtxReq) *Ctx {
	ctx := &Ctx{}
	osSet := udwMap.StringListToSetString(GetOsList())
	archSet := udwMap.StringListToSetString(GetArchList())
	setOs := ""
	setArch := ""
	setOsFn := func(os string) {
		if setOs != "" && os != setOs {
			if (os == "android" && setOs == "linux") || (os == "linux" && setOs == "android") {
				setOs = "android"
				return
			}
			panic("w22mg869ut " + os + " " + setOs)
		}
		setOs = os
	}
	setArchFn := func(arch string) {
		if setArch != "" && arch != setArch {
			panic("pswcmedrjg " + arch + " " + setArch)
		}
		setArch = arch
	}
	if req.TargetOs != "" {
		setOsFn(req.TargetOs)
	}
	if req.TargetArch != "" {
		setArchFn(req.TargetArch)
	}
	if req.TargetOsArch != "" {
		if strings.Contains(req.TargetOsArch, "_") == false {
			panic("7dyqrea3pt " + req.TargetOsArch)
		}
		setOsFn(udwStrings.StringBeforeFirstSubString(req.TargetOsArch, "_"))
		setArchFn(udwStrings.StringAfterFirstSubString(req.TargetOsArch, "_"))
	}
	if len(req.TagList) > 0 {
		for _, tag := range req.TagList {
			_, ok := osSet[tag]
			if ok {
				setOsFn(tag)
			}
			_, ok = archSet[tag]
			if ok {
				setArchFn(tag)
			}
		}
		ctx.tagList = req.TagList
	}
	if setOs != "" {
		ctx.goOs = setOs
	}
	if setArch != "" {
		ctx.goArch = setArch
	}

	if req.GoPathString != "" {
		ctx.goPathList = strings.Split(req.GoPathString, string(filepath.ListSeparator))
	}
	if req.GoRoot != "" {
		ctx.goRoot = req.GoRoot
	}
	if ctx.goRoot == "" {
		ctx.goRoot = GetDefaultGoRoot()
	}

	if len(req.GoPathList) > 0 {
		ctx.goPathList = req.GoPathList
	}

	if len(ctx.goPathList) == 0 && udwProjectPath.HasProjectPath() {
		ctx.goPathList = []string{udwProjectPath.MustGetProjectPath()}
	}
	buildCtx := build.Default

	if len(ctx.goPathList) == 0 && buildCtx.GOPATH != "" {
		ctx.goPathList = strings.Split(buildCtx.GOPATH, string(filepath.ListSeparator))
	}
	if ctx.goArch == "" {
		ctx.goArch = runtime.GOARCH
	}
	if ctx.goOs == "" {
		ctx.goOs = runtime.GOOS
	}
	if len(ctx.tagList) == 0 {
		ctx.tagList = buildCtx.BuildTags
	}
	ctx.BuildTargetPkgPath = req.BuildTargetPkgPath
	return ctx
}

func NewCtxDefault() *Ctx {
	return NewCtx(CtxReq{})
}

func (ctx *Ctx) ToGoBuildCtx() *build.Context {
	buildCtx := build.Default
	buildCtx.GOOS = ctx.goOs
	buildCtx.GOROOT = ctx.goRoot
	buildCtx.GOARCH = ctx.goArch
	buildCtx.GOPATH = strings.Join(ctx.goPathList, string(filepath.ListSeparator))
	buildCtx.BuildTags = ctx.tagList
	buildCtx.CgoEnabled = ctx.DisableCgo == false
	return &buildCtx
}

func MustGetGoBuildCtx(ctx CtxReq) *build.Context {
	return NewCtx(ctx).ToGoBuildCtx()
}
