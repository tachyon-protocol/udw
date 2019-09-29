package udwErr

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestPanicToError(ot *testing.T) {
	flag := 1
	err := PanicToError(func() {
		flag = 2
	})
	udwTest.Equal(flag, 2)
	udwTest.Equal(err, nil)

	err = PanicToError(func() {
		flag = 3
		panic(nil)

	})
	udwTest.Equal(flag, 3)
	udwTest.Equal(err.Error(), "<nil>")

	err = PanicToError(func() {
		panic(1)

	})
	udwTest.Equal(flag, 3)
	udwTest.Ok(err != nil)
	udwTest.Equal(err.Error(), "1")

	err = PanicToError(func() {
		panic("abc\n")
	})
	udwTest.Ok(err != nil)
	udwTest.Equal(err.Error(), "abc\n")
}

func TestPanicToCallback(ot *testing.T) {
	lastErrMsg := ""
	PanicToCallback(func() {
		panic("TestPanicToCallback abc")
	}, func(errMsg string) {
		lastErrMsg = errMsg
	})
	udwTest.Equal(lastErrMsg, "TestPanicToCallback abc")

	lastErrMsg = ""
	PanicToCallback(func() {
		panic("")
	}, func(errMsg string) {
		lastErrMsg = errMsg
	})
	udwTest.Equal(lastErrMsg, `<"">`)

	lastErrMsg = ""
	PanicToCallback(func() {
		panic(nil)
	}, func(errMsg string) {
		lastErrMsg = errMsg
	})
	udwTest.Equal(lastErrMsg, `<nil>`)
}
