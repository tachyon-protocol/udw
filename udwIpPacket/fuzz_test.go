package udwIpPacket_test

import (
	"github.com/tachyon-protocol/udw/udwIpPacket/udwIpPacketFuzzLib"
	"testing"
)

func TestByFuzz(ot *testing.T) {
	for _, cas := range udwIpPacketFuzzLib.GetTestCase() {
		udwIpPacketFuzzLib.Fuzz(cas)
	}
}
