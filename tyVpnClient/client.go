package tyVpnClient

import (
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/tachyon-protocol/udw/tyTls"
	"github.com/tachyon-protocol/udw/tyVpnProtocol"
	"github.com/tachyon-protocol/udw/tyVpnRouteServer/tyVpnRouteClient"
	"github.com/tachyon-protocol/udw/udwBinary"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwIo"
	"github.com/tachyon-protocol/udw/udwIpPacket"
	"github.com/tachyon-protocol/udw/udwLog"
	"github.com/tachyon-protocol/udw/udwNet"
	"github.com/tachyon-protocol/udw/udwNet/udwIPNet"
	"github.com/tachyon-protocol/udw/udwNet/udwTapTun"
	"github.com/tachyon-protocol/udw/udwRand"
	"io"
	"net"
	"sync"
	"time"
	"github.com/tachyon-protocol/udw/udwClose"
)

type Config struct {
	ServerIp   string `json:",omitempty"`
	ServerTKey string `json:",omitempty"`

	IsRelay            bool `json:",omitempty"`
	ExitServerClientId uint64 `json:",omitempty"` //required when IsRelay is true
	ExitServerTKey     string `json:",omitempty"` //required when IsRelay is true

	ServerChk                   string `json:",omitempty"` // if it is "", it will use InsecureSkipVerify
	DisableUsePublicRouteServer bool `json:",omitempty"`
}

type Client struct {
	req                  Config
	clientId             uint64
	clientIdToExitServer uint64
	keepAliveChan        chan uint64
	connLock             sync.Mutex
	directVpnConn        net.Conn
	vpnConn              net.Conn
	tlsConfig            *tls.Config
	tun                  io.ReadWriteCloser
	thisCsCmdId          int
	closer               udwClose.Closer
	rc                   int
	rcLocker             sync.Mutex
	rcWg                 sync.WaitGroup
}

func (c *Client) initObj() (errMsg string){
	tyTls.EnableTlsVersion13()
	c.clientId = tyVpnProtocol.GetClientId(0)
	c.clientIdToExitServer = c.clientId
	if c.req.IsRelay {
		c.clientIdToExitServer = tyVpnProtocol.GetClientId(1)
		if c.req.ExitServerClientId == 0 {
			return "ExitServerClientId can be empty when use relay mode"
		}
	}
	if c.req.ServerChk == "" {
		c.tlsConfig = newInsecureClientTlsConfig()
	} else {
		var errMsg string
		c.tlsConfig, errMsg = tyTls.NewClientTlsConfigWithChk(tyTls.NewClientTlsConfigWithChkReq{
			ServerChk: c.req.ServerChk,
		})
		if errMsg!=""{
			return errMsg
		}
	}
	return ""
}

func (c *Client) connectL1(req Config) {
	//defer udwLog.Log("close connectL1")
	defer c.rcDec()
	c.req = req
	//udwLog.Log("connectL1 1",c.thisCsCmdId)
	setLastError("")
	errMsg := c.initObj()
	if errMsg!=""{
		c.errorDurationConnecting(errMsg)
		return
	}
	if c.closer.IsClose(){
		return
	}
	//udwLog.Log("connectL1 2",c.thisCsCmdId)
	c.tryUseRouteServer()
	//udwLog.Log("connectL1 3",c.thisCsCmdId)
	if c.closer.IsClose(){
		return
	}
	tun, err := createTun(c.req.ServerIp)
	if err!=nil{
		c.errorDurationConnecting(err.Error())
		return
	}
	c.tun = tun
	//udwLog.Log("connectL1 4",c.thisCsCmdId)
	c.closer.AddOnClose(func(){
		c.tun.Close()
	})
	if c.closer.IsClose(){
		return
	}
	//udwLog.Log("connectL1 5",c.thisCsCmdId)
	c.reconnect()
	if c.closer.IsClose(){
		return
	}
	c.initKeepAliveThread()
	c.initTunReadThread()
	c.initConnReadThread()
	c.setInnerCsIfCsCmdIdValid(innerCsConnected)
}

func createTun(vpnServerIp string) (tun io.ReadWriteCloser, err error) {
	vpnClientIp := net.ParseIP("172.21.0.1")
	includeIpNetSet := udwIPNet.NewAllPassIpv4Net()
	includeIpNetSet.RemoveIpString(vpnServerIp)
	tunCreateCtx := &udwTapTun.CreateIpv4TunContext{
		SrcIp:        vpnClientIp,
		DstIp:        vpnClientIp,
		FirstIp:      vpnClientIp,
		DhcpServerIp: vpnClientIp,
		Mtu:          tyVpnProtocol.Mtu,
		Mask:         net.CIDRMask(30, 32),
	}
	err = udwTapTun.CreateIpv4Tun(tunCreateCtx)
	if err != nil {
		return nil, errors.New("[3xa38g7vtd] " + err.Error())
	}
	tunNamed := tunCreateCtx.ReturnTun
	vpnGatewayIp := vpnClientIp
	err = udwErr.PanicToError(func() {
		configLocalNetwork()
		ctx := udwNet.NewRouteContext()
		for _, ipNet := range includeIpNetSet.GetIpv4NetList() {
			goIpNet := ipNet.ToGoIPNet()
			ctx.MustRouteSet(*goIpNet, vpnGatewayIp)
		}
	})
	if err != nil {
		_ = tunNamed.Close()
		return nil, errors.New("[r8y8d5ash4] " + err.Error())
	}
	var closeOnce sync.Once
	return udwIo.StructWriterReaderCloser{
		Reader: tunNamed,
		Writer: tunNamed,
		Closer: udwIo.CloserFunc(func() error {
			closeOnce.Do(func() {
				//trySendPacketToTun()
				_ = tunNamed.Close()
				err := udwErr.PanicToError(func() {
					recoverLocalNetwork()
				})
				if err != nil {
					udwLog.Log("error", "uninstallAllPassRoute", err.Error())
				}
			})
			return nil
		}),
	}, nil
}

//func trySendPacketToTun(){
//	conn,err:=net.Dial("udp","172.21.0.1:10000")
//	if err!=nil{
//		return
//	}
//	defer conn.Close()
//	conn.Write([]byte{0})
//	conn.Close()
//}

func newInsecureClientTlsConfig() *tls.Config {
	return &tls.Config{
		ServerName:         udwRand.MustCryptoRandToReadableAlpha(5) + ".com",
		InsecureSkipVerify: true,
		NextProtos:         []string{"http/1.1", "h2"},
		MinVersion:         tls.VersionTLS12,
	}
}

func (c *Client) tryUseRouteServer() {
	if c.req.ServerIp!=""{
		return
	}
	if c.req.DisableUsePublicRouteServer {
		setLastError("need config ServerIp")
		c.closer.Close()
		return
	}
	//udwLog.Log("connectL1 6",c.thisCsCmdId)
	routeC := tyVpnRouteClient.Rpc_NewClient(tyVpnProtocol.PublicRouteServerAddr)
	//udwLog.Log("connectL1 7",c.thisCsCmdId)
	list, rpcErr := routeC.VpnNodeList()
	//udwLog.Log("connectL1 8",c.thisCsCmdId)
	if rpcErr != nil {
		setLastError(rpcErr.Error())
		c.closer.Close()
		return
	}
	locker := sync.Mutex{}
	var fastNode tyVpnRouteClient.VpnNode
	wg := sync.WaitGroup{}
	for _, node := range list {
		node := node
		wg.Add(1)
		go func() {
			err := Ping(PingReq{
				Ip:        node.Ip,
				ServerChk: node.ServerChk,
			})
			if err == nil {
				locker.Lock()
				if fastNode.Ip == "" {
					fastNode = node
				}
				locker.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	//udwLog.Log("connectL1 9",c.thisCsCmdId)
	if fastNode.Ip == "" {
		setLastError("all ping lost")
		c.closer.Close()
		return
	}
	c.req.ServerIp = fastNode.Ip
	c.req.ServerChk = fastNode.ServerChk
	fmt.Println("ping to get ip [" + c.req.ServerIp + "]")
}

func (c *Client) initTunReadThread(){
	c.rcInc()
	go func(){
		//defer udwLog.Log("close initTunReadThread")
		defer c.rcDec()
		vpnPacket := &tyVpnProtocol.VpnPacket{
			Cmd:              tyVpnProtocol.CmdData,
			ClientIdSender:   c.clientIdToExitServer,
			ClientIdReceiver: c.req.ExitServerClientId,
		}
		buf := make([]byte, 16*1024)
		bufW := udwBytes.NewBufWriter(nil)
		c.connLock.Lock()
		vpnConn := c.vpnConn
		c.connLock.Unlock()
		for {
			//udwLog.Log("before c.tun.Read(buf)")
			n, err := c.tun.Read(buf)
			//udwLog.Log("after c.tun.Read(buf)")
			if c.closer.IsClose(){
				return
			}
			if err != nil {
				fmt.Println("[upe1hcb1q39h] " + err.Error())
				continue
			}
			vpnPacket.Data = buf[:n]
			bufW.Reset()
			vpnPacket.Encode(bufW)
			for {
				err = udwBinary.WriteByteSliceWithUint32LenNoAllocV2(vpnConn, bufW.GetBytes())
				if c.closer.IsClose(){
					return
				}
				if err != nil {
					c.connLock.Lock()
					_vpnConn := c.vpnConn
					c.connLock.Unlock()
					if vpnConn == _vpnConn {
						time.Sleep(time.Millisecond * 50)
					} else {
						vpnConn = _vpnConn
						udwLog.Log("[mpy2nwx1qck] tun read use new vpn conn")
					}
					continue
				}
				break
			}
		}
	}()
}

func (c *Client) initConnReadThread() {
	c.rcInc()
	go func(){
		//defer udwLog.Log("close initConnReadThread")
		defer c.rcDec()
		vpnPacket := &tyVpnProtocol.VpnPacket{}
		buf := udwBytes.NewBufWriter(nil)
		c.connLock.Lock()
		vpnConn := c.vpnConn
		c.connLock.Unlock()
		for {
			buf.Reset()
			for {
				err := udwBinary.ReadByteSliceWithUint32LenToBufW(vpnConn, buf)
				if c.closer.IsClose(){
					return
				}
				if err != nil {
					c.connLock.Lock()
					_vpnConn := c.vpnConn
					c.connLock.Unlock()
					if vpnConn == _vpnConn {
						time.Sleep(time.Millisecond * 50)
					} else {
						vpnConn = _vpnConn
						udwLog.Log("[zdb1mbq1v1kxh] vpn conn read use new vpn conn")
					}
					continue
				}
				break
			}
			err := vpnPacket.Decode(buf.GetBytes())
			if err!=nil{
				udwLog.Log(err.Error())
				continue
			}
			switch vpnPacket.Cmd {
			case tyVpnProtocol.CmdData:
				ipPacket, errMsg := udwIpPacket.NewIpv4PacketFromBuf(vpnPacket.Data)
				if errMsg != "" {
					udwLog.Log("[zdy1mx9y3h]", errMsg)
					continue
				}
				_, err = c.tun.Write(ipPacket.SerializeToBuf())
				if err != nil {
					udwLog.Log("[wmw12fyr9e] TUN Write error", err)
				}
			case tyVpnProtocol.CmdKeepAlive:
				i := binary.LittleEndian.Uint64(vpnPacket.Data)
				select {
					case c.keepAliveChan <- i:
					default:
						continue
				}
			default:
				udwLog.Log("[h67hrf4kda] unexpect cmd", vpnPacket.Cmd)
			}
		}
	}()
}