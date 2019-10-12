package udwVpnPacket

import (
	"encoding/binary"
	"errors"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwNet"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"net"
	"strconv"
	"strings"
)

type ConnAddrTcpUdp struct {
	SrcIp   net.IP
	SrcPort uint16
	DstIp   net.IP
	DstPort uint16
}

func (addr ConnAddrTcpUdp) RevertPeer() ConnAddrTcpUdp {
	return ConnAddrTcpUdp{
		DstIp:   addr.SrcIp,
		DstPort: addr.SrcPort,
		SrcIp:   addr.DstIp,
		SrcPort: addr.DstPort,
	}
}

func (addr ConnAddrTcpUdp) String() string {
	return addr.SrcIp.String() + ":" + strconv.Itoa(int(addr.SrcPort)) + "-" + addr.DstIp.String() + ":" + strconv.Itoa(int(addr.DstPort))
}

func (addr ConnAddrTcpUdp) GetIpv4Array() (outArray [12]byte) {
	srcIp := addr.SrcIp.To4()
	dstIp := addr.DstIp.To4()
	outArray[0] = srcIp[0]
	outArray[1] = srcIp[1]
	outArray[2] = srcIp[2]
	outArray[3] = srcIp[3]

	outArray[4] = byte(addr.SrcPort >> 8)
	outArray[5] = byte(addr.SrcPort)
	outArray[6] = dstIp[0]
	outArray[7] = dstIp[1]
	outArray[8] = dstIp[2]
	outArray[9] = dstIp[3]

	outArray[10] = byte(addr.DstPort >> 8)
	outArray[11] = byte(addr.DstPort)
	return outArray
}

func GetConnAddrTcpUdpFromNetConn(conn net.Conn) (obj ConnAddrTcpUdp, errMsg string) {
	return GetConnAddrTcpUdpFromLocalAndRemoteAddr(conn.LocalAddr(), conn.RemoteAddr())
}

func GetConnAddrTcpUdpFromLocalAndRemoteAddr(localAddr net.Addr, remoteAddr net.Addr) (obj ConnAddrTcpUdp, errMsg string) {
	ip, port, errMsg := udwNet.GetIpAndPortFromNetAddr(localAddr)
	if errMsg != "" {
		return obj, errMsg
	}
	obj.SrcIp = ip
	obj.SrcPort = port
	ip, port, errMsg = udwNet.GetIpAndPortFromNetAddr(remoteAddr)
	if errMsg != "" {
		return obj, errMsg
	}
	obj.DstIp = ip
	obj.DstPort = port
	return obj, ""
}

var ConnAddrTcpUdpStrNotFitErr = errors.New("ConnAddrTcpUdpStrNotFit")

func ParseConnAddrTcpUdpFromString(in string) (out ConnAddrTcpUdp, err error) {
	tmp := strings.Split(in, "-")
	if len(tmp) != 2 {
		err = ConnAddrTcpUdpStrNotFitErr
		return
	}
	src := strings.Split(tmp[0], ":")
	if len(src) != 2 {
		err = ConnAddrTcpUdpStrNotFitErr
		return
	}
	out.SrcIp = net.ParseIP(src[0])
	err = udwErr.PanicToError(func() {
		out.SrcPort = uint16(udwStrconv.MustParseUint64(src[1]))
	})
	if err != nil {
		return
	}
	dst := strings.Split(tmp[1], ":")
	if len(dst) != 2 {
		err = ConnAddrTcpUdpStrNotFitErr
		return
	}
	out.DstIp = net.ParseIP(dst[0])
	err = udwErr.PanicToError(func() {
		out.DstPort = uint16(udwStrconv.MustParseUint64(dst[1]))
	})
	return
}

func (addr ConnAddrTcpUdp) Clone() ConnAddrTcpUdp {
	out := addr
	out.SrcIp = udwNet.IpClone(addr.SrcIp)
	out.DstIp = udwNet.IpClone(addr.DstIp)
	return out
}

func (packet *IpPacket) SetConnAddrTcpUdp(c ConnAddrTcpUdp) {
	packet.SetConnAddrTcpUdp__NoRecomputeCheckSum(c)
	packet.RecomputeCheckSum()
}

func (packet *IpPacket) SetConnAddrTcpUdp__NoRecomputeCheckSum(c ConnAddrTcpUdp) {
	copy(packet.buf[12:16], c.SrcIp.To4())
	copy(packet.buf[16:20], c.DstIp.To4())
	if packet.IsTcp() || packet.IsUdp() {
		ipLen := packet.GetIpHeaderLen()
		binary.BigEndian.PutUint16(packet.buf[ipLen:ipLen+2], c.SrcPort)
		binary.BigEndian.PutUint16(packet.buf[ipLen+2:ipLen+4], c.DstPort)
	}
}

func (packet *IpPacket) SetSrcIp(ip net.IP) {
	copy(packet.buf[12:16], ip.To4())
	packet.RecomputeCheckSum()
}

func (packet *IpPacket) SetSrcIp__NoRecomputeCheckSum(ip net.IP) {
	copy(packet.buf[12:16], ip.To4())
}

func (ipPacket *IpPacket) GetSrcAddrString() string {
	return ipPacket.GetSrcIp().String() + ":" + strconv.Itoa(int(ipPacket.GetSrcPort()))
}

func (packet *IpPacket) SetDstIp(ip net.IP) {
	copy(packet.buf[16:20], ip.To4())
	packet.RecomputeCheckSum()
}
func (packet *IpPacket) SetDstIp__NoRecomputeCheckSum(ip net.IP) {
	copy(packet.buf[16:20], ip.To4())
}

func (ipPacket *IpPacket) GetDstAddrString() string {
	return ipPacket.GetDstIp().String() + ":" + strconv.Itoa(int(ipPacket.GetDstPort()))
}

func (ipPacket *IpPacket) GetSrcDstAddrPeerString() string {
	return ipPacket.GetSrcAddrString() + "-" + ipPacket.GetDstAddrString()
}
func (ipPacket *IpPacket) GetDstSrcAddrPeerString() string {
	return ipPacket.GetDstAddrString() + "-" + ipPacket.GetSrcAddrString()
}
