// +build !js

package udwNet

import (
	"encoding/binary"
	"fmt"
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwLog"
	"net"
	"strconv"
	"syscall"
	"unsafe"
)

const syscallDebug = false

func newRouteContext() *routeContext {
	return &routeContext{}
}

func MustRouteDeleteWithOutputIfIndex(ifIndex int) {
	ctx := newRouteContext()
	ctx.MustRouteDeleteWithOutputIfIndex(ifIndex)
}

type routeContext struct {
	netDeviceList   []*NetDevice
	routeWinRowList []mIB_IPFORWARDROW
	routeList       []*RouteRuleV2
}

func (ctx *routeContext) MustRouteSetWithIfIndex(ipnet net.IPNet, gatewayIp net.IP, IfIndex int) {

	row := mIB_IPFORWARDROW{
		DwForwardDest:    uint32FromIpv4(ipnet.IP),
		DwForwardMask:    uint32FromIpMask(ipnet.Mask),
		DwForwardPolicy:  0,
		DwForwardNextHop: uint32FromIpv4(gatewayIp),
		DwForwardIfIndex: uint32(IfIndex),
		DwForwardType:    MIB_IPROUTE_TYPE_INDIRECT,
		DwForwardProto:   MIB_IPPROTO_NETMGMT,
		DwForwardMetric1: 30,
	}
	mustSetIpForwardEntry(&row)
}

func (ctx *routeContext) mustRouteSet(ipNet net.IPNet, gatewayIP net.IP) {
	ctx.ensureContextExist()

	ifIndex := 0
	for _, dev := range ctx.netDeviceList {
		for _, ipnet := range dev.ipNetList {
			if ipnet.Contains(gatewayIP) {
				ifIndex = dev.GetIndex()
				break
			}
		}
		if ifIndex > 0 {
			break
		}
	}
	if ifIndex == 0 {
		panic(fmt.Errorf("[MustRouteSet] can not get dev from gatewayIp %s", gatewayIP.String()))
	}
	udwCmd.CmdSlice([]string{"route", "add", ipNet.IP.String(),
		"mask", IPNetGetGenMaskString(&ipNet), gatewayIP.String(),
		"metric", "30", "if", strconv.Itoa(ifIndex)}).MustRun()

}

func (ctx *routeContext) MustRouteDeleteWithOutputIfIndex(ifIndex int) {
	ctx.ensureContextExist()
	for _, row := range ctx.routeWinRowList {
		if int(row.DwForwardIfIndex) == ifIndex {
			err := syscallDeleteIpForwardEntry(&row)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (ctx *routeContext) mustRouteDelete(ipnet net.IPNet) {
	ctx.ensureContextExist()
	inputOnes := IPNetGetOnes(ipnet)
	table := ctx.routeWinRowList
	for _, row := range table {
		if !ipv4FromUint32(row.DwForwardDest).Equal(ipnet.IP) {
			continue
		}
		if IPMaskGetOnes(ipMaskFromUint32(row.DwForwardMask)) != inputOnes {
			continue
		}
		err := syscallDeleteIpForwardEntry(&row)
		if err != nil {
			panic(err)
		}
		return
	}
}

func (ctx *routeContext) mustGetRouteTable() []*RouteRuleV2 {
	ctx.netDeviceList = MustGetNetDeviceList()
	interfaceIndexMap := map[int]*NetDevice{}
	for _, dev := range ctx.netDeviceList {
		interfaceIndexMap[dev.GetIndex()] = dev
	}
	rowList := mustSyscallGetIpForwardTable()
	ctx.routeWinRowList = rowList
	outRuleList := make([]*RouteRuleV2, len(rowList))
	for i, row := range rowList {
		thisRule := &RouteRuleV2{}
		thisRule.dstIpNet = net.IPNet{
			IP:   ipv4FromUint32(row.DwForwardDest),
			Mask: ipMaskFromUint32(row.DwForwardMask),
		}
		thisRule.gatewayIp = ipv4FromUint32(row.DwForwardNextHop)

		outputIface, ok := interfaceIndexMap[int(row.DwForwardIfIndex)]
		if !ok {
			panic(fmt.Errorf("[mustParseDarwinRouteV2] can not found interface index %d", row.DwForwardIfIndex))
		}
		thisRule.outputIface = outputIface
		thisRule.originRow = row
		outRuleList[i] = thisRule
	}
	ctx.routeList = outRuleList
	return outRuleList
}

func (ctx *routeContext) mustGetDefaultRouteRule() (rule *RouteRuleV2) {
	routeList := ctx.mustGetRouteTable()
	for _, route := range routeList {
		if IPNetGetOnes(route.dstIpNet) == 0 {
			return route
		}
	}
	panic(fmt.Errorf("[mustGetDefaultRouteRule] can not found default route"))
}

func (ctx *routeContext) ensureContextExist() {
	if ctx.netDeviceList == nil || ctx.routeList == nil {
		ctx.mustGetRouteTable()
	}
}

var (
	modiphlpapi = syscall.NewLazyDLL("iphlpapi.dll")

	procGetIpForwardTable    = modiphlpapi.NewProc("GetIpForwardTable")
	procCreateIpForwardEntry = modiphlpapi.NewProc("CreateIpForwardEntry")
	procSetIpForwardEntry    = modiphlpapi.NewProc("SetIpForwardEntry")
	procDeleteIpForwardEntry = modiphlpapi.NewProc("DeleteIpForwardEntry")
)

const MIB_IPPROTO_NETMGMT = 3
const MIB_IPROUTE_TYPE_DIRECT = 3
const MIB_IPROUTE_TYPE_INDIRECT = 4

const ERROR_OBJECT_ALREADY_EXISTS = 5010

type mIB_IPFORWARDROW struct {
	DwForwardDest      uint32
	DwForwardMask      uint32
	DwForwardPolicy    uint32
	DwForwardNextHop   uint32
	DwForwardIfIndex   uint32
	DwForwardType      uint32
	DwForwardProto     uint32
	DwForwardAge       uint32
	DwForwardNextHopAS uint32
	DwForwardMetric1   uint32
	DwForwardMetric2   uint32
	DwForwardMetric3   uint32
	DwForwardMetric4   uint32
	DwForwardMetric5   uint32
}

func ipv4FromUint32(ipu32 uint32) net.IP {
	return net.IPv4(uint8(ipu32), uint8(ipu32>>8), uint8(ipu32>>16), uint8(ipu32>>24))
}
func uint32FromIpv4(ip net.IP) uint32 {
	ip = ip.To4()
	return binary.LittleEndian.Uint32([]byte(ip))
}
func ipMaskFromUint32(ipu32 uint32) net.IPMask {
	return net.IPv4Mask(uint8(ipu32), uint8(ipu32>>8), uint8(ipu32>>16), uint8(ipu32>>24))
}
func uint32FromIpMask(ipMask net.IPMask) uint32 {
	return binary.LittleEndian.Uint32([]byte(ipMask))
}

func mustSyscallGetIpForwardTable() (rowList []mIB_IPFORWARDROW) {
	var buf []byte
	bufSize := 4096
	for {
		buf = make([]byte, bufSize)
		outSize := uint32(len(buf))
		r0, _, _ := syscall.Syscall(procGetIpForwardTable.Addr(), 3, uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&outSize)),
			uintptr(0))
		if r0 == 0 {
			break
		}
		err := syscall.Errno(r0)
		if err == syscall.ERROR_INSUFFICIENT_BUFFER {
			bufSize = int(outSize)
			continue
		}
		panic(err)
	}
	dwNumEntries := binary.LittleEndian.Uint32(buf[:4])
	bufRowList := buf[4:]
	rowList = make([]mIB_IPFORWARDROW, int(dwNumEntries))
	copy(rowList, *((*[]mIB_IPFORWARDROW)(unsafe.Pointer(&bufRowList))))
	if syscallDebug {
		for i := range rowList {
			udwLog.Log("debug", "[mustSyscallGetIpForwardTable]", rowList[i])
		}
	}
	return rowList
}

func syscallCreateIpForwardEntry(row *mIB_IPFORWARDROW) (err error) {
	r0, _, _ := syscall.Syscall(procCreateIpForwardEntry.Addr(), 1, uintptr(unsafe.Pointer(row)), 0, 0)
	if syscallDebug {
		udwLog.Log("debug", "[syscallCreateIpForwardEntry]", row, r0, syscall.Errno(r0).Error())
	}
	if r0 == 0 {
		return nil
	}
	return syscall.Errno(r0)
}

func syscallSetIpForwardEntry(row *mIB_IPFORWARDROW) (err error) {
	r0, _, _ := syscall.Syscall(procSetIpForwardEntry.Addr(), 1, uintptr(unsafe.Pointer(row)), 0, 0)
	if syscallDebug {
		udwLog.Log("debug", "[syscallSetIpForwardEntry]", row, r0, syscall.Errno(r0).Error())
	}
	if r0 == 0 {
		return nil
	}
	return syscall.Errno(r0)
}

func mustSetIpForwardEntry(r *mIB_IPFORWARDROW) {
	err := syscallCreateIpForwardEntry(r)
	if err == nil {
		return
	}
	if err.(syscall.Errno) != ERROR_OBJECT_ALREADY_EXISTS {
		panic(err)
	}
	err = syscallSetIpForwardEntry(r)
	if err != nil {
		panic(err)
	}
}

func syscallDeleteIpForwardEntry(row *mIB_IPFORWARDROW) (err error) {
	r0, _, _ := syscall.Syscall(procDeleteIpForwardEntry.Addr(), 1, uintptr(unsafe.Pointer(row)), 0, 0)
	if syscallDebug {
		udwLog.Log("debug", "[syscallDeleteIpForwardEntry]", row, r0, syscall.Errno(r0).Error())
	}
	if r0 == 0 {
		return nil
	}
	if r0 == 1168 {
		return nil
	}
	return syscall.Errno(r0)
}

type RouteRuleV2 struct {
	dstIpNet  net.IPNet
	gatewayIp net.IP

	outputIface *NetDevice

	originRow mIB_IPFORWARDROW
}

func (rule *RouteRuleV2) IsIpv4() bool {
	return rule.gatewayIp.To4() != nil
}

func (rule *RouteRuleV2) IsDwForwardTypeDIRECT() bool {
	return rule.originRow.DwForwardType == MIB_IPROUTE_TYPE_DIRECT
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
	if rule.originRow.DwForwardType != 0 {
		s += "DwForwardType:" + dwForwardTypeToString(rule.originRow.DwForwardType) + " "
	}
	if rule.originRow.DwForwardProto != 0 {
		s += "DwForwardProto:" + dwForwardProtoToString(rule.originRow.DwForwardProto) + " "
	}
	if rule.originRow.DwForwardAge != 0 {
		s += "DwForwardAge:" + strconv.Itoa(int(rule.originRow.DwForwardAge)) + " "
	}
	if rule.originRow.DwForwardMetric1 != 0 {
		s += "DwForwardMetric1:" + strconv.Itoa(int(rule.originRow.DwForwardMetric1)) + " "
	}

	if rule.originRow.DwForwardPolicy != 0 {
		s += "DwForwardPolicy:" + strconv.Itoa(int(rule.originRow.DwForwardPolicy)) + " "
	}
	if rule.originRow.DwForwardNextHopAS != 0 {
		s += "DwForwardNextHopAS:" + strconv.Itoa(int(rule.originRow.DwForwardNextHopAS)) + " "
	}
	if rule.originRow.DwForwardMetric2 != 0 {
		s += "DwForwardMetric2:" + strconv.Itoa(int(rule.originRow.DwForwardMetric2)) + " "
	}
	if rule.originRow.DwForwardMetric3 != 0 {
		s += "DwForwardMetric2:" + strconv.Itoa(int(rule.originRow.DwForwardMetric3)) + " "
	}
	if rule.originRow.DwForwardMetric4 != 0 {
		s += "DwForwardMetric2:" + strconv.Itoa(int(rule.originRow.DwForwardMetric4)) + " "
	}
	if rule.originRow.DwForwardMetric5 != 0 {
		s += "DwForwardMetric2:" + strconv.Itoa(int(rule.originRow.DwForwardMetric5)) + " "
	}
	return s
}

func dwForwardTypeToString(typ uint32) string {
	switch typ {
	case 1:
		return "MIB_IPROUTE_TYPE_OTHER"
	case 2:
		return "MIB_IPROUTE_TYPE_INVALID"
	case 3:
		return "MIB_IPROUTE_TYPE_DIRECT"
	case 4:
		return "MIB_IPROUTE_TYPE_INDIRECT"
	default:
		return strconv.Itoa(int(typ))
	}
}

func dwForwardProtoToString(typ uint32) string {
	switch typ {
	case 1:
		return "MIB_IPPROTO_OTHER"
	case 2:
		return "MIB_IPPROTO_LOCAL"
	case 3:
		return "MIB_IPPROTO_NETMGMT"
	case 4:
		return "MIB_IPPROTO_ICMP"
	case 5:
		return "MIB_IPPROTO_EGP"
	case 6:
		return "MIB_IPPROTO_GGP"
	case 7:
		return "MIB_IPPROTO_HELLO"
	case 8:
		return "MIB_IPPROTO_RIP"
	case 9:
		return "MIB_IPPROTO_IS_IS"
	case 10:
		return "MIB_IPPROTO_ES_IS"
	case 11:
		return "MIB_IPPROTO_CISCO"
	case 12:
		return "MIB_IPPROTO_BBN"
	case 13:
		return "MIB_IPPROTO_OSPF"
	case 14:
		return "MIB_IPPROTO_BGP"
	case 10002:
		return "MIB_IPPROTO_NT_AUTOSTATIC"
	case 10006:
		return "MIB_IPPROTO_NT_STATIC"
	case 10007:
		return "MIB_IPPROTO_NT_STATIC_NON_DOD"
	default:
		return strconv.Itoa(int(typ))
	}
}
