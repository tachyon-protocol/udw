package udwTask

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"testing"
)

func TestOne(t *testing.T) {
	tm := NewLimitThreadTaskManager(1)
	result_chan := make(chan int, 1)
	var result_locker sync.Mutex
	var result_var int
	task := TaskFunc(func() {
		result_chan <- 1
		result_locker.Lock()
		result_var = 2
		result_locker.Unlock()
	})
	tm.AddTask(task)
	tm.AddTaskNewThread(task)
	udwTest.Equal(<-result_chan, 1, "not run result_chan not match")
	udwTest.Equal(<-result_chan, 1, "not run result_chan not match")
	tm.Close()
	udwTest.Equal(result_var, 2, "not run result_var not match")
}

func TestOneThread(t *testing.T) {
	tm := NewLimitThreadTaskManager(1)
	var result_var int
	result_chan := make(chan int, 2)
	task := TaskFunc(func() {
		result_var = 2
		result_chan <- 1
	})
	tm.AddTask(task)
	tm.AddTask(task)
	for i := 0; i < 2; i++ {
		if <-result_chan != 1 {
			t.Fatalf("result_chan not match")
		}
	}
	tm.Close()
	if result_var != 2 {
		t.Fatalf("result_var not match")
	}
}

func BenchmarkMulitThread(b *testing.B) {
	tm := NewLimitThreadTaskManager(10)
	tm.AddTask(TaskFunc(func() {
		task := TaskFunc(func() {
		})
		for i := 0; i < b.N; i++ {
			tm.AddTask(task)
		}
	}))

	tm.Close()
}
