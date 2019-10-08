// +build ios macAppStore

package udwNet

/*
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"syscall"
	"unsafe"
)

func NewSupportIpv6OnlyDialer(oldDialer Dialer) Dialer {
	return func(network, address string) (net.Conn, error) {

		conn1, firstError := oldDialer(network, address)
		if firstError == nil {
			return conn1, nil
		}

		if IsNetworkIsUnreachable(firstError) {
			ipS, port, err := net.SplitHostPort(address)
			if err != nil {

				return nil, fmt.Errorf("[NewSupportIpv6OnlyDialer] net.SplitHostPort %s", err.Error())
			}
			ipObjList, _, err := Ipv6GetaddrinfoWithCgo(ipS)
			if err != nil {
				return nil, fmt.Errorf("[NewSupportIpv6OnlyDialer] fail %s", err.Error())
			}
			if len(ipObjList) == 0 {
				return nil, errors.New("[NewSupportIpv6OnlyDialer] len(ipObjList)==0")
			}
			address = net.JoinHostPort(ipObjList[0].String(), port)

			return oldDialer(network, address)
		}
		return nil, firstError
	}
}

type addrinfoErrno int

func (eai addrinfoErrno) Error() string { return C.GoString(C.gai_strerror(C.int(eai))) }

const cgoAddrInfoFlags = (C.AI_CANONNAME | C.AI_V4MAPPED | C.AI_ALL) & C.AI_MASK

func Ipv6GetaddrinfoWithCgo(name string) (addrs []net.IPAddr, cname string, err error) {
	var hints C.struct_addrinfo

	hints.ai_family = C.PF_UNSPEC
	hints.ai_socktype = C.SOCK_STREAM
	hints.ai_flags = C.AI_DEFAULT

	h := C.CString(name)
	defer C.free(unsafe.Pointer(h))
	var res *C.struct_addrinfo
	httpC := C.CString("http")
	gerrno, err := C.getaddrinfo(h, httpC, &hints, &res)
	C.free(unsafe.Pointer(httpC))
	if gerrno != 0 {
		switch gerrno {
		case C.EAI_SYSTEM:
			if err == nil {

				err = syscall.EMFILE
			}
		case C.EAI_NONAME:
			err = errors.New("no such host")
		default:
			err = addrinfoErrno(gerrno)
		}
		return nil, "", &net.DNSError{Err: err.Error(), Name: name}
	}
	defer C.freeaddrinfo(res)

	if res != nil {
		cname = C.GoString(res.ai_canonname)
		if cname == "" {
			cname = name
		}
		if len(cname) > 0 && cname[len(cname)-1] != '.' {
			cname += "."
		}
	}
	for r := res; r != nil; r = r.ai_next {

		switch r.ai_family {

		case C.AF_INET:
			sa := (*syscall.RawSockaddrInet4)(unsafe.Pointer(r.ai_addr))
			addr := net.IPAddr{IP: IpClone(sa.Addr[:])}
			addrs = append(addrs, addr)
		case C.AF_INET6:
			sa := (*syscall.RawSockaddrInet6)(unsafe.Pointer(r.ai_addr))
			addr := net.IPAddr{IP: IpClone(sa.Addr[:]), Zone: zoneToString(int(sa.Scope_id))}
			addrs = append(addrs, addr)
		}
	}
	return addrs, cname, nil
}

func zoneToString(zone int) string {
	if zone == 0 {
		return ""
	}
	if ifi, err := net.InterfaceByIndex(zone); err == nil {
		return ifi.Name
	}
	return strconv.FormatUint(uint64(zone), 10)
}
