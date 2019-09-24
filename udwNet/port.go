package udwNet

import (
	"errors"
	"math"
	"net"
	"strconv"
)

func GetProtocolNotSupportPortError() error {
	return errors.New("Protocol Not Support Port")
}

func PortFromNetAddr(addr net.Addr) (int, error) {
	switch saddr := addr.(type) {
	case *net.TCPAddr:
		return saddr.Port, nil
	case *net.UDPAddr:
		return saddr.Port, nil
	case *net.IPAddr:
		return 0, GetProtocolNotSupportPortError()
	case *net.UnixAddr:
		return -1, GetProtocolNotSupportPortError()
	}
	return PortFromAddrString(addr.String())
}

func PortFromAddrString(addr string) (int, error) {
	_, portS, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}
	portI, err := strconv.Atoi(portS)
	if err != nil {
		return 0, err
	}
	return portI, nil
}

func JoinHostPortInt(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}

func MustGetHostFromAddr(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}
	return host
}

func JoinHostPortNoAlloc(buf []byte, ip string, port int) []byte {
	if len(buf) < 21 {
		buf = make([]byte, 21)
	}
	n := copy(buf, ip)
	buf[n] = byte(':')
	i := 4
	_port := port
	for {
		if i <= 0 {
			break
		}

		d := int(math.Pow10(i))
		p := _port / d
		i--
		if p == 0 {
			continue
		}
		_port = _port - d*p
		n++
		buf[n] = uint8(p) + 48
	}
	n++
	buf[n] = uint8(_port) + 48
	return buf[:n+1]
}
