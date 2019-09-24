package udwDebug

import (
	"encoding/hex"
	"github.com/tachyon-protocol/udw/udwCryptoMd5"
	"strconv"
)

func DebugBinary(b []byte) string {
	s := udwCryptoMd5.Md5Hex(b) + " " + strconv.Itoa(len(b)) + " " + hex.EncodeToString(b)
	return s
}
