package udwSingleFlight

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"testing"
	"time"
)

func TestGroup_Do(t *testing.T) {
	s := Group{}
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			err := udwErr.PanicToErrorMsgWithStack(func() {
				j, err := s.Do("1", func() (interface{}, error) {
					time.Sleep(time.Millisecond * 100)
					panic("one caller crashed!")
				})
				fmt.Println(j, err)
			})
			if err != "" {
				fmt.Println("panic:", err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	udwTest.Ok(true)
}
