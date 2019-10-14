// +build windows

package udwTapTun

import (
	"encoding/binary"
	"fmt"
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwLog"
	"github.com/tachyon-protocol/udw/udwNet"
	"github.com/tachyon-protocol/udw/udwSys/udwSysAsyncFd"
	"strconv"
	"syscall"
	"unsafe"
)

const TAP_WIN_IOCTL_GET_VERSION = 2228232
const TAP_WIN_IOCTL_GET_MTU = 2228236
const TAP_WIN_IOCTL_CONFIG_TUN = 2228264
const TAP_WIN_IOCTL_CONFIG_POINT_TO_POINT = 2228244
const TAP_WIN_IOCTL_SET_MEDIA_STATUS = 2228248
const TAP_WIN_IOCTL_CONFIG_DHCP_MASQ = 2228252
const FILE_FLAG_NO_BUFFERING = 0x20000000
const FILE_FLAG_WRITE_THROUGH = 0x80000000

func mustCreateIpv4Tun(ctx *CreateIpv4TunContext) {
	deviceGuid, actualName := GetDeviceGuidAndActualName()
	path := `\\.\Global\` + deviceGuid + ".tap"
	pathP, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		panic(err)
	}
	h, err := syscall.CreateFile(pathP, syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		0, &syscall.SecurityAttributes{}, syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_SYSTEM|syscall.FILE_FLAG_OVERLAPPED, 0)
	if err != nil {
		panic(err)
	}
	versionInfo := [3]uint32{}
	returned := uint32(0)
	err = syscall.DeviceIoControl(h, TAP_WIN_IOCTL_GET_VERSION,
		(*byte)(unsafe.Pointer(&versionInfo[0])), 12,
		(*byte)(unsafe.Pointer(&versionInfo[0])), 12,
		&returned, nil)
	if err != nil {
		panic(err)
	}

	if versionInfo[0] != 9 && versionInfo[1] != 21 {
		panic(fmt.Errorf("[NewTunNoName] unknow tap-win version %d.%d.%d", versionInfo[0], versionInfo[1], versionInfo[2]))
	}
	err = udwCmd.CmdSlice([]string{
		"netsh",
		"interface",
		"ipv4",
		"set",
		"address",
		actualName,
		"dhcp",
	}).RunAndNotExitStatusCheck()
	if err != nil {
		panic(err)
	}

	ipUint32 := binary.LittleEndian.Uint32([]byte(ctx.SrcIp.To4()))

	maskUint32 := binary.LittleEndian.Uint32([]byte(ctx.Mask))

	firstIpUint32 := binary.LittleEndian.Uint32([]byte(ctx.FirstIp.To4()))

	dhcpServerIpUint32 := binary.LittleEndian.Uint32([]byte(ctx.DhcpServerIp.To4()))

	tunStatus := [3]uint32{
		ipUint32,
		firstIpUint32,
		maskUint32,
	}
	err = syscall.DeviceIoControl(h, TAP_WIN_IOCTL_CONFIG_TUN,
		(*byte)(unsafe.Pointer(&tunStatus[0])), 12,
		(*byte)(unsafe.Pointer(&tunStatus[0])), 12,
		&returned, nil)
	if err != nil {
		panic(fmt.Errorf("[mustCreateIpv4Tun] TAP_WIN_IOCTL_CONFIG_TUN %s", err))
	}

	dhcpServerConfig := [4]uint32{
		ipUint32,
		maskUint32,

		dhcpServerIpUint32,
		uint32(31536000),
	}
	err = syscall.DeviceIoControl(h, TAP_WIN_IOCTL_CONFIG_DHCP_MASQ,
		(*byte)(unsafe.Pointer(&dhcpServerConfig)), 16,
		(*byte)(unsafe.Pointer(&dhcpServerConfig)), 16,
		&returned, nil)
	if err != nil {
		panic(fmt.Errorf("[mustCreateIpv4Tun] TAP_WIN_IOCTL_CONFIG_DHCP_MASQ %s", err))
	}
	setConnected := uint32(1)
	err = syscall.DeviceIoControl(h, TAP_WIN_IOCTL_SET_MEDIA_STATUS,
		(*byte)(unsafe.Pointer(&setConnected)), 4,
		(*byte)(unsafe.Pointer(&setConnected)), 4,
		&returned, nil)
	if err != nil {
		panic(fmt.Errorf("[mustCreateIpv4Tun] TAP_WIN_IOCTL_SET_MEDIA_STATUS %s", err))
	}

	udwCmd.CmdSlice([]string{
		"C:\\WINDOWS\\system32\\netsh.exe",
		"interface",
		"ip",
		"delete",
		"arp",
		actualName,
	}).RunAndNotExitStatusCheck()

	go func() {

		err = DHCPRenew(deviceGuid)
		if err != nil {
			udwLog.Log("error", "DHCPRenew fail", err.Error())
		}
	}()

	udwCmd.CmdSlice([]string{"netsh", "interface", "ipv4", "add", "dnsserver", actualName, "address=8.8.8.8", "index=1"}).
		MustRunAndNotExitStatusCheck()
	udwCmd.CmdSlice([]string{"netsh", "interface", "ipv4", "add", "dnsserver", actualName, "address=8.8.4.4", "index=2"}).
		MustRunAndNotExitStatusCheck()
	udwCmd.CmdSlice([]string{"netsh", "interface", "ipv4", "set", "subinterface", actualName, "mtu=" + strconv.Itoa(ctx.Mtu), "store=persistent"}).
		MustRunAndNotExitStatusCheck()
	udwCmd.MustRun("ipconfig /flushdns")

	tun, err := udwSysAsyncFd.FdToRwc(int(h))
	if err != nil {
		panic(err)
	}
	vpnDevIndex := GetAdapterIndexByGuid(deviceGuid)
	dev, err := udwNet.GetNetDeviceByIndex(int(vpnDevIndex))
	if err != nil {
		panic(err)
	}
	ctx.ReturnTun = tun
	ctx.ReturnDev = dev
}
