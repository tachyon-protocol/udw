package udwVpnPacket_test

import (
	"github.com/tachyon-protocol/udw/udwVpn/udwVpnPacket/udwVpnPacketFuzzLib"
	"testing"
)

func TestByFuzz(ot *testing.T) {
	for _, cas := range udwVpnPacketFuzzLib.GetTestCase() {
		udwVpnPacketFuzzLib.Fuzz(cas)
	}
}
