package udwChan

import (
	"github.com/tachyon-protocol/udw/udwRand"
	"testing"
	"time"
)

func TestChan_ChanV2(t *testing.T) {
	for i := 0; i < 100; i++ {
		_chan := MakeChan(0)
		d := udwRand.IntBetween(1, 10)
		go func() {
			time.Sleep(time.Duration(d) * time.Millisecond)
			_chan.Send(1)
		}()
		go func() {
			time.Sleep(time.Duration(d) * time.Millisecond)
			_chan.Close()
		}()
		ir, isClose := _chan.Receive()
		if isClose {

		} else {
			_i, ok := ir.(int)
			_ = _i
			if ok {

			}
		}
	}
}
