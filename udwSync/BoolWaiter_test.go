package udwSync

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
	"time"
)

func TestBoolWaiter(ot *testing.T) {
	{
		bw := BoolWaiter{}
		bw.Set(true)
		udwTest.Equal(bw.Get(), true)
		udwTest.Equal(bw.WaitTrueWithTimeout(time.Millisecond*100), true)
		bw.Set(false)
		udwTest.Equal(bw.Get(), false)
		udwTest.Equal(bw.WaitTrueWithTimeout(time.Millisecond*100), false)
		bw.Set(false)
		go func() {
			time.Sleep(time.Millisecond * 50)
			bw.Set(true)
		}()
		udwTest.Equal(bw.WaitTrueWithTimeout(time.Millisecond*100), true)
		udwTest.Equal(bw.Get(), true)
		bw.Set(true)
		bw.Set(true)
		bw.Set(false)
		bw.Set(false)

		bw.Set(true)
		go func() {
			time.Sleep(time.Millisecond * 50)
			bw.Set(false)
		}()

		bw.Wait(false)
		bw.Set(false)
		bw.Wait(false)
		go func() {
			time.Sleep(time.Millisecond * 50)
			bw.Set(true)
		}()
		bw.Wait(true)
	}
	{
		bw := BoolWaiter{}
		go func() {
			time.Sleep(time.Millisecond * 50)
			bw.Set(true)
		}()
		bw.Wait(true)
	}
	{
		bw := BoolWaiter{}
		bw.Wait(false)
	}
	{
		bw := BoolWaiter{}
		ret := bw.WaitWithTimeout(false, time.Second)
		udwTest.Equal(ret, true)
		ret = bw.WaitWithTimeout(true, time.Millisecond*50)
		udwTest.Equal(ret, false)
	}
}
