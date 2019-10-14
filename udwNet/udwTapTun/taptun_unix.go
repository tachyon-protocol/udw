// +build darwin,!ios linux darwin,!macAppStore

package udwTapTun

import (
	"github.com/tachyon-protocol/udw/udwNet"
	"net"
)

func mustCreateIpv4Tun(ctx *CreateIpv4TunContext) {
	tunNamed, err := NewTunNoName()
	if err != nil {
		panic(err)
	}
	err = SetP2PIpAndUp(SetP2PIpRequest{
		IfaceName: tunNamed.Name(),
		SrcIp:     ctx.SrcIp.To4(),
		DstIp:     ctx.DstIp.To4(),
		Mtu:       ctx.Mtu,
		Mask:      net.CIDRMask(31, 32),
	})
	if err != nil {
		tunNamed.Close()
		panic(err)
	}
	ctx.ReturnTun = tunNamed
	ctx.ReturnDev, err = udwNet.GetNetDeviceByName(tunNamed.Name())
	if err != nil {
		tunNamed.Close()
		panic(err)
	}

	return
}
