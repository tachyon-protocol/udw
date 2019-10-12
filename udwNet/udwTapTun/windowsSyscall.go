// +build windows

package udwTapTun

import "syscall"

var (
	modiphlpapi      = syscall.NewLazyDLL("iphlpapi.dll")
	ipReleaseAddress = modiphlpapi.NewProc("IpReleaseAddress")
	ipRenewAddress   = modiphlpapi.NewProc("IpRenewAddress")
	getAdapterIndex  = modiphlpapi.NewProc("GetAdapterIndex")
	getInterfaceInfo = modiphlpapi.NewProc("GetInterfaceInfo")
	flushIpNetTable  = modiphlpapi.NewProc("FlushIpNetTable")
)
