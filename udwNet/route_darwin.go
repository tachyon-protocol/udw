// +build !js

package udwNet

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwLog"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

type routeContext struct {
}

func newRouteContext() *routeContext {
	return &routeContext{}
}

func (ctx *routeContext) mustGetRouteTable() (ruleList []*RouteRuleV2) {
	return ctx.mustGetRouteTableWithSyscall()
}

const ErrStringNoDefaultRoute = `can not found default route`

func (ctx *routeContext) mustGetDefaultRouteRule() (rule *RouteRuleV2) {
	routeList := ctx.mustGetRouteTable()
	for _, route := range routeList {
		if IPNetGetOnes(route.dstIpNet) == 0 {
			return route
		}
	}
	panic("[mustGetDefaultRouteRule] " + ErrStringNoDefaultRoute)
}

func (ctx *routeContext) mustRouteSet(ipNet net.IPNet, gateWayIp net.IP) {
	err := routeCmd(ipNet, gateWayIp, syscall.RTM_ADD)
	if err != nil {
		if os.IsExist(err) {
			routeCmd(ipNet, gateWayIp, syscall.RTM_CHANGE)
			return
		}
		panic(err)
	}
}

func (ctx *routeContext) mustRouteDelete(ipNet net.IPNet) {
	err := routeCmd(ipNet, nil, syscall.RTM_DELETE)
	if err != nil {
		if strings.Contains(err.Error(), "no such process") {
			return
		}
		panic(err)
	}
}

var _zero uintptr

func syscallSysctl(mib []int32, old *byte, oldlen *uintptr, new *byte, newlen uintptr) (err error) {
	var _p0 unsafe.Pointer
	if len(mib) > 0 {
		_p0 = unsafe.Pointer(&mib[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	_, _, e1 := syscall.Syscall6(syscall.SYS___SYSCTL, uintptr(_p0), uintptr(len(mib)),
		uintptr(unsafe.Pointer(old)), uintptr(unsafe.Pointer(oldlen)), uintptr(unsafe.Pointer(new)), uintptr(newlen))
	if e1 == 0 {
		return nil
	}
	return e1
}

func (ctx *routeContext) mustGetRouteTableWithSyscall() (ruleList []*RouteRuleV2) {
	debugGetRouteTableWithSyscall := false
	var err error
	for i := 0; i < 3; i++ {
		err = udwErr.PanicToError(func() {
			goNetDeviceList := MustGetNetDeviceList()
			interfaceIndexMap := map[int]*NetDevice{}
			for _, dev := range goNetDeviceList {
				interfaceIndexMap[dev.GetIndex()] = dev
			}
			mib := []int32{
				syscall.CTL_NET,
				syscall.AF_ROUTE,
				0,
				0,
				syscall.NET_RT_DUMP,
				0,
			}
			buflen := uintptr(0)
			err := syscallSysctl(mib, nil, &buflen, nil, 0)
			if err != nil {
				panic(err)
			}
			buf := make([]byte, buflen)
			err = syscallSysctl(mib, &(buf[0]), &buflen, nil, 0)
			if err != nil {
				panic(err)
			}
			ruleList = parseSyscallRouteTable(buf, interfaceIndexMap)
		})
		if err == nil {
			return ruleList
		}
		if debugGetRouteTableWithSyscall {
			udwLog.Log("tmpError", "[routeContext.mustGetRouteTableWithSyscall]", err.Error())
		}
	}
	panic(err)
}

func parseSyscallRouteTable(buf []byte, interfaceIndexMap map[int]*NetDevice) (ruleList []*RouteRuleV2) {
	for {
		if len(buf) == 0 {
			break
		}
		size := int(buf[0]) + (int(buf[1]) << 8)
		thisPart := buf[:size]
		thisRule := parseSyscallOneRouteRow(thisPart, interfaceIndexMap)
		buf = buf[size:]
		if thisRule != nil {
			ruleList = append(ruleList, thisRule)
		}
	}
	return ruleList
}

func parseSyscallOneRouteRow(thisPart []byte, interfaceIndexMap map[int]*NetDevice) (thisRule *RouteRuleV2) {

	flag := binary.LittleEndian.Uint32(thisPart[8:12])

	index := binary.LittleEndian.Uint16(thisPart[4:6])

	thisRule = &RouteRuleV2{}
	thisRule.flags = flag
	dev, ok := interfaceIndexMap[int(index)]
	if !ok {

		panic(fmt.Errorf("[mustGetRouteTableWithSyscall] can not find dev index %d (may happen because of netDevice change and data race)", index))
	}
	thisRule.outputIface = dev

	saPart := thisPart[0x5c:]
	saBufList := [][]byte{}
	for {
		if len(saPart) == 0 {
			break
		}
		size := int(saPart[0])
		if size == 0 {
			break
		}
		thisSaPart := saPart[:size]
		saBufList = append(saBufList, thisSaPart)
		saPart = saPart[size:]
	}

	for _, buf := range saBufList {
		if buf[1] == syscall.AF_LINK {
			return nil
		}
		if buf[1] == syscall.AF_INET6 {
			return nil
		}
	}

	if len(saBufList) == 3 && saBufList[0][1] == 2 && saBufList[1][1] == 2 && saBufList[2][1] == 0xff {
		ip := net.IP(saBufList[0][4:8])
		gateway := net.IP(saBufList[1][4:8])
		IpmaskBuf := make([]byte, 4)
		copy(IpmaskBuf, saBufList[2][4:])
		netmask := net.IPMask(IpmaskBuf)
		thisRule.dstIpNet = net.IPNet{
			IP:   ip,
			Mask: netmask,
		}
		thisRule.gatewayIp = gateway
		return thisRule
	}
	if len(saBufList) == 2 && saBufList[0][1] == 2 && saBufList[1][1] == 2 {
		ip := net.IP(saBufList[0][4:8])
		gateway := net.IP(saBufList[1][4:8])
		var netmask net.IPMask
		if flag&syscall.RTF_HOST > 0 {
			netmask = net.CIDRMask(32, 32)
		} else {
			netmask = net.CIDRMask(0, 32)
		}
		thisRule.dstIpNet = net.IPNet{
			IP:   ip,
			Mask: netmask,
		}
		thisRule.gatewayIp = gateway
		return thisRule
	}
	panic(fmt.Errorf("[parseSyscallOneRouteRow] unsolve sa %s", hex.Dump(thisPart[0x5c:])))

}

type flagDataRow struct {
	valueInt    uint32
	valueString string
}

func routeFlagToString(flag uint32) string {
	out := "routeFlag:"
	flagData := []flagDataRow{
		{syscall.RTF_UP, "UP"},
		{syscall.RTF_GATEWAY, "GATEWAY"},
		{syscall.RTF_HOST, "HOST"},
		{syscall.RTF_REJECT, "REJECT"},
		{syscall.RTF_CLONING, "CLONING"},
		{syscall.RTF_DYNAMIC, "DYNAMIC"},
		{syscall.RTF_MODIFIED, "MODIFIED"},
		{syscall.RTF_DONE, "DONE"},
		{syscall.RTF_DELCLONE, "DELCLONE"},
		{syscall.RTF_XRESOLVE, "XRESOLVE"},
		{syscall.RTF_LLINFO, "LLINFO"},
		{syscall.RTF_STATIC, "STATIC"},
		{0x2000, "NOIFREF"},
		{syscall.RTF_PRCLONING, "PRCLONING"},
		{syscall.RTF_WASCLONED, "WASCLONED"},
		{syscall.RTF_PROTO3, "PROTO3"},
		{syscall.RTF_LOCAL, "LOCAL"},
		{syscall.RTF_MULTICAST, "MULTICAST"},
		{syscall.RTF_IFSCOPE, "IFSCOPE"},
		{syscall.RTF_IFREF, "IFREF"},
		{0x10000000, "ROUTER"},
	}
	for _, row := range flagData {
		if flag&row.valueInt > 0 {
			flag = flag ^ row.valueInt
			out += row.valueString + " "
		}
	}
	if flag > 0 {
		out += fmt.Sprintf("%08x", flag)
	}
	return out
}

type RouteRuleV2 struct {
	dstIpNet  net.IPNet
	gatewayIp net.IP

	outputIface *NetDevice

	flags uint32
}

func (rule *RouteRuleV2) IsIpv4() bool {
	return true
}

func (rule *RouteRuleV2) String() string {
	s := "RouteRule "
	if !IPNetIsZero(rule.dstIpNet) {
		s += "dstIpNet:" + rule.dstIpNet.String() + " "
	}
	if rule.gatewayIp != nil {
		s += "gatewayIp:" + rule.gatewayIp.String() + " "
	}
	if rule.outputIface != nil {
		s += "outputIfaceName:[" + rule.outputIface.GetName() + "] "
		s += "outputIfaceIndex:" + strconv.Itoa(rule.outputIface.GetIndex()) + " "
	}
	if rule.flags > 0 {
		s += routeFlagToString(rule.flags) + " "
	}
	return s
}

var encodeRouteAddMsgTpl = []byte{
	0x84, 0x00,
	0x05,
	0x01,
	0x00, 0x00,
	0x00, 0x00,
	0x03, 0x08, 0x00, 0x00,
	0x07, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x01, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,

	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,

	0x10,
	0x02,
	0x00, 0x00,
	0x72, 0x72, 0x72, 0x72,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

	0x10,
	0x02,
	0x00, 0x00,
	0xac, 0x15, 0x00, 0x01,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

	0x08,
	0x00,
	0x00, 0x00,
	0xff, 0xff, 0xff, 0xff,
}

func encodeRouteAddMsg(ipNet net.IPNet, gateWayIp net.IP) []byte {

	ip := ipNet.IP
	size := IPNetGetOnes(ipNet)
	maskBuf := make([]byte, 4)
	maskValue := uint64((1<<32 - 1) << (uint64(32 - size)))
	binary.BigEndian.PutUint32(maskBuf, uint32(maskValue))
	buf := make([]byte, len(encodeRouteAddMsgTpl))
	copy(buf, encodeRouteAddMsgTpl)
	copy(buf[96:100], []byte(ip.To4()))
	copy(buf[112:116], []byte(gateWayIp.To4()))
	copy(buf[128:132], []byte(maskBuf))
	return buf
}

type routeSetContext struct {
	routeFd int
	locker  sync.Mutex
}

var gRouteContext = routeSetContext{routeFd: -1}

func routeCmd(ipNet net.IPNet, gateWayIp net.IP, cmd byte) (err error) {
	msgByte := encodeRouteAddMsg(ipNet, gateWayIp)
	msgByte[3] = cmd

	gRouteContext.locker.Lock()
	defer gRouteContext.locker.Unlock()

	if gRouteContext.routeFd == -1 {
		gRouteContext.routeFd, err = syscall.Socket(syscall.AF_ROUTE, syscall.SOCK_RAW, 0)
		if err != nil {
			return err
		}
		go func() {

			time.Sleep(5 * time.Second)
			gRouteContext.locker.Lock()
			defer gRouteContext.locker.Unlock()
			if gRouteContext.routeFd != -1 {
				syscall.Close(gRouteContext.routeFd)
				gRouteContext.routeFd = -1
			}
		}()
	}
	_, err = syscall.Write(gRouteContext.routeFd, msgByte)
	if err != nil {
		return err
	}
	return nil
}
