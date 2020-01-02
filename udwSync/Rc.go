package udwSync

import "sync"

type Rc struct {
	locker   sync.Mutex
	wgLocker sync.Mutex
	rc       int
}

func (rc *Rc) Add(i int) {
	rc.locker.Lock()
	oldRc := rc.rc
	rc.rc = oldRc + i
	if i > 0 && oldRc == 0 {
		rc.wgLocker.Lock()
	}
	if i < 0 && rc.rc == 0 {
		rc.wgLocker.Unlock()
	}
	rc.locker.Unlock()
}
func (rc *Rc) Done() {
	rc.Add(-1)
}
func (rc *Rc) Wait() {
	rc.wgLocker.Lock()
	rc.wgLocker.Unlock()
}

func (rc *Rc) Inc() {
	rc.Add(1)
}

func (rc *Rc) Dec() {
	rc.Add(-1)
}

func (rc *Rc) GetRc() int {
	rc.locker.Lock()
	out := rc.rc
	rc.locker.Unlock()
	return out
}
