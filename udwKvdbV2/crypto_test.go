package udwKvdbV2

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwCryptoAesCtrV7"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestEncrypto(t *testing.T) {
	psk := udwCryptoAesCtrV7.Get32PskSha3FromString(`psk`)
	data := cryptoBlock(psk, []byte(`123456`))
	udwTest.Ok(overheadOfCrypto+6 == len(data))
	list, _ := decryptoBlockList(psk, data)
	udwTest.Ok(len(list) == 1 && string(list[0].data) == `123456`)
	udwTest.Ok(overheadOfCrypto+6 == list[0].lengthInFile)

	data = append(data, cryptoBlock(psk, []byte(`1`))...)
	udwTest.Ok(overheadOfCrypto*2+7 == len(data))
	list, _ = decryptoBlockList(psk, data)
	udwTest.Ok(len(list) == 2, len(list))

	list, _ = decryptoBlockList(psk, data[:len(data)-1])
	udwTest.Ok(len(list) == 1)

	tmp := udwBytes.Clone(data)
	tmp[len(tmp)-1]--
	list, _ = decryptoBlockList(psk, tmp)
	udwTest.Ok(len(list) == 1)
	tmp[list[0].lengthInFile+1]++
	list, _ = decryptoBlockList(psk, tmp)
	udwTest.Ok(len(list) == 1)
}

func TestCrypto(t *testing.T) {
	var psk *[32]byte = nil
	var data []byte

	for i := 0; i < 3; i++ {
		blk := bytes.Repeat([]byte{1}, i*100*udwStrconv.MB)
		data = append(data, cryptoBlock(psk, blk)...)
	}
	list, _ := decryptoBlockList(psk, data)
	udwTest.Ok(len(list) == 3)
	list, _ = decryptoBlockList(psk, data[:len(data)-1])
	udwTest.Ok(len(list) == 2)
}
