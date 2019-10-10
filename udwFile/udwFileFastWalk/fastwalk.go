package udwFileFastWalk

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var TraverseLink error

var SkipFiles error

func Walk(root string, walkFn func(path string, typ os.FileMode) error) error {

	pkgInit()
	numWorkers := 4
	if n := runtime.NumCPU(); n > numWorkers {
		numWorkers = n
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	w := &walker{
		fn:       walkFn,
		enqueuec: make(chan walkItem, numWorkers),
		workc:    make(chan walkItem, numWorkers),
		donec:    make(chan struct{}),

		resc: make(chan error, numWorkers),
	}
	defer close(w.donec)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go w.doWork(&wg)
	}
	todo := []walkItem{{dir: root}}
	out := 0
	for {
		workc := w.workc
		var workItem walkItem
		if len(todo) == 0 {
			workc = nil
		} else {
			workItem = todo[len(todo)-1]
		}
		select {
		case workc <- workItem:
			todo = todo[:len(todo)-1]
			out++
		case it := <-w.enqueuec:
			todo = append(todo, it)
		case err := <-w.resc:
			out--
			if err != nil {
				return err
			}
			if out == 0 && len(todo) == 0 {

				select {
				case it := <-w.enqueuec:
					todo = append(todo, it)
				default:
					return nil
				}
			}
		}
	}
}

var gInitOnce sync.Once

func pkgInit() {
	gInitOnce.Do(func() {
		TraverseLink = errors.New("fastwalk: traverse symlink, assuming target is a directory")
		SkipFiles = errors.New("fastwalk: skip remaining files in directory")
	})
}

func (w *walker) doWork(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-w.donec:
			return
		case it := <-w.workc:
			select {
			case <-w.donec:
				return
			case w.resc <- w.walk(it.dir, !it.callbackDone):
			}
		}
	}
}

type walker struct {
	fn func(path string, typ os.FileMode) error

	donec    chan struct{}
	workc    chan walkItem
	enqueuec chan walkItem
	resc     chan error
}

type walkItem struct {
	dir          string
	callbackDone bool
}

func (w *walker) enqueue(it walkItem) {
	select {
	case w.enqueuec <- it:
	case <-w.donec:
	}
}

func (w *walker) onDirEnt(dirName, baseName string, typ os.FileMode) error {
	joined := dirName + string(os.PathSeparator) + baseName
	if typ == os.ModeDir {
		w.enqueue(walkItem{dir: joined})
		return nil
	}

	err := w.fn(joined, typ)
	if typ == os.ModeSymlink {
		if err == TraverseLink {

			w.enqueue(walkItem{dir: joined, callbackDone: true})
			return nil
		}
		if err == filepath.SkipDir {

			return nil
		}
	}
	return err
}

func (w *walker) walk(root string, runUserCallback bool) error {
	if runUserCallback {
		err := w.fn(root, os.ModeDir)
		if err == filepath.SkipDir {
			return nil
		}
		if err != nil {
			return err
		}
	}

	return readDir(root, w.onDirEnt)
}
