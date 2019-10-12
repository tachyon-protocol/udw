package udwDnsPacket

import (
	"testing"
)

func TestDnsPacketFuzz(t *testing.T) {
	buf1 := []byte("0000\x00\x00000000\x010")
	_, _ = PacketReadFromByteSlice(buf1)
}
