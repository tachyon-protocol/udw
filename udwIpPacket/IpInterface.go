package udwIpPacket

type TunEventOnTunProcessor struct {
	SetWriteToTunCallback func(f func(ipPacket IpPacket))

	ReadFromTun func(ipPacket IpPacket)

	Close func()
}
