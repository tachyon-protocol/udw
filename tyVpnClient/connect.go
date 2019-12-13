package tyVpnClient

import (
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/tachyon-protocol/udw/tyVpnProtocol"
	"github.com/tachyon-protocol/udw/udwBinary"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwLog"
	"net"
	"strconv"
	"time"
)

func (c *Client) connect() error {
	vpnConn, err := net.Dial("tcp", c.req.ServerIp+":"+strconv.Itoa(tyVpnProtocol.VpnPort))
	if err != nil {
		return errors.New("[w7syh9d1zgd] " + err.Error())
	}
	vpnConn = tls.Client(vpnConn, c.tlsConfig)
	var (
		handshakeVpnPacket = tyVpnProtocol.VpnPacket{
			Cmd:            tyVpnProtocol.CmdHandshake,
			ClientIdSender: c.clientId,
			Data:           []byte(c.req.ServerTKey),
		}
		handshakeBuf = udwBytes.NewBufWriter(nil)
	)
	handshakeVpnPacket.Encode(handshakeBuf)
	err = udwBinary.WriteByteSliceWithUint32LenNoAllocV2(vpnConn, handshakeBuf.GetBytes())
	if err != nil {
		_ = vpnConn.Close()
		return errors.New("[52y73b9e89] " + err.Error())
	}
	c.connLock.Lock()
	c.directVpnConn = vpnConn
	c.connLock.Unlock()
	c.closer.AddOnClose(func(){
		c.connLock.Lock()
		directVpnConn:=c.directVpnConn
		c.connLock.Unlock()
		directVpnConn.Close()
	})
	serverType := "DIRECT"
	if c.req.IsRelay {
		serverType = "RELAY"
		var (
			connRelaySide, plain = tyVpnProtocol.NewInternalConnectionDual(nil, nil)
			relayConn            = vpnConn
		)
		c.closer.AddOnClose(func() {
			connRelaySide.Close()
			plain.Close()
		})
		vpnConn = tls.Client(plain, c.tlsConfig)
		c.rcInc()
		//read from relay conn, write to vpn conn
		go func() {
			defer c.rcDec()
			var (
				buf       = udwBytes.NewBufWriter(nil)
				vpnPacket = &tyVpnProtocol.VpnPacket{}
			)
			for {
				buf.Reset()
				err := udwBinary.ReadByteSliceWithUint32LenToBufW(relayConn, buf)
				if c.closer.IsClose(){
					return
				}
				if err != nil {
					udwLog.Log("[wua1j5ps1pam] close 3 connections", err)
					_ = connRelaySide.Close()
					_ = plain.Close()
					_ = vpnConn.Close()
					return
				}
				err = vpnPacket.Decode(buf.GetBytes())
				if err != nil {
					udwLog.Log("[kj4v98z1fzc] close 3 connections", err)
					_ = connRelaySide.Close()
					_ = plain.Close()
					_ = vpnConn.Close()
					return
				}
				switch vpnPacket.Cmd {
				case tyVpnProtocol.CmdForward:
					_, err := connRelaySide.Write(vpnPacket.Data)
					if err != nil {
						udwLog.Log("[8gys171bvm] close 3 connections", err)
						_ = connRelaySide.Close()
						_ = plain.Close()
						_ = vpnConn.Close()
						return
					}
				case tyVpnProtocol.CmdKeepAlive:
					c.keepAliveChan <- binary.LittleEndian.Uint64(vpnPacket.Data)
				default:
					fmt.Println("[a3t7vfh1ms] Unexpected Cmd[", vpnPacket.Cmd, "]")
				}
			}
		}()
		c.rcInc()
		//read from vpn conn, write to relay conn
		go func() {
			c.rcDec()
			vpnPacket := &tyVpnProtocol.VpnPacket{
				Cmd:              tyVpnProtocol.CmdForward,
				ClientIdSender:   c.clientId,
				ClientIdReceiver: c.req.ExitServerClientId,
			}
			buf := make([]byte, 16*1024)
			bufW := udwBytes.NewBufWriter(nil)
			for {
				n, err := connRelaySide.Read(buf)
				if c.closer.IsClose(){
					return
				}
				if err != nil {
					udwLog.Log("[e9erq1bwd1] close 3 connections", err)
					_ = connRelaySide.Close()
					_ = plain.Close()
					_ = vpnConn.Close()
					return
				}
				vpnPacket.Data = buf[:n]
				bufW.Reset()
				vpnPacket.Encode(bufW)
				err = udwBinary.WriteByteSliceWithUint32LenNoAllocV2(relayConn, bufW.GetBytes())
				if err != nil {
					udwLog.Log("[n2cvu3w1cb] close 3 connections", err)
					_ = connRelaySide.Close()
					_ = plain.Close()
					_ = vpnConn.Close()
					return
				}
			}
		}()
		udwLog.Log("send handshake to ExitServer...")
		handshakeVpnPacket.ClientIdSender = c.clientIdToExitServer
		handshakeVpnPacket.Data = []byte(c.req.ExitServerTKey)
		handshakeBuf.Reset()
		handshakeVpnPacket.Encode(handshakeBuf)
		err = udwBinary.WriteByteSliceWithUint32LenNoAllocV2(vpnConn, handshakeBuf.GetBytes())
		if err != nil {
			_ = vpnConn.Close()
			return errors.New("[q3nwv1ebx1cd] " + err.Error())
		}
		udwLog.Log("sent handshake to ExitServer ✔")
	}
	fmt.Println("Connected to", serverType, "Server ✔")
	c.connLock.Lock()
	c.vpnConn = vpnConn
	c.connLock.Unlock()
	c.closer.AddOnClose(func(){
		c.connLock.Lock()
		vpnConn:=c.vpnConn
		c.connLock.Unlock()
		vpnConn.Close()
	})
	return nil
}

func (c *Client) initKeepAliveThread() {
	c.keepAliveChan = make(chan uint64, 10)
	c.rcInc()
	go func() {
		//defer udwLog.Log("close initKeepAliveThread")
		defer c.rcDec()
		i := uint64(0)
		vpnPacket := &tyVpnProtocol.VpnPacket{
			Cmd:            tyVpnProtocol.CmdKeepAlive,
			ClientIdSender: c.clientId,
		}
		bufW := udwBytes.NewBufWriter(nil)
		const timeout = time.Second * 2
		c.closer.SleepDur(timeout / 2)
		if c.closer.IsClose(){
			return
		}
		timer := time.NewTimer(timeout)
		for {
			bufW.Reset()
			c.connLock.Lock()
			directVpnConn := c.directVpnConn
			c.connLock.Unlock()
			vpnPacket.Data = vpnPacket.Data[:0]
			vpnPacket.Encode(bufW)
			bufW.WriteLittleEndUint64(i)
			err := udwBinary.WriteByteSliceWithUint32LenNoAllocV2(directVpnConn, bufW.GetBytes())
			if c.closer.IsClose(){
				return
			}
			if err != nil {
				c.reconnect()
				continue
			}
			timer.Reset(timeout)
			select {
			case <-timer.C:
				udwLog.Log("[snc1hhr1ems1q] keepAlive timeout", i)
				c.reconnect()
			case _i := <-c.keepAliveChan:
				if _i == i {
					i++
					c.closer.SleepDur(timeout / 2)
					if c.closer.IsClose(){
						return
					}
					continue
				}
				udwLog.Log("[snc1hhr1ems1q] keepAlive error: i not match, expect", i, "but got", _i)
				c.reconnect()
			case <-c.closer.GetCloseChan():
				return
			}
		}
	}()
}

func (c *Client) reconnect() {
	c.connLock.Lock()
	if c.vpnConn != nil {
		_ = c.vpnConn.Close()
	}
	if c.directVpnConn != nil {
		_ = c.directVpnConn.Close()
	}
	c.connLock.Unlock()
	for {
		if c.closer.IsClose(){
			return
		}
		udwLog.Log("[ruu1n967nwm] RECONNECT...")
		err := c.connect()
		if err != nil {
			udwLog.Log("[ruu1n967nwm] RECONNECT Failed", err)
			time.Sleep(time.Millisecond * 500)
			continue
		}
		udwLog.Log("[ruu1n967nwm] RECONNECT Succeed ✔")
		return
	}
}
