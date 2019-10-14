// +build windows

package udwTapTun

import (
	"encoding/binary"
	"fmt"
	"github.com/tachyon-protocol/udw/udwErr"
	"syscall"
	"unsafe"
)

func GetAdapterNameByGuid(guid string) (adapterName [256]uint16) {
	_name, err := syscall.UTF16FromString("\\DEVICE\\TCPIP_" + guid)
	udwErr.PanicIfError(err)
	for i, b := range _name {
		adapterName[i] = b
	}
	adapterName[255] = 0
	return adapterName
}

func GetAdapterIndexByGuid(guid string) (index uint32) {
	name := GetAdapterNameByGuid(guid)
	r1, _, _ := syscall.Syscall(getAdapterIndex.Addr(), 2, uintptr(unsafe.Pointer(&name)), uintptr(unsafe.Pointer(&index)), 0)
	if r1 != 0 {
		fmt.Println("GetAdapterIndexByGuid Error")
		panic(syscall.Errno(r1))
	}
	return index
}

type IP_ADAPTER_INDEX_MAP struct {
	Index uint32
	Name  [syscall.MAX_ADAPTER_NAME_LENGTH]uint16
}

type IP_INTERFACE_INFO struct {
	NumAdapters uint32
	Adapter     []IP_ADAPTER_INDEX_MAP
}

func GetIP_ADAPTER_INDEX_MAPByAdapterIndex(index uint32) *IP_ADAPTER_INDEX_MAP {

	size := uint32(0)
	buf := make([]byte, 4096)
	r1, _, _ := syscall.Syscall(getInterfaceInfo.Addr(), 2, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)), 0)
	if syscall.Errno(r1) != syscall.ERROR_INSUFFICIENT_BUFFER {
		fmt.Println("GetInterfaceInfo Error")
		panic(syscall.Errno(r1))
	}
	buf = make([]byte, size)
	r1, _, _ = syscall.Syscall(getInterfaceInfo.Addr(), 2, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)), 0)
	if r1 != 0 {
		fmt.Println("GetInterfaceInfo Error")
		panic(syscall.Errno(r1))
	}
	num := binary.LittleEndian.Uint32(buf[:4])
	fmt.Println("binary.LittleEndian.Uint32(buf[:4])", num)
	listBuf := buf[4:]
	list := make([]IP_ADAPTER_INDEX_MAP, int(num))
	copy(list, *((*[]IP_ADAPTER_INDEX_MAP)(unsafe.Pointer(&listBuf))))
	info := &IP_INTERFACE_INFO{
		NumAdapters: num,
		Adapter:     list,
	}
	for _, i := range info.Adapter {
		n := []uint16{}
		for _, b := range i.Name {
			n = append(n, b)
		}
		if i.Index == uint32(32) {
			fmt.Println("bingo")
		}
		fmt.Println(i.Index, syscall.UTF16ToString(n))
	}
	return nil
}

func NewIP_ADAPTER_INDEX_MAP(guid string) *IP_ADAPTER_INDEX_MAP {
	return &IP_ADAPTER_INDEX_MAP{
		Index: GetAdapterIndexByGuid(guid),
		Name:  GetAdapterNameByGuid(guid),
	}
}

func DHCPRelease(guid string) {
	r1, _, _ := syscall.Syscall(ipReleaseAddress.Addr(), 1, uintptr(unsafe.Pointer(NewIP_ADAPTER_INDEX_MAP(guid))), 0, 0)
	if r1 != 0 {
		fmt.Println("DHCPRelease Error")
		panic(syscall.Errno(r1))
	}
}

func DHCPRenew(guid string) (err error) {
	r1, _, _ := syscall.Syscall(ipRenewAddress.Addr(), 1, uintptr(unsafe.Pointer(NewIP_ADAPTER_INDEX_MAP(guid))), 0, 0)
	if r1 != 0 {
		return syscall.Errno(r1)
	}
	return nil
}

func FlushIpNetTable(index uint32) {
	fmt.Println(index)
	r1, _, _ := syscall.Syscall(flushIpNetTable.Addr(), 1, uintptr(unsafe.Pointer(&index)), 0, 0)
	if r1 != 0 {
		fmt.Println("FlushIpNetTable Error")
		panic(syscall.Errno(r1))
	}
}
