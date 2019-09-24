package udwFile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type WalkContext struct {
	StatAndFullPath
	nextAction int
	root       string
}

func (ctx *WalkContext) StopWalk() {
	ctx.nextAction = 1
}

func (ctx *WalkContext) SkipDir() {
	ctx.nextAction = 2
}

func (ctx *WalkContext) GetRelPath() string {
	rel, err := filepath.Rel(ctx.root, ctx.StatAndFullPath.FullPath)
	if err != nil {
		panic(err)
	}
	return rel
}

func MustFileWalk(root string, fn func(ctx *WalkContext)) {
	root, err := GetFullPath(root)
	if err != nil {
		panic(err)
	}
	ctx := &WalkContext{
		root: root,
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}
		ctx.StatAndFullPath = StatAndFullPath{
			FullPath: path,
			Fi:       info,
		}
		ctx.nextAction = 0
		fn(ctx)
		switch ctx.nextAction {
		case 0:
			return nil
		case 1:
			return io.EOF
		case 2:
			return filepath.SkipDir
		default:
			panic(fmt.Errorf("[MustAllFileAndDirectoryWalk] unexpected ctx.nextAction"))
		}
	})
	if err != nil && err != io.EOF {
		panic(err)
	}
	return
}
