package udwVpnPacket

import (
	"encoding/hex"
	"fmt"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwHex"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestNewIpPacketFromBuf(t *testing.T) {
	sample := `4500003c000040003406b0db68767f84ac1501d101bbc302504dc095da538312a0127120635c0000020405b40402080a2cd057b41294137501030305`
	buf := udwHex.MustDecodeStringToByteArray(sample)
	fmt.Println(">>> sample")
	fmt.Println(hex.Dump(buf))
	ipPacket, errMsg := NewIpv4PacketFromBuf(buf)
	udwErr.PanicIfErrorMsg(errMsg)
	ipPacket.RecomputeCheckSum()
	fmt.Println(">>> RecomputeCheckSum")
	fmt.Println(hex.Dump(ipPacket.SerializeToBuf()))
	udwTest.Ok(sample == ipPacket.HexDump())

}
