package udwTapTun

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwBytes"
	"io"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

type TunTapObj struct {
	deviceType   DeviceType
	file         *os.File
	fd           uintptr
	name         string
	tunWriteChan chan *udwBytes.BufWriter
}

func NewTap(ifName string) (ifce *TunTapObj, err error) {
	return newTunTapObj(true)
}

func NewTun(ifName string) (ifce *TunTapObj, err error) {
	return newTunTapObj(false)
}

func NewTunNoName() (ifce *TunTapObj, err error) {
	return NewTun("")
}

func newTunTapObj(isTap bool) (ifce *TunTapObj, err error) {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	fd := file.Fd()
	flag := uint16(cIFF_TUN | cIFF_NO_PI)
	if isTap {
		flag = cIFF_TAP | cIFF_NO_PI
	}
	name, err := createInterface(fd, "", flag)
	if err != nil {
		file.Close()
		return nil, err
	}
	ifce = &TunTapObj{deviceType: DeviceTypeTun, file: file, name: name, fd: fd}
	return
}

func (ifce *TunTapObj) GetDeviceType() DeviceType {
	return ifce.deviceType
}

func (ifce *TunTapObj) Name() string {
	return ifce.name
}

func (ifce *TunTapObj) Write(p []byte) (n int, err error) {

	return ifce.file.Write(p)
}

func (ifce *TunTapObj) WriteFastNoLock(p []byte) (err error) {
	num, _, errno := syscall.Syscall(syscall.SYS_WRITE, ifce.fd, uintptr(unsafe.Pointer(&p[0])), uintptr(len(p)))
	numi := int(num)
	if numi == -1 {
		err = errno
		if errno == 0 {
			panic("[iInterface.Write] numi==-1 errno==0")
		}
	}
	return err
}

func (ifce *TunTapObj) Read(p []byte) (n int, err error) {
	return ifce.file.Read(p)
}
func (ifce *TunTapObj) Close() (err error) {
	return ifce.file.Close()
}

const (
	cIFF_TUN   = 0x0001
	cIFF_TAP   = 0x0002
	cIFF_NO_PI = 0x1000
)

type ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func createInterface(fd uintptr, ifName string, flags uint16) (createdIFName string, err error) {
	var req ifReq
	req.Flags = flags
	copy(req.Name[:], ifName)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		err = errno
		return
	}
	createdIFName = strings.Trim(string(req.Name[:]), "\x00")
	return
}

var gTunTcpWritePool udwBytes.BufWriterPool

func (ifce *TunTapObj) StartWriteBufferThread(printErrorMsg string) {
	ifce.tunWriteChan = make(chan *udwBytes.BufWriter, 4096)
	go func() {
		for {
			bw := <-ifce.tunWriteChan
			err := ifce.WriteFastNoLock(bw.GetBytes())
			gTunTcpWritePool.Put(bw)
			if err != nil && printErrorMsg != "" {
				fmt.Println("[TunTapObj.StartWriteTunBuffer] 3fukb7mxv9", printErrorMsg, err)
			}
		}
	}()
}

func (ifce *TunTapObj) WriteWithBuffer(buf []byte) {
	ifce.tunWriteChan <- gTunTcpWritePool.GetAndCloneFromByteSlice(buf)
}

func (ifce *TunTapObj) writev(v *[][]byte, bufIovec []syscall.Iovec) (int64, error) {
	iovecs := bufIovec[:0]
	maxVec := 1024

	var n int64
	var err error
	for len(*v) > 0 {
		iovecs = iovecs[:0]
		for _, chunk := range *v {
			if len(chunk) == 0 {
				continue
			}
			iovecs = append(iovecs, syscall.Iovec{Base: &chunk[0]})
			iovecs[len(iovecs)-1].SetLen(len(chunk))
			if len(iovecs) == maxVec {
				break
			}
		}
		if len(iovecs) == 0 {
			break
		}

		wrote, _, e0 := syscall.Syscall(syscall.SYS_WRITEV,
			uintptr(ifce.fd),
			uintptr(unsafe.Pointer(&iovecs[0])),
			uintptr(len(iovecs)))
		if wrote == ^uintptr(0) {
			wrote = 0
		}
		n += int64(wrote)
		consume(v, int64(wrote))

		if e0 != 0 {
			err = syscall.Errno(e0)
		}
		if err != nil {
			break
		}
		if n == 0 {
			err = io.ErrUnexpectedEOF
			break
		}
	}
	return n, err
}

func consume(v *[][]byte, n int64) {
	for len(*v) > 0 {
		ln0 := int64(len((*v)[0]))
		if ln0 > n {
			(*v)[0] = (*v)[0][n:]
			return
		}
		n -= ln0
		*v = (*v)[1:]
	}
}
