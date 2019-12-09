package udwNet

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwErr"
	"strings"
)

const ErrMsgSocketCloseError = "use of closed network connection"

func GetSocketCloseError() error {
	return errors.New("use of closed network connection")
}

func IsSocketCloseError(err error) bool {
	return IsErrorMsgSocketClose(udwErr.ErrorToMsg(err))

}

func IsErrorMsgSocketClose(errMsg string) bool {
	if errMsg == "" {
		return false
	}
	return errMsg == "EOF" || errMsg == "io: read/write on closed pipe" ||
		strings.Contains(errMsg, "use of closed network connection") ||
		strings.Contains(errMsg, "socket is not connected") ||
		strings.Contains(errMsg, "broken pipe") ||
		strings.Contains(errMsg, "reset by peer")
}

func IsUseOfClosedNetworkConnection(err error) bool {
	return err != nil && strings.Contains(err.Error(), "use of closed network connection")
}

func IsBadFileDescripter(err error) bool {
	return err != nil && strings.Contains(err.Error(), "bad file descriptor")
}

func IsConnectionRefused(err error) bool {
	return err != nil && IsConnectionRefusedString(err.Error())
}

func IsConnectionRefusedString(errMsg string) bool {
	return errMsg != "" && strings.Contains(errMsg, "connection refused")
}

func IsNoRouteToHost(err error) bool {
	return err != nil && IsNoRouteToHostString(err.Error())
}

func IsNoRouteToHostString(errMsg string) bool {
	return errMsg != "" && strings.Contains(strings.ToLower(errMsg), "no route to host")
}

func IsNetworkIsUnreachable(err error) bool {
	return err != nil && ErrorMsgIsNetworkIsUnreachable(err.Error())
}

func ErrorMsgIsNetworkIsUnreachable(s string) bool {
	return strings.Contains(strings.ToLower(s), "network is unreachable")
}

func IsResourceBusy(err error) bool {
	return err != nil && strings.Contains(err.Error(), "resource busy")
}

func IsTimeoutError(err error) bool {
	return err != nil && IsTimeoutErrorString(err.Error())
}

func IsTimeoutErrorString(errMsg string) bool {
	return errMsg != "" && (strings.Contains(errMsg, "i/o timeout") ||
		strings.Contains(errMsg, "Client.Timeout exceeded while awaiting headers"))
}

func IsNoBufferSpaceAvailable(err error) bool {
	return err != nil && strings.Contains(err.Error(), "no buffer space available")
}

func IsMessageTooLong(err error) bool {
	return err != nil && strings.Contains(err.Error(), "message too long")
}

func IsInterruptedSystemCall(err error) bool {
	return err != nil && strings.Contains(err.Error(), "interrupted system call")
}

func NetErrorMsgSimple(errorMsg string) string {
	if strings.Contains(errorMsg, "i/o timeout") {
		errorMsg = "i/o timeout"
	} else if strings.Contains(errorMsg, "connection timed out") {
		errorMsg = "connection timed out"
	} else if strings.Contains(errorMsg, "connect: network is unreachable") {
		errorMsg = "connect: network is unreachable"
	} else if strings.Contains(errorMsg, "network is unreachable") {
		errorMsg = "network is unreachable"
	} else if strings.Contains(errorMsg, "Network is unreachable") {
		errorMsg = "network is unreachable"
	} else if strings.Contains(errorMsg, "write: invalid argument") {
		errorMsg = "write: invalid argument"
	} else if strings.Contains(errorMsg, "no route to host") {
		errorMsg = "no route to host"
	} else if strings.Contains(errorMsg, "no buffer space available") {
		errorMsg = "no buffer space available"
	} else if strings.Contains(errorMsg, "can't assign requested address") {
		errorMsg = "can't assign requested address"
	} else if strings.Contains(errorMsg, "read: socket is not connected") {
		errorMsg = "read: socket is not connected"
	} else if strings.Contains(errorMsg, "connection refused") {
		errorMsg = "connection refused"
	} else if strings.Contains(errorMsg, "operation timed out") {
		errorMsg = "operation timed out"
	} else if strings.Contains(errorMsg, "connection reset by peer") {
		errorMsg = "connection reset by peer"
	} else if strings.Contains(errorMsg, "bad file descriptor") {
		errorMsg = "bad file descriptor"
	} else if strings.Contains(errorMsg, "write: operation not permitted") {
		errorMsg = "write: operation not permitted"
	}
	return errorMsg
}

func IsErrorMsgNetwork(errorMsg string) bool {
	return strings.Contains(errorMsg, "unreachable") || strings.Contains(errorMsg, "i/o timeout") || strings.Contains(errorMsg, "An established connection was aborted by the software in your host machine") || strings.Contains(errorMsg, "[udwRpcGoClientTcp.sendRequest]")
}
