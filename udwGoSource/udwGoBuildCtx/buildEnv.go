package udwGoBuildCtx

import (
	"path/filepath"
	"strings"
)

func (ctx *Ctx) UpdateFromGoEnv(env map[string]string) {
	if len(env) == 0 {
		return
	}
	var value string
	goRoot := env["GOROOT"]
	if goRoot != "" {
		ctx.goRoot = goRoot
	}
	goOs := env["GOOS"]
	if goOs != "" {
		ctx.goOs = goOs
	}
	goArch := env["GOARCH"]
	if goArch != "" {
		ctx.goArch = goArch
	}
	value = env["GOPATH"]
	if value != "" {
		ctx.goPathList = strings.Split(value, string(filepath.ListSeparator))
	}
	cGO_ENABLED := env["CGO_ENABLED"]
	if cGO_ENABLED == "1" {
		ctx.DisableCgo = false
	} else if cGO_ENABLED == "0" {
		ctx.DisableCgo = true
	}
	value = env["CC"]
	if value != "" {
		ctx.BuildCcPath = value
	}
	value = env["CXX"]
	if value != "" {
		ctx.BuildCxxPath = value
	}
	value = env["GOARM"]
	if value != "" {
		ctx.BuildGoArm = value
	}
	value = env["GO386"]
	if value != "" {
		ctx.BuildGo386 = value
	}
	value = env["ANDROID_HOME"]
	if value != "" {
		ctx.BuildAndroidHome = value
	}
	value = env["CGO_CFLAGS"]
	if value != "" {
		ctx.BuildCgoCflags = value
	}
	value = env["CGO_CPPFLAGS"]
	if value != "" {
		ctx.BuildCgoCppFlags = value
	}
	value = env["CGO_LDFLAGS"]
	if value != "" {
		ctx.BuildCgoLdFlags = value
	}
}
