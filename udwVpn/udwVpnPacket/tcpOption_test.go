package udwVpnPacket

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwHex"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestTcpOptionV2_GetTcpOption(t *testing.T) {
	ipBuf := udwHex.MustDecodeStringToByteArray(`45000040ef2d40004006c341ac1500013b6fa0c3c10201bb6f551b1100000000b002ffff1826000002040543010303050101080a151b34c00000000004020000`)
	packet := &IpPacket{buf: ipBuf}
	udwTest.Ok(packet.TcpGetMss() == 1347)
}

func TestCheckTcpOptionCorrect(t *testing.T) {
	tcpOptionBuf := []byte{2, 0, 0}
	errMsg := checkTcpOptionCorrect(tcpOptionBuf)
	udwTest.Equal(errMsg, "tcp optionSize<2 0")
}

func TestTcpFixMss(t *testing.T) {

	ipBuf := udwHex.MustDecodeStringToByteArray(`45000040ef2d40004006c341ac1500013b6fa0c3c10201bb6f551b1100000000b002ffff1826000002040543010303050101080a151b34c00000000004020000`)
	packet := &IpPacket{buf: ipBuf}
	udwTest.Equal(packet.TcpGetMss(), uint16(1347))
	packet.TcpFixMss(1300)
	udwTest.Equal(packet.TcpGetMss(), uint16(1300))
	packet.TcpFixMss(1400)
	udwTest.Equal(packet.TcpGetMss(), uint16(1300))

	ipBuf = udwHex.MustDecodeStringToByteArray(`45100034d34240004006f19dac15000123b9a604df040016387b66fe7ccb906180100fd60c4500000101080a6d92ffa20022ebb5`)
	packet = &IpPacket{buf: ipBuf}
	udwTest.Equal(packet.TcpGetMss(), uint16(0))
	packet.TcpFixMss(1300)
	udwTest.Equal(packet.TcpGetMss(), uint16(0))

	ipBuf = udwHex.MustDecodeStringToByteArray(`4500003abf070000ff116bb0ac15000172727272e2ee00350026780dea4a01000001000000000000057371696d6702717103636f6d0000010001`)
	packet = &IpPacket{buf: ipBuf}
	udwTest.Equal(packet.TcpGetMss(), uint16(0))
	packet.TcpFixMss(1300)
	udwTest.Equal(packet.TcpGetMss(), uint16(0))

	ipBuf = udwHex.MustDecodeStringToByteArray(`45000040ef2d40004006c341ac1500013b6fa0c3c10201bb6f551b1100000000b002ffff1826000002040543010303050101080a151b34c00000000004020000`)
	const benchNum = 1e4
	packetList := make([]IpPacket, benchNum)
	for i := 0; i < benchNum; i++ {
		packetList[i].buf = udwBytes.Clone(ipBuf)
	}

	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(benchNum)
		for i := 0; i < benchNum; i++ {
			packetList[i].TcpFixMss(1300)
		}
	})
}
