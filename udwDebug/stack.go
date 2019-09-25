package udwDebug

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"runtime"
	"strconv"
)

type StackFuncCall struct {
	File     string
	Line     int
	Pc       uintptr
	FuncName string
}

type Stack []StackFuncCall

func (s *Stack) ToString() (output string) {
	output = ""
	for _, call := range *s {
		output += fmt.Sprintf("%s\n\t%s:%d:%x\n", call.FuncName, call.File, call.Line, call.Pc)
	}
	return
}

func GetCurrentStack(skip int) (stack *Stack) {
	pcs := make([]uintptr, 32)
	thisLen := runtime.Callers(skip+2, pcs)
	s := make(Stack, thisLen)
	stack = &s
	for i := 0; i < thisLen; i++ {
		f := runtime.FuncForPC(pcs[i])
		file, line := f.FileLine(pcs[i])
		(*stack)[i] = StackFuncCall{
			Pc:       pcs[i],
			FuncName: f.Name(),
			File:     file,
			Line:     line - 1,
		}
	}
	return
}

func GetCurrentOneStackString(skip int) string {
	var pcs [1]uintptr
	thisLen := runtime.Callers(skip+2, pcs[:])
	if thisLen == 0 {
		return "[No stack]"
	}
	f := runtime.FuncForPC(pcs[0])
	file, line := f.FileLine(pcs[0])
	return file + ":" + strconv.Itoa(line)
}

func GetCurrentAllStackString(skip int) string {
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

func PrintCurrentAllStack(skip int) {
	fmt.Println(GetCurrentAllStackString(skip))
}

func GetAllStack() []byte {
	buf := make([]byte, 32*1024)
	n := runtime.Stack(buf, true)
	return buf[:n]
}

func PrintCurrentStackV2() {
	buf := make([]byte, 32*1024)
	n := runtime.Stack(buf, false)
	fmt.Println(string(buf[:n]))
}

func PrintAllStackV2() {
	buf := make([]byte, 32*1024)
	n := runtime.Stack(buf, true)
	fmt.Println(string(buf[:n]))
}
func PrintCallerPos() {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d %s\n", file, line, f.Name())
}
func GetCallerPos() string {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return fmt.Sprintf("%s:%d %s\n", file, line, f.Name())
}
