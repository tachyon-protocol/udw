// +build !darwin

package udwFile

import "os"

func FixedTmpDir() string {
	return os.TempDir()
}
