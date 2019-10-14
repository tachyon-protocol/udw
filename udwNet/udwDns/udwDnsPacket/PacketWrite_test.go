package udwDnsPacket

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestWriteDomainToDnsPacket(ot *testing.T) {
	buf := make([]byte, 1024)
	pos, errMsg := writeDomainToDnsPacket("abc.com", buf, 0)
	udwTest.Equal(errMsg, "")
	udwTest.Equal(pos, 9)

	udwTest.Equal(buf[:pos], []byte{
		0x03, 0x61, 0x62, 0x63,
		0x03, 0x63, 0x6f, 0x6d,
		0x00,
	})

	pos, errMsg = writeDomainToDnsPacket("", buf, 0)
	udwTest.Equal(errMsg, "")
	udwTest.Equal(pos, 1)

	udwTest.Equal(buf[:pos], []byte{
		0x00,
	})
}

func TestNewDnsPacketForRequest(t *testing.T) {
	packet := NewDnsPacketForRequest("abc.com", DnsTypeA)
	udwTest.Equal(packet.HasEdns4096, true)
	udwTest.Equal(packet.QuestionName, "abc.com")
	buf, errMsg := PacketWriteToByteSlice(packet, nil)
	udwTest.Equal(errMsg, "")
	packet2, errMsg := PacketReadFromByteSlice(buf)
	udwTest.Equal(errMsg, "")
	udwTest.Equal(packet2.HasEdns4096, true)
	udwTest.Equal(packet2.QuestionName, "abc.com")
}
