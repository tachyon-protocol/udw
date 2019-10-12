package udwProjectPath

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"os"
	"path/filepath"
	"sync"
)

func MustGetProjectPath() string {
	if getIsDisallowGetProjectPathThisProcess() {
		panic("tgvbv7ckqv disallow get project path")
	}
	doInit()
	if gProjectPath == "" {
		panic("34thp9f7n3 ProjectPath not config")
	}
	return gProjectPath
}

func HasProjectPath() bool {
	if getIsDisallowGetProjectPathThisProcess() {
		return false
	}
	doInit()
	return gProjectPath != ""
}

func MustPathInProject(p string) string {
	return filepath.Join(MustGetProjectPath(), p)
}

func DisallowGetProjectPathThisProcess() {
	isDisallowGetProjectPathThisProcessLocker.Lock()
	isDisallowGetProjectPathThisProcess = true
	isDisallowGetProjectPathThisProcessLocker.Unlock()
}

var isDisallowGetProjectPathThisProcess = false
var isDisallowGetProjectPathThisProcessLocker sync.Mutex

func getIsDisallowGetProjectPathThisProcess() bool {
	isDisallowGetProjectPathThisProcessLocker.Lock()
	out := isDisallowGetProjectPathThisProcess
	isDisallowGetProjectPathThisProcessLocker.Unlock()
	return out
}

var gProjectPathOnce sync.Once
var gProjectPath string

func doInit() {
	gProjectPathOnce.Do(func() {
		p, err := os.Getwd()
		if err != nil {
			return
		}
		p, err = udwFile.SearchFileInParentDir(p, ".project_root.sign")
		if err != nil {

			p, err = udwFile.SearchFileInParentDir(p, ".udw.yml")
			if err != nil {
				return
			}
		}
		gProjectPath, err = filepath.Abs(p)
		if err != nil {
			return
		}
	})
}
