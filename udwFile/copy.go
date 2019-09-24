package udwFile

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) (err error) {
	if src == dst {

		return
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		if ErrorIsFileNotFound(err) {
			err = os.MkdirAll(filepath.Dir(dst), os.FileMode(0777))
			if err != nil {
				return err
			}
			out, err = os.Create(dst)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("[CopyFile] createDst err[%s]", err.Error())
		}
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}

	return
}

func MustCopyFile(src, dst string) {
	err := CopyFile(src, dst)
	if err != nil {
		panic(err)
	}
}

func MustCopyFileIgnoreNotExist(src, dst string) {
	err := CopyFile(src, dst)
	if ErrorIsFileNotFound(err) {
		return
	}
	if err != nil {
		panic(err)
	}
}

func MustCopy(src string, dst string) {
	MustCopyV2(MustCopyV2Request{
		Src: src,
		Dst: dst,
	})

}

func MustCopyWithFilter(src string, dst string, filter func(fullpath string) bool) {
	MustCopyV2(MustCopyV2Request{
		Src: src,
		Dst: dst,
		IsCopyFilter: func(stat StatAndFullPath) bool {
			return filter(stat.FullPath)
		},
	})
}

func MustCheckContentCopyWithoutMerge(src, dst string) {
	src = MustGetFullPath(src)
	dst = MustGetFullPath(dst)
	if !MustFileExist(dst) {
		MustCopy(src, dst)
		return
	}

	dstFileList := MustGetAllFileAndDirectoryStat(dst)
	dstFileMap := map[string]StatAndFullPath{}
	for _, stat := range dstFileList {
		dstFileMap[stat.FullPath] = stat
	}
	srcFileList := MustGetAllFileAndDirectoryStat(src)
	srcFileMap := map[string]StatAndFullPath{}
	for _, stat := range srcFileList {
		srcFileMap[stat.FullPath] = stat
	}

	for _, stat := range dstFileList {
		rel, err := filepath.Rel(dst, stat.FullPath)
		if err != nil {
			panic(err)
		}
		srcFilePath := filepath.Join(src, rel)
		_, ok := srcFileMap[srcFilePath]
		if !ok {
			MustDelete(stat.FullPath)
		}
	}

	for _, srcStat := range srcFileList {
		rel, err := filepath.Rel(src, srcStat.FullPath)
		if err != nil {
			panic(err)
		}
		dstFilePath := filepath.Join(dst, rel)
		dstStat, ok := dstFileMap[dstFilePath]
		if ok && MustFileSystemItemInfoEqual(srcStat, dstStat) {

			continue
		}

		MustDelete(dstFilePath)
		if FileInfoIsSymlink(srcStat.Fi) {
			link := MustReadSymlink(srcStat.FullPath)
			MustSymlink(link, dstFilePath)
			continue
		}
		if srcStat.Fi.IsDir() {
			MustMkdir(dstFilePath)
			continue
		}
		MustCopyFile(srcStat.FullPath, dstFilePath)
		continue
	}
}

func MustFileSystemItemInfoEqual(a StatAndFullPath, b StatAndFullPath) bool {
	if FileInfoIsSymlink(a.Fi) && FileInfoIsSymlink(b.Fi) {
		return MustReadSymlink(a.FullPath) == MustReadSymlink(b.FullPath)
	}
	if a.Fi.IsDir() && b.Fi.IsDir() {
		return true
	}
	if !a.Fi.IsDir() && !b.Fi.IsDir() {
		if a.Fi.Size() != b.Fi.Size() {
			return false
		}
		if a.Fi.Size() < 64*1024 {

			return bytes.Equal(MustReadFile(a.FullPath), MustReadFile(b.FullPath))
		}

		return mustFileContentEqual(a.FullPath, b.FullPath)
	}

	return false
}

func mustFileContentEqual(aPath string, bPath string) bool {
	fa, err := os.Open(aPath)
	if err != nil {
		panic(err)
	}
	defer fa.Close()
	fb, err := os.Open(bPath)
	if err != nil {
		panic(err)
	}
	defer fb.Close()
	bufA := make([]byte, 32*1024)
	bufB := make([]byte, 32*1024)
	for {
		nA, errA := fa.Read(bufA)
		if errA != nil && errA != io.EOF {
			panic(errA)
		}
		nB, errB := fb.Read(bufB)
		if errB != nil && errB != io.EOF {
			panic(errB)
		}

		if nA != nB {
			panic("[mustFileContentEqual] nA!=nB")
		}
		if bytes.Equal(bufA[:nA], bufB[:nB]) == false {
			return false
		}
		if errA == io.EOF && errB == io.EOF {
			break
		}
		if errA == io.EOF || errB == io.EOF {
			panic("[mustFileContentEqual] eof not equal.")
		}
	}
	return true
}
