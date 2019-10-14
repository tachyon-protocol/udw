package udwPcapIpPacket

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"strconv"
	"time"
)

func ReadPcapToIpPacket(buf []byte, cb func(item Item)) (errMsg string) {
	if len(buf) < 24 {
		return "t46j3xg2fp"
	}
	reader := udwBytes.NewBufReaderWithOk(buf)
	magic := reader.ReadLittleEndUint32()
	if magic != 0xa1b2c3d4 {
		return "b77crjtzyf"
	}
	version_major := reader.ReadLittleEndUint16()
	if version_major != 2 {
		return "7dbm9ezcys"
	}
	version_minor := reader.ReadLittleEndUint16()
	if version_minor != 4 {
		return "zrwah22vvf"
	}
	time_zone := reader.ReadLittleEndUint32()
	if time_zone != 0 {
		return "wnh4p4cc7q"
	}
	sigfigs := reader.ReadLittleEndUint32()
	if sigfigs != 0 {
		return "5hw6ev9rwe"
	}
	reader.ReadSliceBySize(4)
	network := reader.ReadLittleEndUint32()
	if network != 0 && network != 1 {
		return "ezvn2xfxsy " + strconv.Itoa(int(network))
	}
	for {
		if reader.IsEof() {
			return ""
		}
		t1 := reader.ReadLittleEndUint32()
		t2 := reader.ReadLittleEndUint32()
		item := Item{
			T: time.Unix(int64(t1), int64(t2)*1000).UTC(),
		}
		incl_len := reader.ReadLittleEndUint32()
		if incl_len <= 4 || int(incl_len)+4 > reader.GetRemainSize() {
			return "cbmntrw4pd"
		}
		reader.ReadSliceBySize(4)
		if network == 0 {
			loopbackHeader := reader.ReadLittleEndUint32()
			if loopbackHeader != 2 {
				return "ssf4sjtbvq"
			}
			item.IpPacket = reader.ReadSliceBySize(int(incl_len) - 4)
		} else if network == 1 {
			reader.ReadSliceBySize(12)
			ipType := reader.ReadBigEndUint16()
			data := reader.ReadSliceBySize(int(incl_len) - 14)
			if ipType != 0x800 && ipType != 0x86dd {

				continue
			}
			item.IpPacket = data
		}
		if reader.IsOk() == false {
			return "5rqhj8k6nf"
		}
		cb(item)
	}
	return ""
}
