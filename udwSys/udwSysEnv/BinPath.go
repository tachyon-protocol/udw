package udwSysEnv

import (
	"github.com/tachyon-protocol/udw/udwMath"
	"github.com/tachyon-protocol/udw/udwPlatform"
	"github.com/tachyon-protocol/udw/udwStrings"
	"os"
	"path/filepath"
	"strings"
)

func RecoverPath() {
	if udwPlatform.IsDarwin() || udwPlatform.IsLinux() {
		const (
			a = "/usr/local/bin"
			b = "/bin"
			c = "/usr/bin"
		)
		targetPathOrderList := [3]string{a, b, c}
		targetPathIndexMap := make(map[string]int, 3)
		for _, p := range targetPathOrderList {
			targetPathIndexMap[p] = -1
		}

		pathToIndex := map[string]int{}
		originList := GetBinPathList()
		for i, p := range originList {
			for _, tp := range targetPathOrderList {
				if p == tp {
					targetPathIndexMap[tp] = i
					break
				}
			}

			pathToIndex[p] = i
		}
		needRewrite := false
		allTargetExist := true
		for _, i := range targetPathIndexMap {
			if i == -1 {
				allTargetExist = false
			}
		}
		if allTargetExist {
			for i, p := range targetPathOrderList {
				if i == len(targetPathOrderList)-1 {
					break
				}
				if targetPathIndexMap[p] > targetPathIndexMap[targetPathOrderList[i+1]] {
					needRewrite = true
					break
				}
			}
		} else {
			needRewrite = true
		}

		if !needRewrite {
			return
		}

		size := udwMath.IntMax([]int{len(pathToIndex) - 3, 3})
		newPathList := make([]string, 0, size)
		for _, p := range targetPathOrderList {
			delete(pathToIndex, p)
			newPathList = append(newPathList, p)
		}

		for _, p := range originList {
			if pathToIndex[p] == -1 {
				continue
			}
			if p == a || p == b || p == c {
				continue
			}
			newPathList = append(newPathList, p)
			pathToIndex[p] = -1
		}
		os.Setenv("PATH", strings.Join(newPathList, ":"))
	}
	if udwPlatform.IsWindows() {
		pathEnv := os.Getenv("PATH")
		pathList := GetBinPathList()
		change := false

		for _, needPath := range []string{
			`c:\windows\system32`,
			`c:\windows\system32\wbem`,
		} {
			if !udwStrings.IsInSlice(pathList, needPath) {
				change = true
				pathEnv += ";" + needPath
			}
		}
		if change {
			os.Setenv("PATH", pathEnv)
		}
		return
	}
}

func GetBinPathList() []string {
	return filepath.SplitList(os.Getenv("PATH"))
}
