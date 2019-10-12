package udwTapTun

import (
	"bytes"
	"errors"
	"github.com/tachyon-protocol/udw/udwCmd"
	"net"
	"os/exec"
	"runtime"
	"strconv"
)

type SetP2PIpRequest struct {
	IfaceName string
	SrcIp     net.IP
	DstIp     net.IP
	Mtu       int
	Mask      net.IPMask
}

func SetP2PIpAndUp(req SetP2PIpRequest) error {
	switch runtime.GOOS {
	case "darwin":
		cmdSlice := []string{"ifconfig", req.IfaceName, req.SrcIp.String(), req.DstIp.String()}
		if req.Mask != nil {
			cmdSlice = append(cmdSlice, "netmask", netmaskDotListString(req.Mask))
		}
		if req.Mtu > 0 {
			cmdSlice = append(cmdSlice, "mtu", strconv.Itoa(req.Mtu))
		}
		cmdSlice = append(cmdSlice, "up")
		err := udwCmd.StdioSliceRun(cmdSlice)
		if err != nil {
			return err
		}
		return nil
	case "linux":
		path, err := exec.LookPath("ip")
		if err == nil {
			err = udwCmd.StdioSliceRun([]string{path, "link", "set", req.IfaceName, "up", "mtu", strconv.Itoa(req.Mtu)})
			if err != nil {
				return err
			}
			addrNet := net.IPNet{
				IP:   req.SrcIp,
				Mask: req.Mask,
			}
			err = udwCmd.StdioSliceRun([]string{path, "addr", "add", addrNet.String(), "dev", req.IfaceName})
			if err != nil {
				return err
			}
			return nil
		}
		path, err = exec.LookPath("ifconfig")
		if err == nil {
			cmdSlice := []string{path, req.IfaceName, req.SrcIp.String(), "pointopoint", req.DstIp.String()}
			if req.Mask != nil {
				cmdSlice = append(cmdSlice, "netmask", netmaskDotListString(req.Mask))
			}
			if req.Mtu > 0 {
				cmdSlice = append(cmdSlice, "mtu", strconv.Itoa(req.Mtu))
			}
			cmdSlice = append(cmdSlice, "up")
			return udwCmd.StdioSliceRun(cmdSlice)
		}
		return errors.New(`Can not find binary "ifconfig" or "ip", install them or set currect $PATH may solve this.`)
	default:
		return GetErrPlatformNotSupport()
	}
}

func netmaskDotListString(mask net.IPMask) string {
	buf := &bytes.Buffer{}
	for i, b := range mask {
		buf.WriteString(strconv.Itoa(int(b)))
		if i != len(mask)-1 {
			buf.WriteString(".")
		}
	}
	return buf.String()
}

func GetErrPlatformNotSupport() error {
	return errors.New("Platform Not Support")
}
