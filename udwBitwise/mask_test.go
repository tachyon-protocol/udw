package udwBitwise

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestUint16SetPartWithMask(t *testing.T) {
	udwTest.Equal(Uint16SetPartWithMask(0, 1, 1), uint16(1))
	udwTest.Equal(Uint16SetPartWithMask(2, 1, 1), uint16(3))
	udwTest.Equal(Uint16SetPartWithMask(4, 3, 3), uint16(7))

	udwTest.Equal(Uint16SetPartWithMask(3, 2, 3), uint16(2))
}
