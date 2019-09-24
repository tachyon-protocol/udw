package udwNet

import (
	"encoding/binary"
	"net"
)

func ParseIpv4ToUint32(s string) (b uint32, ok bool) {
	var p [4]byte
	for i := 0; i < 4; i++ {
		if len(s) == 0 {

			return 0, false
		}
		if i > 0 {
			if s[0] != '.' {
				return 0, false
			}
			s = s[1:]
		}
		n, c, ok := dtoi(s)
		if !ok || n > 0xFF {
			return 0, false
		}
		s = s[c:]
		p[i] = byte(n)
	}
	if len(s) != 0 {
		return 0, false
	}
	return binary.BigEndian.Uint32(p[:]), true
}

func ParseIpToBuf(s string, buf []byte) (ip net.IP) {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return parseIpv4ToBuf(s, buf)
		case ':':
			return parseIpv6ToBuf(s, buf)
		}
	}
	return nil
}

func parseIpv4ToBuf(s string, buf []byte) (ip net.IP) {
	if len(buf) >= 4 {
		ip = net.IP(buf[:4])
	} else {
		ip = make(net.IP, 4)
	}
	for i := 0; i < 4; i++ {
		if len(s) == 0 {

			return nil
		}
		if i > 0 {
			if s[0] != '.' {
				return nil
			}
			s = s[1:]
		}
		n, c, ok := dtoi(s)
		if !ok || n > 0xFF {
			return nil
		}
		s = s[c:]
		ip[i] = byte(n)
	}
	if len(s) != 0 {
		return nil
	}

	return ip
}

func dtoiv4(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n > 0xFF {
			return 0xFF, i, false
		}
	}
	if i == 0 {
		return 0, 0, false
	}
	return n, i, true
}

func parseIpv6ToBuf(s string, buf []byte) (ip net.IP) {
	if len(buf) >= net.IPv6len {
		ip = net.IP(buf[:net.IPv6len])
	} else {
		ip = make(net.IP, net.IPv6len)
	}
	ellipsis := -1

	if len(s) >= 2 && s[0] == ':' && s[1] == ':' {
		ellipsis = 0
		s = s[2:]

		if len(s) == 0 {
			return ip
		}
	}

	i := 0
	for i < net.IPv6len {

		n, c, ok := xtoi(s)
		if !ok || n > 0xFFFF {
			return nil
		}

		if c < len(s) && s[c] == '.' {
			if ellipsis < 0 && i != net.IPv6len-net.IPv4len {

				return nil
			}
			if i+net.IPv4len > net.IPv6len {

				return nil
			}
			ip4 := parseIpv4ToBuf(s, buf)
			if ip4 == nil {
				return nil
			}
			ip[i] = ip4[0]
			ip[i+1] = ip4[1]
			ip[i+2] = ip4[2]
			ip[i+3] = ip4[3]
			s = ""
			i += net.IPv4len
			break
		}

		ip[i] = byte(n >> 8)
		ip[i+1] = byte(n)
		i += 2

		s = s[c:]
		if len(s) == 0 {
			break
		}

		if s[0] != ':' || len(s) == 1 {
			return nil
		}
		s = s[1:]

		if s[0] == ':' {
			if ellipsis >= 0 {
				return nil
			}
			ellipsis = i
			s = s[1:]
			if len(s) == 0 {
				break
			}
		}
	}

	if len(s) != 0 {
		return nil
	}

	if i < net.IPv6len {
		if ellipsis < 0 {
			return nil
		}
		n := net.IPv6len - i
		for j := i - 1; j >= ellipsis; j-- {
			ip[j+n] = ip[j]
		}
		for j := ellipsis + n - 1; j >= ellipsis; j-- {
			ip[j] = 0
		}
	} else if ellipsis >= 0 {

		return nil
	}
	return ip
}

func xtoi(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s); i++ {
		if '0' <= s[i] && s[i] <= '9' {
			n *= 16
			n += int(s[i] - '0')
		} else if 'a' <= s[i] && s[i] <= 'f' {
			n *= 16
			n += int(s[i]-'a') + 10
		} else if 'A' <= s[i] && s[i] <= 'F' {
			n *= 16
			n += int(s[i]-'A') + 10
		} else {
			break
		}
		if n >= big {
			return 0, i, false
		}
	}
	if i == 0 {
		return 0, i, false
	}
	return n, i, true
}

const big = 0xFFFFFF

func dtoi(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= big {
			return big, i, false
		}
	}
	if i == 0 {
		return 0, 0, false
	}
	return n, i, true
}
