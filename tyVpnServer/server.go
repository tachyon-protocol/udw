package tyVpnServer

import (
	"crypto/tls"
	"fmt"
	"github.com/tachyon-protocol/udw/tyTls"
	"github.com/tachyon-protocol/udw/tyVpnProtocol"
	"github.com/tachyon-protocol/udw/tyVpnRouteServer/tyVpnRouteClient"
	"github.com/tachyon-protocol/udw/udwBinary"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwClose"
	"github.com/tachyon-protocol/udw/udwConsole"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwIpPacket"
	"github.com/tachyon-protocol/udw/udwIpToCountryV2"
	"github.com/tachyon-protocol/udw/udwLog"
	"github.com/tachyon-protocol/udw/udwNet"
	"github.com/tachyon-protocol/udw/udwNet/udwTapTun"
	"github.com/tachyon-protocol/udw/udwStrings"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/tachyon-protocol/udw/udwJson"
	"github.com/tachyon-protocol/udw/tyVpnClient"
)

func RunCmd(req ServerRunReq){
	server := &Server{}
	server.Run(req)
}

type ServerRunReq struct {
	UseRelay        bool
	RelayServerIp   string
	RelayServerTKey string
	RelayServerChk  string
	SelfTKey                   string
	NoNeedSelfTKey bool
	BlockCountryCodeListS      string // empty string do not block any country code, look like "KP,IR,RU"
	SelfIp string
	RegisterToRouteServer bool
}

type Server struct {
	clientId               uint64
	tun                    *udwTapTun.TunTapObj
	relayConnKeepAliveChan chan uint64

	lock           sync.Mutex
	clientMap      map[uint64]*vpnClient
	vpnIpList      [maxCountVpnIp]*vpnClient
	nextVpnIpIndex int
	relayConn      net.Conn

	req                  ServerRunReq
	blockCountryCodeList []string
}

func (s *Server) Run(config ServerRunReq) {
	tyTls.EnableTlsVersion13()
	s.req = config
	s.clientId = tyVpnProtocol.GetClientId(0)
	fmt.Println("ClientId:", s.clientId,"ProtocolVersion: ",tyVpnProtocol.Version)

	tun, err := udwTapTun.NewTun("")
	udwErr.PanicIfError(err)
	err = udwTapTun.SetP2PIpAndUp(udwTapTun.SetP2PIpRequest{
		IfaceName: tun.Name(),
		SrcIp:     udwNet.Ipv4AddAndCopyWithBuffer(READONLY_vpnIpStart, 2, nil),
		DstIp:     udwNet.Ipv4AddAndCopyWithBuffer(READONLY_vpnIpStart, 1, nil),
		Mtu:       tyVpnProtocol.Mtu,
		Mask:      net.CIDRMask(16, 32),
	})
	udwErr.PanicIfError(err)
	s.tun = tun
	networkConfig()
	serverStorageInfo:=getServerStorageInfo()
	if s.req.SelfTKey=="" && s.req.NoNeedSelfTKey==false{
		s.req.SelfTKey = serverStorageInfo.SelfTKey
	}
	routeC := tyVpnRouteClient.Rpc_NewClient(tyVpnProtocol.PublicRouteServerAddr)
	if s.req.SelfIp==""{
		ip,rpcErr := routeC.GetIp()
		if rpcErr!=nil{
			panic(rpcErr.Error())
		}
		s.req.SelfIp = ip
	}
	sTlsConfig, errMsg := tyTls.NewServerTlsConfigWithChk(tyTls.NewServerTlsConfigWithChkReq{
		ServerCert: *serverStorageInfo.ServerTlsCert,
	})
	udwErr.PanicIfErrorMsg(errMsg)
	//fmt.Println("ServerChk: " + serverStorageInfo.ServerChk)
	clientConfig:=tyVpnClient.Config{
		ServerIp: s.req.SelfIp,
		ServerChk: serverStorageInfo.ServerChk,
		ServerTKey: s.req.SelfTKey,
	}
	clientConnectConfigS:=udwJson.MustMarshalToString(clientConfig)
	fmt.Println("ServerConnectConfig: "+tyVpnClient.MarshalConfig(clientConfig))
	fmt.Println("ServerConnectConfigJson: "+clientConnectConfigS)
	fmt.Println("Server started ✔")
	if s.req.BlockCountryCodeListS != "" {
		s.blockCountryCodeList = strings.Split(s.req.BlockCountryCodeListS, ",")
	}
	if s.req.RegisterToRouteServer {
		go func() {
			for {
				err1, err2 := routeC.VpnNodeRegister(tyVpnRouteClient.VpnNode{
					ServerChk: serverStorageInfo.ServerChk,
				})
				if err1 != "" {
					udwLog.Log("4etc1ghe1khj " + err1)
				}
				if err2 != nil {
					udwLog.Log("yew68bub3a " + err2.Error())
				}
				time.Sleep(time.Second * 30)
			}
		}()
	}

	//read thread from TUN
	go func() {
		bufR := make([]byte, 16*1024)
		bufW := udwBytes.NewBufWriter(nil)
		vpnPacket := &tyVpnProtocol.VpnPacket{
			ClientIdSender: s.clientId,
			Cmd:            tyVpnProtocol.CmdData,
		}
		for {
			n, err := tun.Read(bufR)
			if err != nil {
				udwLog.Log("[m7j1pw1vr7] TUN Read failed", err)
				continue
			}
			packetBuf := bufR[:n]
			ipPacket, errMsg := udwIpPacket.NewIpv4PacketFromBuf(packetBuf)
			if errMsg != "" {
				udwLog.Log("[wj1nz633mg] TUN Read parse IPv4 failed", errMsg)
				continue
			}
			ip := ipPacket.GetDstIp()
			if udwNet.IsPrivateNetwork(ipPacket.GetSrcIp()) {
				//udwLog.Log("[ye723euu1ah] private IP address is not allowed", ipPacket.GetSrcAddrString())
				continue
			}
			client := s.getClientByVpnIp(ip)
			if client == nil {
				udwLog.Log("[r1tp9rk84m] TUN Read no such client", ipPacket.GetSrcAddrString())
				continue
			}
			//vpnPacket := &tachyonVpnProtocol.VpnPacket{
			//	ClientIdSender:   s.clientId,
			//	ClientIdReceiver: client.id,
			//	Cmd:              tachyonVpnProtocol.CmdData,
			//}
			vpnPacket.ClientIdReceiver = client.id
			if tyVpnProtocol.Debug {
				fmt.Println("read from tun, write to client", vpnPacket.ClientIdSender, "->", vpnPacket.ClientIdReceiver)
			}
			ipPacket.SetDstIp__NoRecomputeCheckSum(READONLY_vpnIpClient)
			ipPacket.TcpFixMss__NoRecomputeCheckSum(tyVpnProtocol.Mss)
			ipPacket.RecomputeCheckSum()
			vpnPacket.Data = ipPacket.SerializeToBuf()
			bufW.Reset()
			vpnPacket.Encode(bufW)
			_ = udwBinary.WriteByteSliceWithUint32LenNoAllocV2(client.getConnToClient(), bufW.GetBytes()) //TODO
		}
	}()
	closer := udwClose.NewCloser()
	//two methods to accept new vpn conn
	if s.req.UseRelay {
		err := s.connectToRelay()
		udwErr.PanicIfError(err)
		s.relayConnKeepAliveThread()
	} else {
		closerFn := udwNet.TcpNewListener(":"+strconv.Itoa(tyVpnProtocol.VpnPort), func(conn net.Conn) {
			if s.clientConnFilter(conn) == false {
				_ = conn.Close()
				return
			}
			conn = tls.Server(conn, sTlsConfig)
			s.clientTcpConnHandle(conn)
		})
		closer.AddOnClose(closerFn)
	}
	udwConsole.WaitForExit()
	closer.Close()
}

// return true as pass
func (s *Server) clientConnFilter(connToClient net.Conn) bool {
	if len(s.blockCountryCodeList) > 0 {
		ip, _, errMsg := udwNet.GetIpAndPortFromNetAddr(connToClient.RemoteAddr())
		if errMsg == "" {
			cc := udwIpToCountryV2.MustGetCountryIsoCode(ip)
			if cc != "" && udwStrings.IsInSlice(s.blockCountryCodeList, cc) {
				return false
			}
		}
	}
	return true
}

func (s *Server) clientTcpConnHandle(connToClient net.Conn) {
	vpnPacket := &tyVpnProtocol.VpnPacket{}
	bufW := udwBytes.NewBufWriter(nil)
	for {
		bufW.Reset()
		err := udwBinary.ReadByteSliceWithUint32LenToBufW(connToClient, bufW)
		if err != nil {
			udwLog.Log("[tw1me5hux3] close conn", err, connToClient.RemoteAddr())
			_ = connToClient.Close()
			return
		}
		err = vpnPacket.Decode(bufW.GetBytes())
		if err != nil {
			udwLog.Log("[m1ds6vv2n8] close conn", err)
			_ = connToClient.Close()
			return
		}
		writePacketFn:=func(cmd byte,data []byte) (err error){
			bufW.Reset()
			vpnPacket.Cmd = cmd
			vpnPacket.ClientIdReceiver = vpnPacket.ClientIdSender
			vpnPacket.ClientIdSender = s.clientId
			vpnPacket.Data = data
			vpnPacket.Encode(bufW)
			err = udwBinary.WriteByteSliceWithUint32LenNoAllocV2(connToClient, bufW.GetBytes())
			return err
		}
		switch vpnPacket.Cmd {
		case tyVpnProtocol.CmdPing, tyVpnProtocol.CmdKeepAlive:
			err:=writePacketFn(vpnPacket.Cmd,vpnPacket.Data)
			if err!=nil{
				udwLog.Log("[2cpj1sbv37s] close conn", err)
				_ = connToClient.Close()
				return
			}
		case tyVpnProtocol.CmdHandshake:
			isOk:=false
			if s.req.SelfTKey == "" {
				udwLog.Log("[4z734vc9pn] New client sent handshake ✔ server not require TKey", connToClient.RemoteAddr())
				isOk = true
			}else if len(s.req.SelfTKey) == len(string(vpnPacket.Data)) && s.req.SelfTKey == string(vpnPacket.Data) {
				udwLog.Log("[agz7rzq1kr9] New client TKey matched ✔", connToClient.RemoteAddr())
				isOk = true
			}else {
				udwLog.Log("[wzh56ty1bur] New client TKey not match ✘ close conn", connToClient.RemoteAddr())
				isOk = false
			}
			if isOk{
				s.newOrUpdateClientFromDirectConn(vpnPacket.ClientIdSender, connToClient)
				err:=writePacketFn(tyVpnProtocol.CmdErr,nil)
				if err!=nil{
					udwLog.Log("[n4wbkq9q68] close conn", err)
					_ = connToClient.Close()
					return
				}
			} else {
				err:=writePacketFn(tyVpnProtocol.CmdErr,[]byte("tegryjenqb TKey not match"))
				if err!=nil{
					udwLog.Log("[fsegxt36f4] close conn", err)
					return
				}
				_ = connToClient.Close()
			}
		case tyVpnProtocol.CmdData:
			client := s.getClient(vpnPacket.ClientIdSender)
			if client == nil {
				_ = connToClient.Close()
				udwLog.Log("[k692xqw1d2n] CmdData close conn cause no such client", vpnPacket.ClientIdSender, connToClient.RemoteAddr())
				return
			}
			//client := s.getOrNewClientFromDirectConn(vpnPacket.ClientIdSender, connToClient)
			ipPacket, errMsg := udwIpPacket.NewIpv4PacketFromBuf(vpnPacket.Data)
			if errMsg != "" {
				udwLog.Log("[txd5xn4ex7] close conn", errMsg, "ipPacket.IsIpv4:", ipPacket.IsIpv4(), "ipPacket.Ipv4HasMoreFragments:", ipPacket.Ipv4HasMoreFragments())
				_ = connToClient.Close()
				return
			}
			vpnIp := udwNet.Ipv4AddAndCopyWithBuffer(READONLY_vpnIpStart, uint32(client.vpnIpOffset), bufW)
			ipPacket.SetSrcIp__NoRecomputeCheckSum(vpnIp)
			ipPacket.TcpFixMss__NoRecomputeCheckSum(tyVpnProtocol.Mss)
			ipPacket.RecomputeCheckSum()
			_, err = s.tun.Write(ipPacket.SerializeToBuf())
			if err != nil {
				udwLog.Log("[x8z73fm1x5] close conn", err)
				_ = connToClient.Close()
				return
			}
		case tyVpnProtocol.CmdForward:
			client := s.getClient(vpnPacket.ClientIdSender)
			if client == nil {
				_ = connToClient.Close()
				udwLog.Log("[be8meu1vhm1d] CmdForward close conn cause no such client", vpnPacket.ClientIdSender, connToClient.RemoteAddr())
				return
			}
			nextPeer := s.getClient(vpnPacket.ClientIdReceiver)
			if nextPeer == nil {
				udwLog.Log("[4tz1d2932g] forward failed nextPeer[", vpnPacket.ClientIdReceiver, "] == nil")
				continue
			}
			err := udwBinary.WriteByteSliceWithUint32LenNoAllocV2(nextPeer.getConnToClient(), bufW.GetBytes()) //TLS layer
			if err != nil {
				udwLog.Log("[va1gz58zm3] forward failed", err)
				continue
			}
		default:
			_ = connToClient.Close()
			udwLog.Log("[rjb3nay1ezg] Cmd unknown", vpnPacket.Cmd, "close conn", connToClient.RemoteAddr())
			return
		}
	}
}

