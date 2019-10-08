package udwNet

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestMustLookupDomainInAddrString(ot *testing.T) {
	out := MustLookupDomainInAddrString("www.baidu.com:80")
	ip, port := MustSplitIpPort(out)
	udwTest.Equal(port, uint16(80))
	udwTest.Ok(ip != nil)
	out = MustLookupDomainInAddrString("127.0.0.1:80")
	udwTest.Equal(out, "127.0.0.1:80")
}
