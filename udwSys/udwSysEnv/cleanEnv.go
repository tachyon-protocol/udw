package udwSysEnv

import (
	"github.com/tachyon-protocol/udw/udwPlatform"
	"os"
)

func CleanEnv() {
	if udwPlatform.IsDarwin() {
		env := MustNewEnvFromCurrentProcess()
		needEnvKeyMap := map[string]bool{
			"USER":     true,
			"TEMP_DIR": true,
			"TMPDIR":   true,
			"GOPATH":   true,
			"HOME":     true,
		}
		for k := range env.ToMap() {
			if needEnvKeyMap[k] {
				continue
			}
			os.Unsetenv(k)
		}
		os.Setenv("LC_ALL", "en_US.UTF-8")
		os.Setenv("LC_CTYPE", "UTF-8")
		os.Setenv("SHELL", "/bin/bash")
		os.Setenv("PATH", "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin")
	} else {
		panic("[CleanEnv] TODO")
	}
}
