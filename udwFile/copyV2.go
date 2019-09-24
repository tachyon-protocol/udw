package udwFile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type MustCopyV2Request struct {
	Src             string
	Dst             string
	IsCopyFilter    func(stat StatAndFullPath) bool
	AfterCopyFilter func(before StatAndFullPath, afterPath string) bool
}

func MustCopyV2(req MustCopyV2Request) {
	if req.IsCopyFilter == nil {
		req.IsCopyFilter = func(stat StatAndFullPath) bool {
			return true
		}
	}
	req.Src = MustGetFullPath(req.Src)
	req.Dst = MustGetFullPath(req.Dst)
	for _, thisStatAndFullPath := range MustGetAllFileAndDirectoryStat(req.Src) {
		if !req.IsCopyFilter(thisStatAndFullPath) {
			continue
		}
		rel, err := filepath.Rel(req.Src, thisStatAndFullPath.FullPath)
		if err != nil {
			panic(err)
		}
		targetFilePath := filepath.Join(req.Dst, rel)
		if FileInfoIsSymlink(thisStatAndFullPath.Fi) {
			link := MustReadSymlink(thisStatAndFullPath.FullPath)
			MustSymlink(link, targetFilePath)
			continue
		}
		if thisStatAndFullPath.Fi.IsDir() {
			MustMkdir777(targetFilePath)
			continue
		}
		err = copyOneFileV2(thisStatAndFullPath, targetFilePath)
		if err != nil {
			panic(err)
		}
		if req.AfterCopyFilter != nil {
			req.AfterCopyFilter(thisStatAndFullPath, targetFilePath)
		}
	}
}

func copyOneFileV2(srcStat StatAndFullPath, dst string) (err error) {
	in, err := os.Open(srcStat.FullPath)
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
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	if srcStat.Fi.Mode().Perm()&0100 > 0 {
		err := Chmod(dst, 0777)
		if err != nil {
			return err
		}
	} else {
		err := Chmod(dst, 0666)
		if err != nil {
			return err
		}
	}

	return
}
