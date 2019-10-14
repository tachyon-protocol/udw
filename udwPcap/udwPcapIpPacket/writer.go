package udwPcapIpPacket

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwFile"
	"os"
	"time"
)

type Item struct {
	T        time.Time
	IpPacket []byte
}

func MarshalToBuffer(_buf *udwBytes.BufWriter, itemList []Item) {
	_buf.Write_(gGlobalHeader)
	for _, item := range itemList {
		AppendItemToBuffer(_buf, item)
	}
}

var gGlobalHeader = []byte{
	0xd4, 0xc3, 0xb2, 0xa1,
	0x02, 0x00, 0x04, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x04, 0x00,
	0x00, 0x00, 0x00, 0x00,
}

var gLoopbackHeader = []byte{
	2, 0, 0, 0,
}

func MustAppendIpPacketToFile(filePath string, t time.Time, IpPacket []byte) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0777))
	if err != nil {
		if udwFile.ErrorIsFileNotFound(err) {
			udwFile.MustMkdirForFile777(filePath)
			f, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0777))
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		panic(err)
	}
	_buf := &udwBytes.BufWriter{}
	if st.Size() == 0 {
		_buf.Write_(gGlobalHeader)
	}
	AppendItemToBuffer(_buf, Item{
		T:        t,
		IpPacket: IpPacket,
	})
	f.Write(_buf.GetBytes())
	return
}

func GetGlobalHeader() []byte {
	return gGlobalHeader
}

func AppendItemToBuffer(_buf *udwBytes.BufWriter, item Item) {
	_buf.WriteLittleEndUint32(uint32(item.T.Unix()))
	_buf.WriteLittleEndUint32(uint32(item.T.Nanosecond() / 1000))
	thisLen := len(item.IpPacket) + 4
	_buf.WriteLittleEndUint32(uint32(thisLen))
	_buf.WriteLittleEndUint32(uint32(thisLen))
	_buf.Write_(gLoopbackHeader)
	_buf.Write_(item.IpPacket)
}
