package udwCryptoSha3

import (
	"crypto/hmac"
	"encoding/hex"
	"fmt"
	"github.com/tachyon-protocol/udw/udwTest"
	"math"
	"strconv"
	"testing"
)

func TestSha3(t *testing.T) {
	udwTest.Equal(AlphaNumByString("abc", 15), "9f28mufw8dts775")
	udwTest.Equal(len(AlphaNumByString("abc", 15)), 15)
	udwTest.Equal(AlphaNumByString("abc", 15), AlphaNumByString("abc", 15))
	udwTest.Ok(AlphaNumByString("abc", 15) != AlphaNumByString("abcd", 15))
	for i := 0; i < 100; i++ {
		udwTest.Equal(len(AlphaNumByString("abcf", i)), i)
	}

	udwTest.Equal(IntnByString("abc", 15), 13)
	udwTest.Equal(IntnByString("abcd", 15), 5)

}

func ExampleFloat64ByString() {
	rangeMap := map[float64]int{}
	for i := 0; i < 800; i++ {
		out := Float64ByString(strconv.Itoa(i))
		thisRange := out - math.Mod(out, 1.0/16)
		rangeMap[thisRange]++
	}
	fmt.Println(len(rangeMap))
	for thisRange, size := range rangeMap {
		fmt.Println(thisRange, size)
	}
}

func ExampleIntnByString() {
	rangeMap := map[int]int{}
	for i := 0; i < 1000; i++ {
		out := IntnByString(strconv.Itoa(i), 5)
		rangeMap[out]++
	}
	fmt.Println(len(rangeMap))
	for thisRange, size := range rangeMap {
		fmt.Println(thisRange, size)
	}
}

func TestSha3512(ot *testing.T) {
	type t_cas struct {
		in     string
		outHex string
	}
	data := []t_cas{
		{"", "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26"},
		{"123", "48c8947f69c054a5caa934674ce8881d02bb18fb59d5a63eeaddff735b0e9801e87294783281ae49fc8287a0fd86779b27d7972d3e84f0fa0d826d7cb67dfefc"},
		{"\x00", "7127aab211f82a18d06cf7578ff49d5089017944139aa60d8bee057811a15fb55a53887600a3eceba004de51105139f32506fe5b53e1913bfa6b32e716fe97da"},
	}
	for _, cas := range data {
		udwTest.Equal(Sha3512ToHexStringFromString(cas.in), cas.outHex)
	}
}

func TestHmacSha3512(ot *testing.T) {
	h := hmac.New(New512, []byte("1"))
	h.Write([]byte("1"))
	udwTest.Equal(hex.EncodeToString(h.Sum(nil)), "6c033a2e0d093408050e39de3f2674ac00b4810c54f006336f843a593792ff90be34ae3d5b742ecedec86682b9a219488fc38371e4d76f0ebe82acc9a95de918")
}
