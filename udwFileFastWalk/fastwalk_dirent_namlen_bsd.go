package udwFileFastWalk

import "syscall"

func direntNamlen(dirent *syscall.Dirent) uint64 {
	return uint64(dirent.Namlen)
}
