package udwNet

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestGetIpStringNoPortOrInput(t *testing.T) {
	udwTest.Equal(GetIpStringNoPortOrInput(":"), "")
	udwTest.Equal(GetIpStringNoPortOrInput("127.0.0.1:1"), "127.0.0.1")
	udwTest.Equal(GetIpStringNoPortOrInput("127.0.0.1"), "127.0.0.1")
}
