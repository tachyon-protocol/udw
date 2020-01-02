// +build !darwin go1.12

package udwFile

import "os"

func FileSync(f *os.File) error { return f.Sync() }
