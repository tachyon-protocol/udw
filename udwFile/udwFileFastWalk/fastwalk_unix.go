// +build linux darwin freebsd openbsd netbsd
// +build !appengine

package udwFileFastWalk

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const blockSize = 8 << 10

const unknownFileMode os.FileMode = os.ModeNamedPipe | os.ModeSocket | os.ModeDevice

func readDir(dirName string, fn func(dirName, entName string, typ os.FileMode) error) error {
	fd, err := syscall.Open(dirName, 0, 0)
	if err != nil {
		return &os.PathError{Op: "open", Path: dirName, Err: err}
	}
	defer syscall.Close(fd)

	buf := make([]byte, blockSize)
	bufp := 0
	nbuf := 0
	skipFiles := false
	for {
		if bufp >= nbuf {
			bufp = 0
			nbuf, err = syscall.ReadDirent(fd, buf)
			if err != nil {
				return os.NewSyscallError("readdirent", err)
			}
			if nbuf <= 0 {
				return nil
			}
		}
		consumed, name, typ := parseDirEnt(buf[bufp:nbuf])
		bufp += consumed
		if name == "" || name == "." || name == ".." {
			continue
		}

		if typ == unknownFileMode {
			fi, err := os.Lstat(dirName + "/" + name)
			if err != nil {

				if os.IsNotExist(err) {
					continue
				}
				return err
			}
			typ = fi.Mode() & os.ModeType
		}
		if skipFiles && typ.IsRegular() {
			continue
		}
		if err := fn(dirName, name, typ); err != nil {
			if err == SkipFiles {
				skipFiles = true
				continue
			}
			return err
		}
	}
}

func parseDirEnt(buf []byte) (consumed int, name string, typ os.FileMode) {

	dirent := (*syscall.Dirent)(unsafe.Pointer(&buf[0]))
	if v := unsafe.Offsetof(dirent.Reclen) + unsafe.Sizeof(dirent.Reclen); uintptr(len(buf)) < v {
		panic(fmt.Sprintf("buf size of %d smaller than dirent header size %d", len(buf), v))
	}
	if len(buf) < int(dirent.Reclen) {
		panic(fmt.Sprintf("buf size %d < record length %d", len(buf), dirent.Reclen))
	}
	consumed = int(dirent.Reclen)
	if direntInode(dirent) == 0 {
		return
	}
	switch dirent.Type {
	case syscall.DT_REG:
		typ = 0
	case syscall.DT_DIR:
		typ = os.ModeDir
	case syscall.DT_LNK:
		typ = os.ModeSymlink
	case syscall.DT_BLK:
		typ = os.ModeDevice
	case syscall.DT_FIFO:
		typ = os.ModeNamedPipe
	case syscall.DT_SOCK:
		typ = os.ModeSocket
	case syscall.DT_UNKNOWN:
		typ = unknownFileMode
	default:

		return
	}

	nameBuf := (*[unsafe.Sizeof(dirent.Name)]byte)(unsafe.Pointer(&dirent.Name[0]))
	nameLen := direntNamlen(dirent)

	if nameLen == 1 && nameBuf[0] == '.' {
		name = "."
	} else if nameLen == 2 && nameBuf[0] == '.' && nameBuf[1] == '.' {
		name = ".."
	} else {
		name = string(nameBuf[:nameLen])
	}
	return
}
