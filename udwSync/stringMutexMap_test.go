package udwSync

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"strconv"
	"sync"
	"testing"
)

func TestStringMutexMap(ot *testing.T) {
	m := StringMutexMap{}
	m.LockByString("1")
	m.LockByString("2")
	m.UnlockByString("2")
	m.UnlockByString("1")

	wg := sync.WaitGroup{}
	num := 0
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			m.LockByString("1")
			m.LockByString("2")
			num++
			m.UnlockByString("2")
			m.UnlockByString("1")
			wg.Done()
		}()
	}
	wg.Wait()
	udwTest.Equal(num, 10)

	num = 0
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			m.LockByString("1")
			num++
			m.UnlockByString("1")
			wg.Done()
		}()
	}
	wg.Wait()
	udwTest.Equal(num, 10)

	for i := 0; i < 1024*2; i++ {
		iS := strconv.Itoa(i)
		m.LockByString(iS)
	}
	for i := 0; i < 1024*2; i++ {
		iS := strconv.Itoa(i)
		m.UnlockByString(iS)
	}

	const benchNum = 1024 * 2
	iList := []string{}
	for i := 0; i < benchNum; i++ {
		iList = append(iList, strconv.Itoa(i))
	}
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(benchNum)
		for i := 0; i < benchNum; i++ {
			m.LockByString(iList[i])
			m.UnlockByString(iList[i])
		}
	})
}
