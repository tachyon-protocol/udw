package udwChan

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwRand"
	"testing"
	"time"
)

func TestChan_ChanV2(t *testing.T) {
	for i := 0; i < 10000; i++ {
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
		i, isClose := _chan.Receive()
		if isClose {
			fmt.Println("closed")
		} else {
			_i, ok := i.(int)
			if ok {
				fmt.Println("✔", _i)
			}
		}
	}
}
