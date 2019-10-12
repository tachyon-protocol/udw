// +build darwin,!ios darwin,!macAppStore

package udwTapTun

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwHex"
	"github.com/tachyon-protocol/udw/udwNet"
	"github.com/tachyon-protocol/udw/udwSys/udwSysAsyncFd"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

type TunTapObj struct {
	deviceType  DeviceType
	file        io.ReadWriteCloser
	name        string
	isUtun      bool
	writeBuf    []byte
	writeLocker sync.Mutex
}

func NewTap(ifName string) (ifce *TunTapObj, err error) {
	return newDeviceGeneric(ifName, DeviceTypeTap)
}

func NewTun(ifName string) (ifce *TunTapObj, err error) {
	return newDeviceGeneric(ifName, DeviceTypeTun)
}

func NewTunNoName() (ifce *TunTapObj, err error) {
	for i := 0; i < 255; i++ {
		fd, err := utunOpenHelper(i)
		if err != nil {
			if udwNet.IsResourceBusy(err) {

				continue
			}
			return nil, err
		}
		utunname := [20]byte{}
		utunname_len := uint32(20)
		err = SyscallGetSockopt(fd, SYSPROTO_CONTROL, UTUN_OPT_IFNAME, uintptr(unsafe.Pointer(&utunname[0])), &utunname_len)
		if err != nil {
			return nil, fmt.Errorf("[udwTapTun.NewTunNoName]Opening utun Error retrieving utun interface name %s", err.Error())
		}
		name := string(utunname[:int(utunname_len)-1])

		file, err := udwSysAsyncFd.FdToRwc(fd)
		if err != nil {
			return nil, fmt.Errorf("[udwTapTun.NewTunNoName]udwSysAsyncFd.FdToRwc %s", err.Error())
		}

		return &TunTapObj{
			deviceType: DeviceTypeTun,

			file:     file,
			name:     name,
			isUtun:   true,
			writeBuf: make([]byte, 2048),
		}, nil
	}
	return nil, getErrAllDeviceBusy()
}

const SYSPROTO_CONTROL = 2
const AF_SYS_CONTROL = 2

const CTLIOCGINFO = 0xc0644e03
const MAX_KCTL_NAME = 96

type ctl_info struct {
	ctl_id   uint32
	ctl_name [MAX_KCTL_NAME]byte
}
type sockaddr_ctl struct {
	sc_len      byte
	sc_family   byte
	ss_sysaddr  uint16
	sc_id       uint32
	sc_unit     uint32
	sc_reserved [5]uint32
}

const UTUN_CONTROL_NAME = "com.apple.net.utun_control\x00"
const UTUN_OPT_IFNAME = 2

func utunOpenHelper(num int) (fd int, err error) {
	fd, err = syscall.Socket(syscall.AF_SYSTEM, syscall.SOCK_DGRAM, SYSPROTO_CONTROL)
	if err != nil {
		return 0, fmt.Errorf("Opening utun socket(SYSPROTO_CONTROL) %s ", err.Error())
	}
	ctlInfo := ctl_info{}
	copy(ctlInfo.ctl_name[:], UTUN_CONTROL_NAME)
	err = SyscallIoctl(fd, CTLIOCGINFO, uintptr(unsafe.Pointer(&ctlInfo)))
	if err != nil {
		syscall.Close(fd)
		return 0, fmt.Errorf("Opening utun ioctl(CTLIOCGINFO) %s ", err.Error())
	}
	sa := sockaddr_ctl{
		sc_id:      ctlInfo.ctl_id,
		sc_family:  syscall.AF_SYSTEM,
		ss_sysaddr: AF_SYS_CONTROL,
		sc_unit:    uint32(num + 1),
	}
	sa.sc_len = byte(unsafe.Sizeof(sa))

	err = SyscallConnect(fd, uintptr(unsafe.Pointer(&sa)), uintptr(sa.sc_len))
	if err != nil {
		syscall.Close(fd)
		return 0, fmt.Errorf("Opening utun connect(AF_SYS_CONTROL) %s ", err.Error())
	}
	err = syscall.SetNonblock(fd, false)
	if err != nil {
		syscall.Close(fd)
		return 0, fmt.Errorf("Opening utun SetNonblock %s ", err.Error())
	}
	_, err = SyscallFcntl(fd, syscall.F_SETFD, syscall.FD_CLOEXEC)
	if err != nil {
		syscall.Close(fd)
		return 0, fmt.Errorf("Opening utun Fcntl %s ", err.Error())
	}
	return fd, nil
}

func newDeviceGeneric(ifName string, deviceType DeviceType) (ifce *TunTapObj, err error) {
	iifce := &TunTapObj{
		deviceType: deviceType,
		name:       ifName,
		writeBuf:   make([]byte, 2048),
	}

	devTypeString := deviceType.String()
	if ifName != "" {
		if !strings.HasPrefix(ifName, devTypeString) {
			return nil, fmt.Errorf("name should look like %s0 .", devTypeString)
		}
		iifce.file, err = os.OpenFile("/dev/"+ifName, os.O_RDWR, 0)
		if err != nil {
			return nil, err
		}
		return iifce, nil
	} else {

		for i := 0; i <= 15; i++ {
			iifce.name = devTypeString + strconv.Itoa(i)
			fmt.Println("Open tun ", iifce.name, "start")
			iifce.file, err = os.OpenFile("/dev/"+iifce.name, os.O_RDWR, 0)
			fmt.Println("Open tun ", iifce.name, "finish")

			if err == nil {
				return iifce, nil
			}
			if err != nil {
				if strings.Contains(err.Error(), "resource busy") {
					continue
				}
				return nil, err
			}
		}
		return nil, getErrAllDeviceBusy()
	}
}

func (ifce *TunTapObj) GetDeviceType() DeviceType {
	return ifce.deviceType
}

func (ifce *TunTapObj) Name() string {
	return ifce.name
}

func (ifce *TunTapObj) Write(buf []byte) (n int, err error) {
	ifce.writeLocker.Lock()
	inLen := len(buf)
	firstByte := buf[0]
	if ifce.isUtun {
		if inLen+4 > len(ifce.writeBuf) {
			ifce.writeBuf = make([]byte, inLen+1000)
		}
		ifce.writeBuf[0] = 0
		ifce.writeBuf[1] = 0
		ifce.writeBuf[2] = 0
		if firstByte&0xf0 == 0x40 {
			ifce.writeBuf[3] = 2
		} else if firstByte&0xf0 == 0x60 {
			ifce.writeBuf[3] = 0x1e
		} else {
			ifce.writeLocker.Unlock()
			panic(fmt.Errorf("unexpect first byte %d,%s", firstByte, udwHex.EncodeBytesToString(buf)))
		}
		copy(ifce.writeBuf[4:], buf)
		buf = ifce.writeBuf[:4+inLen]
	}
	_, err = ifce.file.Write(buf)
	ifce.writeLocker.Unlock()
	if err != nil {
		return 0, err
	}
	return inLen, nil
}

func (ifce *TunTapObj) Read(p []byte) (nr int, err error) {
	nr, err = ifce.file.Read(p)
	if err != nil {
		return
	}
	if ifce.isUtun {
		copy(p[:nr-4], p[4:nr])
		return nr - 4, nil
	} else {
		return nr, nil
	}
}
func (ifce *TunTapObj) Close() (err error) {
	return ifce.file.Close()
}

func (ifce *TunTapObj) StartWriteBufferThread(printErrorMsg string) {
}

func (ifce *TunTapObj) WriteWithBuffer(buf []byte) {
	ifce.Write(buf)
}
