package udwKvdbV2

import (
	"encoding/binary"
	"errors"
	"github.com/tachyon-protocol/udw/udwCryptoAesCtrV7"
	"github.com/tachyon-protocol/udw/udwStrconv"
)

const overheadOfCrypto = 36 + 32

func cryptoBlock(psk *[32]byte, data []byte) []byte {
	length := uint32(len(data))
	var lengthBytes [4]byte
	binary.LittleEndian.PutUint32(lengthBytes[:], length)
	if psk == nil {
		return append(lengthBytes[:], data...)
	} else {
		return append(udwCryptoAesCtrV7.Encrypt32(psk, lengthBytes[:]), udwCryptoAesCtrV7.Encrypt32(psk, data)...)
	}
}

type decryptItem struct {
	data         []byte
	lengthInFile int
}

func decryptoBlockList(psk *[32]byte, data []byte) (list []decryptItem, err error) {
	if psk == nil {
		for len(data) >= 4 {
			length := binary.LittleEndian.Uint32(data)
			lastIdx := int(length + 4)
			if lastIdx > len(data) {
				break
			}
			list = append(list, decryptItem{
				data:         data[4:lastIdx],
				lengthInFile: lastIdx,
			})
			data = data[lastIdx:]
		}
	} else {
		for len(data) >= overheadOfCrypto {
			lengthBytes, em := udwCryptoAesCtrV7.Decrypt32(psk, data[:36])
			if em != `` {
				break
			}
			length := binary.LittleEndian.Uint32(lengthBytes)
			lastIdx := int(length + overheadOfCrypto)
			if lastIdx > len(data) {
				break
			}
			one, em := udwCryptoAesCtrV7.Decrypt32(psk, data[36:lastIdx])
			if em != `` {
				break
			}
			list = append(list, decryptItem{
				data:         one,
				lengthInFile: lastIdx,
			})
			data = data[lastIdx:]
		}
	}
	if len(data) > 0 {
		return list, errors.New(`ju2zgemkfk len(data)` + udwStrconv.GbFromInt(len(data)))
	}
	return list, nil
}
