package udwSync

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"testing"
)

func TestWaitGroup(ot *testing.T) {
	{
		wg := Rc{}
		wg.Add(1)
		wg.Done()
		wg.Wait()
	}
	{
		i := Int{}
		wg := Rc{}
		wg.Inc()
		go func() {
			i.Inc()
			wg.Dec()
		}()
		wg.Wait()
		udwTest.Equal(i.Get(), 1)
	}
	{
		counter := Int{}
		wg := Rc{}
		wg.Inc()
		go func() {
			defer wg.Dec()
			for i := 0; i < 10; i++ {
				wg.Inc()
				go func() {
					counter.Inc()
					wg.Dec()
				}()
			}
		}()
		for i := 0; i < 10; i++ {
			wg.Wait()
			udwTest.Equal(counter.Get(), 10)
		}
	}
	{
		for i := 0; i < 10; i++ {
			wg := sync.WaitGroup{}
			wg.Add(1)
			for j := 0; j < 10; j++ {
				go func() {
					wg.Wait()
				}()
			}
			wg.Add(10)
			for j := 0; j < 10; j++ {
				go func() {
					wg.Add(1)
					go func() {
						wg.Done()
					}()
					wg.Done()
				}()
			}
			wg.Done()
			wg.Wait()
		}
	}

	{

		wg := Rc{}
		wg.Add(1)
		go func() {
			wg.Wait()
		}()
		go func() {
			wg.Add(1)
			wg.Done()
		}()
		wg.Done()
		wg.Wait()
	}
}
