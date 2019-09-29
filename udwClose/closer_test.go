package udwClose

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwSync"
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"testing"
	"time"
)

func TestCloser1(ot *testing.T) {
	c := Closer{}
	udwTest.Equal(c.IsClose(), false)
	num := udwSync.NewInt(0)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(10 * time.Millisecond)
			c.CloseWithCallback(func() {
				num.Add(1)
			})
			udwTest.Equal(num.Get(), 1)
			wg.Done()
		}()
	}
	wg.Wait()
	udwTest.Equal(num.Get(), 1)
	udwTest.Equal(c.IsClose(), true)
}

func TestCloser2(ot *testing.T) {
	c := Closer{}
	udwTest.Equal(c.IsClose(), false)
	go func() {
		time.Sleep(10 * time.Millisecond)
		c.Close()
	}()
	select {
	case <-time.After(time.Second):
		panic("[TestCloser2] timeout")
	case <-c.GetCloseChan():
	}
	udwTest.Equal(c.IsClose(), true)
}

func TestCloser3(ot *testing.T) {
	c := Closer{}
	udwTest.Equal(c.IsClose(), false)
	c.Close()
	select {
	case <-time.After(time.Second):
		panic("[TestCloser2] timeout")
	case <-c.GetCloseChan():
	}
	udwTest.Equal(c.IsClose(), true)
}

func TestCloser4(ot *testing.T) {
	c := Closer{}
	udwTest.Equal(c.IsClose(), false)
	t := time.Now()
	for i := 0; i < 1000000; i++ {
		c.IsClose()
	}
	fmt.Println("IsClose false", time.Since(t)/1000000)
	c.Close()
	t = time.Now()
	for i := 0; i < 1000000; i++ {
		c.Close()
	}
	fmt.Println("Close", time.Since(t)/1000000)
	t = time.Now()
	for i := 0; i < 1000000; i++ {
		c.GetCloseChan()
	}
	fmt.Println("GetCloseChan", time.Since(t)/1000000)
	udwTest.Equal(c.IsClose(), true)

	t = time.Now()
	for i := 0; i < 1000000; i++ {
		c := Closer{}
		c.Close()
		c.IsClose()
	}
	fmt.Println("all 1", time.Since(t)/1000000)
}

func TestClose5(ot *testing.T) {
	c := Closer{}
	closeCallList := []string{}
	c.AddOnClose(func() {
		closeCallList = append(closeCallList, "1")
	})
	c.AddOnClose(func() {
		closeCallList = append(closeCallList, "2")
	})
	c.CloseWithCallback(func() {
		closeCallList = append(closeCallList, "3")
	})
	udwTest.Equal(closeCallList, []string{"3", "2", "1"})
	c.Close()
}
func TestCloser_AddUpperCloser(t *testing.T) {
	c := Closer{}
	c2 := Closer{}
	c2.AddUpperCloser(&c)
	c.Close()
	udwTest.Equal(c2.IsClose(), true)
	udwTest.Equal(c.IsClose(), true)
}
