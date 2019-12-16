package tyVpnServer

import (
	"crypto/tls"
	"fmt"
	"github.com/tachyon-protocol/udw/tyTlsPacketDebugger"
	"github.com/tachyon-protocol/udw/tyVpnProtocol"
	"github.com/tachyon-protocol/udw/udwBinary"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwLog"
	"github.com/tachyon-protocol/udw/udwTlsSelfSignCertV2"
	"net"
	"sync"
	"time"
)

type vpnClient struct {
	id          uint64
	vpnIpOffset int
	vpnIp       net.IP

	connLock      sync.Mutex
	connToClient  net.Conn
	connRelaySide net.Conn

	connLastRwTimeLock sync.Mutex
	connLastRwTime time.Time
}

func (s *Server) gcClientThread() {
	go func() {
		for {
			time.Sleep(time.Minute * 5)
			s.lock.Lock()
			now := time.Now()
			for _, client := range s.clientMap {
				client.connLastRwTimeLock.Lock()
				t := client.connLastRwTime
				client.connLastRwTimeLock.Unlock()
				if now.Sub(t) < time.Minute*15 {
					continue
				}
				udwLog.Log("[dzr1zb5e3wz] gc remove client", client.id)
				if s.clientMap != nil {
					delete(s.clientMap, client.id)
				}
				s.vpnIpList[client.vpnIpOffset] = nil
				client.connLock.Lock()
				if client.connToClient != nil {
					_ = client.connToClient.Close()
				}
				if client.connRelaySide != nil {
					_ = client.connRelaySide.Close()
				}
				client.connLock.Unlock()
			}
			s.lock.Unlock()
		}
	}()
}

func (vc *vpnClient) getConnToClient() net.Conn {
	vc.connLock.Lock()
	conn := vc.connToClient
	vc.connLock.Unlock()
	return conn
}

func (s *Server) getClient(clientId uint64) *vpnClient {
	s.lock.Lock()
	if s.clientMap == nil {
		s.clientMap = map[uint64]*vpnClient{}
	}
	client := s.clientMap[clientId]
	s.lock.Unlock()
	if client != nil {
		client.connLastRwTimeLock.Lock()
		client.connLastRwTime = time.Now()
		client.connLastRwTimeLock.Unlock()
	}
	return client
}

func (s *Server) newOrUpdateClientFromDirectConn(clientId uint64, connToClient net.Conn) {
	s.lock.Lock()
	if s.clientMap == nil {
		s.clientMap = map[uint64]*vpnClient{}
	}
	client := s.clientMap[clientId]
	if client != nil {
		client.connLock.Lock()
		client.connToClient = connToClient //reconnect
		client.connLock.Unlock()
		s.lock.Unlock()
		return
	}
	client = &vpnClient{
		id:           clientId,
		connToClient: connToClient,
	}
	s.clientMap[client.id] = client
	err := s.clientAllocateVpnIp_NoLock(client)
	s.lock.Unlock()
	if err != nil {
		panic("[ub4fm53v26] " + err.Error())
	}
	return
}

func (s *Server) getOrNewClientFromRelayConn(clientId uint64) *vpnClient {
	s.lock.Lock()
	if s.clientMap == nil {
		s.clientMap = map[uint64]*vpnClient{}
	}
	client := s.clientMap[clientId]
	if client != nil {
		s.lock.Unlock()
		return client
	}
	client = &vpnClient{
		id: clientId,
	}
	left, right := tyVpnProtocol.NewInternalConnectionDual(func() {
		s.lock.Lock()
		delete(s.clientMap, clientId)
		s.lock.Unlock()
	}, nil)
	right = tls.Server(right, &tls.Config{
		Certificates: []tls.Certificate{ //TODO optimize allocate
			*udwTlsSelfSignCertV2.GetTlsCertificate(),
		},
		NextProtos: []string{"http/1.1"},
		MinVersion: tls.VersionTLS12,
	})
	client.connToClient = right
	client.connRelaySide = left
	s.clientMap[client.id] = client
	err := s.clientAllocateVpnIp_NoLock(client)
	go s.clientTcpConnHandle(client.getConnToClient())
	s.lock.Unlock()
	if err != nil {
		panic("[ub4fm53v26] " + err.Error())
	}
	go func() {
		vpnPacket := &tyVpnProtocol.VpnPacket{
			Cmd:              tyVpnProtocol.CmdForward,
			ClientIdSender:   s.clientId,
			ClientIdReceiver: clientId,
		}
		buf := make([]byte, 16*1024)
		bufW := udwBytes.NewBufWriter(nil)
		for {
			n, err := client.connRelaySide.Read(buf)
			if err != nil {
				udwLog.Log("[cz2xvv1smx] close conn", err)
				_ = client.connRelaySide.Close()
				return
			}
			if tyVpnProtocol.Debug {
				fmt.Println("read from connRelaySide write to relayConn", vpnPacket.ClientIdSender, "->", vpnPacket.ClientIdReceiver)
				if tyVpnProtocol.Debug {
					tyTlsPacketDebugger.Dump("---", buf[:n])
				}
			}
			vpnPacket.Data = buf[:n]
			bufW.Reset()
			vpnPacket.Encode(bufW)
			err = udwBinary.WriteByteSliceWithUint32LenNoAllocV2(s.getRelayConn(), bufW.GetBytes()) //TODO lock
			if err != nil {
				udwLog.Log("[ar1nr4wf3s]", err)
				continue
			}
		}
	}()
	return client
}
