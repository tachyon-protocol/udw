package tyVpnProtocol

import (
	"encoding/binary"
	"errors"
)

func (packet *VpnPacket) Decode(buf []byte) error {
	if len(buf) < 17 {
		return errors.New("len(buf) < 17")
	}
	packet.Cmd = buf[0]
	packet.ClientIdSender = binary.BigEndian.Uint64(buf[1:9])
	packet.ClientIdReceiver = binary.BigEndian.Uint64(buf[9:17])
	if len(buf) > 17 {
		packet.Data = buf[17:]
	}
	return nil
}
