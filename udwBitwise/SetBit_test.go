package udwBitwise

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestUint16SetBit(t *testing.T) {
	udwTest.Equal(Uint16SetBit0(1, 0), uint16(0))
	udwTest.Equal(Uint16SetBit0(3, 0), uint16(2))

	udwTest.Equal(Uint16SetBit1(0, 0), uint16(1))
	udwTest.Equal(Uint16SetBit1(0, 1), uint16(2))
	udwTest.Equal(Uint16SetBit1(1, 1), uint16(3))

}

func TestUint32GetBitTo1Or0Uint8(t *testing.T) {
	udwTest.Equal(Uint32GetBitTo1Or0Uint8(1, 0), uint8(1))
	udwTest.Equal(Uint32GetBitTo1Or0Uint8(1, 1), uint8(0))
	udwTest.Equal(Uint32GetBitTo1Or0Uint8(2, 0), uint8(0))
}
