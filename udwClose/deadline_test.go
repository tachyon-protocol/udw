package udwClose

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"testing"
	"time"
)

func TestCloser_SetTimeout1(t *testing.T) {
	iLock := &sync.Mutex{}
	i := 0
	closer := NewCloser()
	closer.AddOnClose(func() {
		iLock.Lock()
		i++
		iLock.Unlock()
	})
	closer.SetTimeoutFromStart(time.Millisecond * 100)
	time.Sleep(time.Millisecond * 150)
	iLock.Lock()
	_i := i
	iLock.Unlock()
	udwTest.Ok(_i == 1)
	udwTest.Ok(closer.IsClose())
}

func TestCloser_SetTimeout2(t *testing.T) {
	iLock := &sync.Mutex{}
	i := 0
	closer := NewCloser()
	closer.AddOnClose(func() {
		iLock.Lock()
		i++
		iLock.Unlock()
	})
	closer.SetTimeoutFromStart(time.Millisecond * 150)
	time.Sleep(time.Millisecond * 100)
	iLock.Lock()
	_i := i
	iLock.Unlock()
	udwTest.Ok(_i == 0)
	udwTest.Ok(closer.IsClose() == false)
	closer.SetTimeoutFromStart(time.Millisecond * 220)
	time.Sleep(time.Millisecond * 100)
	iLock.Lock()
	_i = i
	iLock.Unlock()
	udwTest.Ok(_i == 0)
	udwTest.Ok(closer.IsClose() == false)
	time.Sleep(time.Millisecond * 30)
	iLock.Lock()
	_i = i
	iLock.Unlock()
	udwTest.Ok(_i == 1)
	udwTest.Ok(closer.IsClose() == true)
}

func TestCloser_SetTimeout3(t *testing.T) {
	closer := NewCloser()
	closer.SetTimeoutFromStart(time.Millisecond * 150)
	time.Sleep(time.Millisecond * 100)
	udwTest.Ok(closer.IsClose() == false)
	closer.SetTimeoutFromStart(time.Millisecond * 10)
	udwTest.Ok(closer.IsClose() == true)
}

func TestCloser_ClearTimeout(t *testing.T) {
	closer0 := NewCloser()
	closer0.SetTimeoutFromStart(time.Millisecond * 50)
	time.Sleep(time.Millisecond * 25)
	closer0.ClearTimeout()
	time.Sleep(time.Millisecond * 50)
	udwTest.Ok(closer0.IsClose() == false)

	closer1 := NewCloser()
	closer1.ClearTimeout()
	udwTest.Ok(closer1.IsClose() == false)

	closer2 := NewCloser()
	closer2.SetTimeoutFromStart(time.Millisecond)
	time.Sleep(time.Millisecond * 2)
	closer2.ClearTimeout()
	udwTest.Ok(closer2.IsClose() == true)
}
