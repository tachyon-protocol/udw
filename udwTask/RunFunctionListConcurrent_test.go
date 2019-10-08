package udwTask

import (
	"github.com/tachyon-protocol/udw/udwSync"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestRunFunctionListConcurrent(ot *testing.T) {
	counter := udwSync.NewInt(0)
	addOneFn := func() {
		counter.Add(1)
	}
	addTwoFn := func() {
		counter.Add(2)
	}
	fnList := []func(){
		addOneFn,
		addOneFn,
		addTwoFn,
	}
	RunFunctionListConcurrent(fnList...)
	udwTest.Equal(counter.Get(), 4)

	fnList = []func(){}
	udwTest.BenchmarkWithRepeatNum(100, func() {
		RunFunctionListConcurrent(fnList...)
	})
	fnList = append(fnList, func() {}, func() {})
	benchFn := func() {
		RunFunctionListConcurrent(fnList...)
	}
	udwTest.BenchmarkWithRepeatNum(100, benchFn)
}
