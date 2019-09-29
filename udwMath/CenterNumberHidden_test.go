package udwMath

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestCenterString(ot *testing.T) {
	test1 := HiddenCenterString("32112")
	udwTest.Equal(test1, "32**2")
	test2 := HiddenCenterString("12345678901")
	udwTest.Equal(test2, "1234****901")
	test3 := HiddenCenterString("123")
	udwTest.Equal(test3, "1**")
}
