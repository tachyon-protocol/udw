package udwIpToCountryV2

import (
	"github.com/tachyon-protocol/udw/udwIpToCountryV2/udwIpCountryV2Map"
	"net"
	"sync"
	"sync/atomic"
	"unsafe"
)

func EnsureInit() {
	gEnsureInitOnce.Do(func() {

		thisReader := getGeoip2Reader()
		SetReader(thisReader)
	})
}

func MustGetCountryIsoCode(ip net.IP) (code string) {
	EnsureInit()
	return GetReader().MustGetCountryIsoCode(ip)
}

func MustGetCountryIsoCodeByString(ip string) (code string) {
	EnsureInit()
	return GetReader().MustGetCountryIsoCodeByString(ip)
}

func GetCountryCodeList() []string {
	EnsureInit()
	return GetReader().GetAllCountryCode()
}

var gReaderL2 unsafe.Pointer

func SetReader(r *udwIpCountryV2Map.Reader) {
	if r == nil {
		panic("[udwIpToCountryV2.SetReader] r==nil")
	}
	r.CompatibleToOldVersion()
	atomic.StorePointer(&gReaderL2, unsafe.Pointer(r))

}

func GetReader() *udwIpCountryV2Map.Reader {
	thisReaderL1 := atomic.LoadPointer(&gReaderL2)
	return (*udwIpCountryV2Map.Reader)(thisReaderL1)

}

var gEnsureInitOnce sync.Once
