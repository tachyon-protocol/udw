package udwSync

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"testing"
	"time"
)

func TestCond(ot *testing.T) {

	c := Cond{}
	c.InLock(func() {})
	c.Signal()
	c.Broadcast()
	c.Close()
}

func TestCond2(ot *testing.T) {
	c := Cond{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	a := 0
	go func() {
		c.WaitCheckAndDo(func() bool {
			return a == 1
		}, func() {
			wg.Done()
		})
	}()
	time.Sleep(10 * time.Millisecond)
	c.InLock(func() {
		a = 1
	})
	c.Signal()
	wg.Wait()
}

func TestCond3(ot *testing.T) {
	c := Cond{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	a := 0
	isReturnLocker := sync.Mutex{}
	isReturn := false
	go func() {
		c.WaitCheckAndDo(func() bool {
			return a == 1
		}, func() {
		})
		wg.Done()
		isReturnLocker.Lock()
		isReturn = true
		isReturnLocker.Unlock()
	}()
	time.Sleep(10 * time.Millisecond)
	isReturnLocker.Lock()
	udwTest.Equal(isReturn, false)
	isReturnLocker.Unlock()
	c.Close()
	wg.Wait()
}
