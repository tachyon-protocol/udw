package udwSys

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwPlatform"
)

func IsIpForwardOn() bool {
	if !udwPlatform.IsLinux() {
		panic("[IsIpForwardOn] only support linux now")
	}
	b := udwFile.MustReadFile("/proc/sys/net/ipv4/ip_forward")
	if bytes.Contains(b, []byte{'0'}) {
		return false
	}
	if bytes.Contains(b, []byte{'1'}) {
		return true
	}
	panic(fmt.Errorf("[IsIpForwardOn] unable to understand info in /proc/sys/net/ipv4/ip_forward %#v", b))
}

func SetIpForwardOn() {
	if !udwPlatform.IsLinux() {
		panic("[SetIpForwardOn] only support linux now")
	}
	udwFile.MustWriteFile("/proc/sys/net/ipv4/ip_forward", []byte("1"))

	if !bytes.Contains(udwFile.MustReadFile("/etc/sysctl.conf"), []byte("\nnet.ipv4.ip_forward = 1")) {
		udwFile.MustAppendFile("/etc/sysctl.conf", []byte("\nnet.ipv4.ip_forward = 1"))
	}
}
