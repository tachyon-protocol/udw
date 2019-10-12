package udwTapTun

import (
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestTun(ot *testing.T) {
	tun, err := NewTun("")
	if err != nil {
		ot.Skip("you need root permission to run this test.")
		return
	}
	udwTest.Equal(err, nil)
	defer tun.Close()
	udwTest.Equal(tun.GetDeviceType(), DeviceTypeTun)

	cmd := udwCmd.CmdString("ping 10.209.34.2").GetExecCmd()
	err = cmd.Start()
	udwTest.Equal(err, nil)
	defer cmd.Process.Kill()

	buf := make([]byte, 4096)
	n, err := tun.Read(buf)
	udwTest.Equal(err, nil)
	udwTest.Ok(n > 0)

}

func TestTap(ot *testing.T) {
	tap, err := NewTap("")
	if err != nil {
		ot.Skip("you need root permission to run this test.")
		return
	}
	udwTest.Equal(err, nil)
	defer tap.Close()
	udwTest.Equal(tap.GetDeviceType(), DeviceTypeTap)

	err = udwCmd.CmdString("ifconfig " + tap.Name() + " 10.209.34.1 up").GetExecCmd().Run()
	udwTest.Equal(err, nil)

}
