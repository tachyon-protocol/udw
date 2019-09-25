package udwRandNewId

import (
	"encoding/binary"
	"github.com/tachyon-protocol/udw/udwCryptoSha3"
	"github.com/tachyon-protocol/udw/udwHex"
	"github.com/tachyon-protocol/udw/udwRand"
	"strings"
	"sync"
	"time"
)

func NewId() (outS string) {
	gNewIdObjectPoolEntryLocker.Lock()
	if len(gNewIdObjectPoolEntry.buf) == 0 {
		gNewIdObjectPoolEntry.buf = make([]byte, 72)
		gNewIdObjectPoolEntry.shakeHash = udwCryptoSha3.NewShake256()
	} else {
		gNewIdObjectPoolEntry.shakeHash.Reset()
	}
	entry2 := gNewIdObjectPoolEntry

	entry2.shakeHash.Write(entry2.buf[:48])
	now := time.Now()
	binary.LittleEndian.PutUint64(entry2.buf[:8], uint64(now.UnixNano()))

	udwRand.MustCryptoRandBytesWithBuf(entry2.buf[8:64])
	entry2.shakeHash.Write(entry2.buf[:64])
	entry2.shakeHash.Read(entry2.buf[:48])
	outS = udwRand.EncodeReadableAlphaNumForRandNoAlloc(entry2.buf[:48], entry2.buf[48:72])
	gNewIdObjectPoolEntryLocker.Unlock()
	return outS

}

func NewIdLikeUuid() string {
	byteId := newIdByte()
	return strings.Join([]string{
		udwHex.EncodeBytesToString(byteId[:4]),
		udwHex.EncodeBytesToString(byteId[4:6]),
		udwHex.EncodeBytesToString(byteId[6:8]),
		udwHex.EncodeBytesToString(byteId[8:10]),
		udwHex.EncodeBytesToString(byteId[10:16]),
	}, "-")
}

func NewIdLikeWindowsGUID() string {
	byteId := newIdByte()
	data := strings.ToUpper(udwHex.EncodeBytesToString(byteId[:16]))
	return "{" + data[0:8] + "-" + data[8:12] + "-" + data[12:16] + "-" + data[16:20] + "-" + data[20:32] + "}"
}

func newIdByte() []byte {
	var _buf [64]byte
	now := time.Now()
	nowS := now.String()
	udwRand.MustCryptoRandBytesWithBuf(_buf[:])
	h := udwCryptoSha3.New512()
	h.Write([]byte(nowS))
	h.Write(_buf[:])
	return h.Sum(nil)
}

var gNewIdObjectPoolEntry newIdObjectPoolEntry
var gNewIdObjectPoolEntryLocker sync.Mutex

type newIdObjectPoolEntry struct {
	buf       []byte
	shakeHash udwCryptoSha3.ShakeHash
}
