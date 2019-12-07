package udwIpToCountryV2

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"net"
	"runtime/debug"
	"testing"
)

func TestGetCountryIsoCode(ot *testing.T) {
	for _, cas := range []struct {
		ip   string
		code string
	}{
		{"180.97.33.107", "CN"},
		{"127.0.0.1", ""},
		{"10.1.1.1", ""},
		{"173.194.127.50", "US"},
		{"2404:6800:4005:801::200e", "AU"},
		{"43.250.12.38", "CN"},
		{"35.229.165.10", "US"},
		{"61.14.4.0", "HK"},
		{"122.8.59.125", "PK"},
		{"122.8.32.0", ""},
		{"255.255.255.255", ""},
	} {
		udwTest.Equal(MustGetCountryIsoCode(net.ParseIP(cas.ip)), cas.code, cas.ip)
		udwTest.Equal(MustGetCountryIsoCodeByString(cas.ip), cas.code, cas.ip)
	}
	udwTest.Equal(MustGetCountryIsoCodeByString(""), "")
	udwTest.Equal(MustGetCountryIsoCodeByString(" "), "")
	udwTest.Equal(MustGetCountryIsoCode(nil), "")
}

func TestBenchMustGetCountryIsoCodeByString(ot *testing.T) {
	debug.FreeOSMemory()
	MustGetCountryIsoCodeByString("173.194.127.50")

	udwTest.BenchmarkWithRepeatNum(1e6, func() {
		MustGetCountryIsoCodeByString("173.194.127.50")
	})

}
