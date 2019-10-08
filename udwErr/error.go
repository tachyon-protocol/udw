package udwErr

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"runtime"
	"runtime/debug"
	"strconv"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicIfErrorMsg(errMsg string) {
	if errMsg != "" {
		panic(errMsg)
	}
}

func LogErrorWithStack(err error) {
	if err == nil {
		return
	}
	s := ""
	if err != nil {
		s = err.Error()
	}
	println("error", s, getCurrentAllStackString(1))
}

func LogUserErrorWithStack(err error) {
	if err == nil {
		return
	}
	s := ""
	if err != nil {
		s = err.Error()
	}
	println("userError", s, getCurrentAllStackString(1))
}

func LogError(err error) {
	s := ""
	if err != nil {
		s = err.Error()
	}
	println("error", s)
}

var PrintStack = debug.PrintStack

type PanicErr struct {
	PanicObj interface{}
}

func (e PanicErr) Error() string {
	return fmt.Sprintf("%#v", e.PanicObj)
}

func PanicToError(f func()) (err error) {
	hasFinish := false
	defer func() {
		if hasFinish {
			return
		}
		out := recover()
		err = interfaceToError(out)
	}()
	f()
	hasFinish = true
	return nil
}

func PanicToErrorMsg(f func()) (errMsg string) {
	hasFinish := false
	defer func() {
		if hasFinish {
			return
		}
		out := recover()
		errMsg = interfaceToStringNotEmpty(out)
	}()
	f()
	hasFinish = true
	return errMsg
}

func PanicToErrorAndLog(f func()) (err error) {
	hasFinish := false
	defer func() {
		if hasFinish {
			return
		}
		out := recover()
		err = interfaceToError(out)
		LogErrorWithStack(err)
	}()
	f()
	hasFinish = true
	return nil
}

func interfaceToError(i interface{}) error {
	switch out := i.(type) {
	case error:
		return out
	case string:
		return errors.New(out)
	default:
		return PanicErr{PanicObj: i}
	}
}

func interfaceToStringNotEmpty(i interface{}) (outS string) {
	switch out := i.(type) {
	case error:
		outS = out.Error()
	case string:
		outS = out
	default:
		return fmt.Sprintf("%#v", out)
	}
	if outS == "" {
		outS = `<"">`
	}
	return outS
}

func PanicToErrorWithStackAndLog(f func()) (err error) {
	hasFinish := false
	defer func() {
		if hasFinish {
			return
		}
		out := recover()
		err = interfaceToError(out)
		LogErrorWithStack(err)
		err = errors.New(err.Error() + "\n" + getCurrentAllStackString(1))
	}()
	f()
	hasFinish = true
	return nil
}

func ErrorToMsg(err error) (errMsg string) {
	if err == nil {
		return ""
	}
	return err.Error()
}
func ErrorMsgToErr(errMsg string) (err error) {
	if errMsg == "" {
		return nil
	}
	return errors.New(errMsg)
}

func PanicToErrorMsgWithStack(f func()) (errMsg string) {
	hasFinish := false
	defer func() {
		if hasFinish {
			return
		}
		out := recover()
		err := interfaceToError(out)
		errMsg = err.Error() + "\n" + getCurrentAllStackString(1)
	}()
	f()
	hasFinish = true
	return ""
}

func PanicToErrorMsgWithStackAndLog(f func()) {
	hasFinish := false
	defer func() {
		if hasFinish {
			return
		}
		out := recover()
		err := interfaceToError(out)
		errMsg := err.Error() + "\n" + getCurrentAllStackString(1)
		fmt.Println(errMsg)
	}()
	f()
	hasFinish = true
	return
}

func PanicToErrorMsgAndStack(f func()) (errMsg string, stack string) {
	hasFinish := false
	defer func() {
		if hasFinish {
			return
		}
		out := recover()
		errMsg = interfaceToStringNotEmpty(out)
		stack = getCurrentAllStackString(1)
	}()
	f()
	hasFinish = true
	return "", ""
}

func PanicToCallback(f func(), panicFn func(errMsg string)) {
	hasFinish := false
	defer func() {
		if hasFinish {
			return
		}
		out := recover()
		errMsg := interfaceToStringNotEmpty(out)
		panicFn(errMsg)
	}()
	f()
	hasFinish = true
	return
}

type stringWrapError string

func (err stringWrapError) Error() string {
	return string(err)
}

func New(s string) error {
	return stringWrapError(s)
}

type RpcErr struct {
	RemoteLogicErr string

	NetworkErr string
}

func (e *RpcErr) Error() string {
	str := ""
	if e.RemoteLogicErr != "" {
		str += e.RemoteLogicErr
	}
	if e.NetworkErr != "" {
		str += " " + e.NetworkErr
	}
	return str
}

func getCurrentAllStackString(skip int) string {
	pcs := make([]uintptr, 32)
	thisLen := runtime.Callers(skip+2, pcs)
	frames := runtime.CallersFrames(pcs[:thisLen])
	buf := bytes.Buffer{}
	for {
		frame, more := frames.Next()
		buf.WriteString("    ")
		buf.WriteString(frame.File)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(frame.Line))
		buf.WriteByte(' ')
		buf.WriteString(frame.Function)
		buf.WriteByte(' ')
		buf.WriteString(udwStrconv.FormatUint64Hex(uint64(frame.PC)))
		buf.WriteByte('\n')
		if more == false {
			break
		}
	}
	return buf.String()
}
