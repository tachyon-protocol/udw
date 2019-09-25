package udwBytes

import (
	"bytes"
	"io"
)

type ReadUtilRequest struct {
	R         io.Reader
	ToReadBuf []byte
	Sub       []byte
}
type ReadUtilResponse struct {
	TotalReadLen      int
	SubStringStartPos int
	ErrMsg            string
}

func ReadUtil(req ReadUtilRequest) (resp ReadUtilResponse) {
	toReadBuf := req.ToReadBuf
	for {
		if len(toReadBuf) == 0 {
			resp.ErrMsg = "xz7r2ramfj"
			return resp
		}
		n, err := req.R.Read(toReadBuf)
		if err != nil {
			resp.ErrMsg = err.Error()
			return resp
		}
		resp.TotalReadLen += n
		toReadBuf = toReadBuf[n:]
		resp.SubStringStartPos = bytes.Index(req.ToReadBuf[:resp.TotalReadLen], req.Sub)
		if resp.SubStringStartPos != -1 {
			break
		}
	}
	return resp
}

func ReadCStringFromBufToByte(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			return b[:i]
		}
	}
	return b
}

func ReadCStringFromBuf(b []byte) string {
	return string(ReadCStringFromBufToByte(b))
}
