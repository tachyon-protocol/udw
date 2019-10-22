// +build darwin dragonfly freebsd linux netbsd openbsd

package udwSsh

import (
	osUser "os/user"
	"path/filepath"
)

func certFilePathList() []string {
	u, err := osUser.Current()
	if err != nil {
		return nil
	}
	return []string{
		filepath.Join(u.HomeDir, ".ssh", "id_rsa"),
	}
}
