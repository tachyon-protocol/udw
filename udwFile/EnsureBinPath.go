package udwFile

import (
	"os"
	"path/filepath"
	"time"

	"fmt"
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwPlatform"
	"github.com/tachyon-protocol/udw/udwRand"
	"github.com/tachyon-protocol/udw/udwStrings"
	"github.com/tachyon-protocol/udw/udwTime"
	"strings"
)

func MustEnsureBinPath(finalPath string) {
	if !udwPlatform.IsDarwin() && !udwPlatform.IsLinux() {
		panic("[MustEnsureBinPath] only support darwin and linux")
	}
	finalPath = MustGetFullPath(finalPath)
	filePathSymLinkList := MustGetAllSymlinkPathList(finalPath)
	basePath := filepath.Base(finalPath)
	if !MustFileExist(finalPath) {
		panic(fmt.Errorf("[MustEnsureBinPath] finalPath %s file not exist", finalPath))
	}
	pathList := getPathList()
	found := false
	for _, path := range pathList {
		endPath := filepath.Join(path, basePath)
		if endPath == finalPath {
			found = true
			break
		}
	}
	if found == false {
		panic(fmt.Errorf("[MustEnsureBinPath] finalPath %s is not in path", finalPath))
	}
	for _, path := range pathList {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		endPath := filepath.Join(path, basePath)
		if !MustFileExist(filepath.Join(path, basePath)) {
			continue
		}
		if udwStrings.IsInSlice(filePathSymLinkList, endPath) {
			continue
		}
		backPathDir := "/var/backup/bin/" + basePath + "_" + time.Now().Format(udwTime.FormatFileName) + "_" + udwRand.MustCryptoRandToReadableAlphaNum(6)
		MustMkdirForFile(backPathDir)
		udwCmd.MustRun("mv " + endPath + " " + backPathDir)
	}
}

func getPathList() []string {
	pathenv := os.Getenv("PATH")
	pathList := strings.Split(pathenv, ":")
	outPathList := []string{}
	for _, path := range pathList {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		outPathList = append(outPathList, path)
	}
	return outPathList
}
