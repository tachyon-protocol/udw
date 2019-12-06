package udwRpc2

import (
	"fmt"
)

const (
	errCodeNormal  = 0
	errCodeNetwork = 1
	errCodeOther   = 2
)

type RpcError struct {
	errCode byte
	errMsg  string
}

func (err *RpcError) Error() string {
	return err.errMsg
}

func (err *RpcError) IsNetworkError() bool {
	return err.errCode == errCodeNetwork
}

func NewNetworkError(errMsg string) *RpcError {
	return &RpcError{
		errCode: errCodeNetwork,
		errMsg:  errMsg,
	}
}

func NewOtherError(errMsg string) *RpcError {
	return &RpcError{
		errCode: errCodeOther,
		errMsg:  errMsg,
	}
}

func PanicToErrMsg(f func()) (errMsg string) {
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
