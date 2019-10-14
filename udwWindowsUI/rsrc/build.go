package rsrc

import (
	"github.com/tachyon-protocol/udw/udwCache"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwImage/udwImageIco"
	"github.com/tachyon-protocol/udw/udwJson"
	"path/filepath"
)

type MustBuildRequest struct {
	ManifestFileContent []byte
	IconContentList     [][]byte
	OutputFilePath      string
	Arch                string
}

func MustBuildWithCache(req MustBuildRequest) {
	cacheFileList := []string{
		req.OutputFilePath,
		"src/github.com/tachyon-protocol/udw/udwWindowsUI/rsrc",
	}

	udwCache.MustMd5FileChangeCache("rsrc_"+udwJson.MustMarshalToString(req),
		cacheFileList,
		func() {
			if len(req.ManifestFileContent) == 0 {
				req.ManifestFileContent = gRootManifestContent
			}
			MustBuild(req)
		})
}

type MustBuildToWin32WithCacheRequest struct {
	IconPngContent          []byte
	TargetPackagePath       string
	IsRequireRootPermission bool
	GopathString            string
}

func MustBuildToWin32WithCache(req MustBuildToWin32WithCacheRequest) {
	sysoPath := filepath.Join(req.GopathString, "src", req.TargetPackagePath, "zzzig_rsrc.syso")
	udwCache.MustMd5FileChangeCache("rsrc_"+udwJson.MustMarshalToString(req),
		[]string{
			sysoPath,
			"src/github.com/tachyon-protocol/udw/udwWindowsUI/rsrc",
		},
		func() {
			udwFile.MustDelete(filepath.Join(req.GopathString, "src", req.TargetPackagePath, "~i_rsrc.syso"))
			buildReq := MustBuildRequest{
				OutputFilePath: sysoPath,
				Arch:           "386",
			}
			if len(req.IconPngContent) > 0 {
				icoContent := udwImageIco.MustEncodePngContentToIcoContent(req.IconPngContent)
				buildReq.IconContentList = [][]byte{icoContent}
			}

			if req.IsRequireRootPermission {
				buildReq.ManifestFileContent = gRootManifestContent
			} else {
				buildReq.ManifestFileContent = gNormalManifestContent
			}

			MustBuild(buildReq)
		})
}

func MustBuild(req MustBuildRequest) {
	err := run(req)
	if err != nil {
		panic(err)
	}
}
