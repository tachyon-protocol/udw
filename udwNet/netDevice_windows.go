// +build !js

package udwNet

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwNet/udwNetWindowsSyscall"
	"net"
	"os"
	"syscall"
	"unsafe"
)

func mustGetNetDeviceList() []*NetDevice {
	iaa, err := getAdapters()
	if err != nil {
		panic(err)
	}
	out := []*NetDevice{}
	for {
		if iaa == nil {
			break
		}

		thisDev := &NetDevice{
			index:     int(iaa.IfIndex),
			name:      stringFromUTF16Ptr(iaa.FriendlyName),
			firstIpv4: getIpv4FromIaa(iaa),
		}
		puni := iaa.FirstUnicastAddress
		for ; puni != nil; puni = puni.Next {
			sa, err := puni.Address.Sockaddr.Sockaddr()
			if err != nil {
				continue
			}
			switch sav := sa.(type) {
			case *syscall.SockaddrInet4:
				ip := make(net.IP, net.IPv4len)
				copy(ip, sav.Addr[:])
				thisDev.ipList = append(thisDev.ipList, ip)
				thisDev.ipNetList = append(thisDev.ipNetList, net.IPNet{
					IP:   ip,
					Mask: net.CIDRMask(int(puni.Address.SockaddrLength), 8*net.IPv4len),
				})
			case *syscall.SockaddrInet6:
				ip := make(net.IP, net.IPv6len)
				copy(ip, sav.Addr[:])
				thisDev.ipList = append(thisDev.ipList, ip)
				thisDev.ipNetList = append(thisDev.ipNetList, net.IPNet{
					IP:   ip,
					Mask: net.CIDRMask(int(puni.Address.SockaddrLength), 8*net.IPv6len),
				})
			}
		}
		out = append(out, thisDev)
		iaa = iaa.Next
	}
	return out
}

func stringFromCstringPrt(b *byte) string {
	buf := &bytes.Buffer{}
	for {
		if *b == 0 {
			break
		}
		buf.WriteByte(*b)
		b = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(b)) + 1))
	}
	return buf.String()
}

func stringFromUTF16Ptr(b *uint16) string {
	return syscall.UTF16ToString((*(*[10000]uint16)(unsafe.Pointer(b)))[:])
}

func getAdapters() (*udwNetWindowsSyscall.IpAdapterAddresses, error) {
	block := uint32(unsafe.Sizeof(udwNetWindowsSyscall.IpAdapterAddresses{}))

	size := uint32(15000)

	var addrs []udwNetWindowsSyscall.IpAdapterAddresses
	for {
		addrs = make([]udwNetWindowsSyscall.IpAdapterAddresses, size/block+1)
		err := udwNetWindowsSyscall.GetAdaptersAddresses(syscall.AF_UNSPEC, udwNetWindowsSyscall.GAA_FLAG_INCLUDE_PREFIX, 0, &addrs[0], &size)
		if err == nil {
			break
		}
		if err.(syscall.Errno) != syscall.ERROR_BUFFER_OVERFLOW {
			return nil, os.NewSyscallError("getadaptersaddresses", err)
		}
	}
	return &addrs[0], nil
}

func getIpv4FromIaa(iaa *udwNetWindowsSyscall.IpAdapterAddresses) net.IP {
	puni := iaa.FirstUnicastAddress
	for ; puni != nil; puni = puni.Next {
		sa, err := puni.Address.Sockaddr.Sockaddr()
		if err != nil {
			continue
		}
		switch sav := sa.(type) {
		case *syscall.SockaddrInet4:
			ip := make(net.IP, net.IPv4len)
			copy(ip, sav.Addr[:])
			return ip
		}
	}
	return nil
}
