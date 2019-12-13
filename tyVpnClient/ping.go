package tyVpnClient

import (
	"crypto/tls"
	"errors"
	"github.com/tachyon-protocol/udw/tyTls"
	"github.com/tachyon-protocol/udw/tyVpnProtocol"
	"github.com/tachyon-protocol/udw/udwBinary"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwLog"
	"net"
	"strconv"
	"time"
)

type PingReq struct {
	Ip        string
	ServerChk string // if it is "", it will use InsecureSkipVerify
	Count     int
	DebugLog  bool
}

//TODO relay mode
func Ping(req PingReq) error {
	var tlsConfig *tls.Config
	if req.ServerChk == "" {
		tlsConfig = newInsecureClientTlsConfig()
	} else {
		var errMsg string
		tlsConfig, errMsg = tyTls.NewClientTlsConfigWithChk(tyTls.NewClientTlsConfigWithChkReq{
			ServerChk: req.ServerChk,
		})
		if errMsg != "" {
			return errors.New(errMsg)
		}
	}
	conn, err := net.Dial("tcp", req.Ip+":"+strconv.Itoa(tyVpnProtocol.VpnPort))
	if err != nil {
		return err
	}
	conn = tls.Client(conn, tlsConfig)
	var (
		pingPacket = tyVpnProtocol.VpnPacket{
			Cmd: tyVpnProtocol.CmdKeepAlive,
		}
		buf = udwBytes.NewBufWriter(nil)
	)
	for i := 0; i < req.Count; i++ {
		buf.Reset()
		pingPacket.Encode(buf)
		start := time.Now()
		err = udwBinary.WriteByteSliceWithUint32LenNoAllocV2(conn, buf.GetBytes())
		if err != nil {
			return err
		}
		if req.DebugLog {
			udwLog.Log("-> ...")
		}
		buf.Reset()
		err := udwBinary.ReadByteSliceWithUint32LenToBufW(conn, buf)
		if err != nil {
			return err
		}
		err = pingPacket.Decode(buf.GetBytes())
		if err != nil {
			return err
		}
		if pingPacket.Cmd!=tyVpnProtocol.CmdKeepAlive{
			return errors.New("tvhfx76ynx")
		}
		if req.DebugLog {
			udwLog.Log("<- ✔", time.Now().Sub(start))
		}
	}
	return nil
}
