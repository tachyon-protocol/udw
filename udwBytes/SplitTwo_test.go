package udwBytes

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestSplitTwoBetweenFirst(t *testing.T) {
	var b []byte
	var a []byte
	b, a = SplitTwoBetweenFirst([]byte("123"), []byte("2"))
	udwTest.Equal(b, []byte("1"))
	udwTest.Equal(a, []byte("3"))
}
