package udwTapTun

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwNet"
	"io"
	"net"
)

type DeviceType string

func (s DeviceType) String() string {
	return string(s)
}

var DeviceTypeTap DeviceType = "tap"
var DeviceTypeTun DeviceType = "tun"

func getErrAllDeviceBusy() error {
	return errors.New("tun/tap: all dev is busy.")
}

type TunTapInterface interface {
	io.ReadWriteCloser
	GetDeviceType() DeviceType
	Name() string
}

type CreateIpv4TunContext struct {
	SrcIp net.IP

	DstIp net.IP

	FirstIp net.IP

	DhcpServerIp net.IP

	Mask net.IPMask

	Mtu int

	ReturnTun io.ReadWriteCloser
	ReturnDev *udwNet.NetDevice
}

func MustCreateIpv4Tun(ctx *CreateIpv4TunContext) {
	mustCreateIpv4Tun(ctx)
}

func CreateIpv4Tun(ctx *CreateIpv4TunContext) (err error) {
	return udwErr.PanicToErrorAndLog(func() {
		MustCreateIpv4Tun(ctx)
	})
}
