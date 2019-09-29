// +build linux android
// +build !js

package udwNet

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwLog"
	"net"
	"os/exec"
	"strconv"
	"sync"
	"syscall"
	"unsafe"
)

type RouteRuleV2 struct {
	srcIpNet    net.IPNet
	dstIpNet    net.IPNet
	gatewayIp   net.IP
	prefSrcIp   net.IP
	inputIface  *NetDevice
	outputIface *NetDevice
	priority    uint32
	ipVersion   uint8
	tableId     uint32
	flags       uint32
	typ         byte
}

func (rule *RouteRuleV2) IsIpv4() bool {
	return rule.ipVersion == 4
}
func (rule *RouteRuleV2) String() string {
	s := "RouteRule "
	if !IPNetIsZero(rule.srcIpNet) {
		s += "srcIpNet:" + rule.srcIpNet.String() + " "
	}
	if !IPNetIsZero(rule.dstIpNet) {
		s += "dstIpNet:" + rule.dstIpNet.String() + " "
	}
	if rule.gatewayIp != nil {
		s += "gatewayIp:" + rule.gatewayIp.String() + " "
	}
	if rule.prefSrcIp != nil {
		s += "prefSrcIp:" + rule.prefSrcIp.String() + " "
	}
	if rule.inputIface != nil {
		s += "inputIface:" + rule.inputIface.GetName() + " "
	}
	if rule.outputIface != nil {
		s += "outputIface:" + rule.outputIface.GetName() + " "
	}
	if rule.priority > 0 {
		s += "priority:" + strconv.Itoa(int(rule.priority)) + " "
	}
	if rule.tableId > 0 {
		s += "tableId:" + strconv.Itoa(int(rule.tableId)) + " "
	}
	if rule.typ == syscall.RTN_LOCAL {
		s += "type:local "
	} else if rule.typ == syscall.RTN_UNICAST {
		s += "type:unicast "
	} else if rule.typ == syscall.RTN_BROADCAST {
		s += "type:broadcast "
	} else if rule.typ == syscall.RTN_UNREACHABLE {
		s += "type:unreachable "
	} else {
		s += "type:" + strconv.Itoa(int(rule.typ)) + " "
	}
	if rule.flags > 0 {
		s += routeFlagString(rule.flags) + " "
	}
	s += "ipv" + strconv.Itoa(int(rule.ipVersion))
	return s
}

func routeFlagString(flags uint32) string {
	s := "flags:"
	if (flags & syscall.RTM_F_CLONED) > 0 {
		s += "Cloned,"
		flags = flags & (^uint32(syscall.RTM_F_CLONED))
	}
	if flags > 0 {
		s += strconv.Itoa(int(flags))
	}
	return s
}

func interfaceString(dev net.Interface) string {
	s := "net.Interface "
	s += "Index:" + strconv.Itoa(dev.Index) + " "
	s += "Mtu:" + strconv.Itoa(dev.MTU) + " "
	s += "Name:" + dev.Name + " "
	s += "HardwareAddr:" + dev.HardwareAddr.String() + " "
	s += "Flags:" + dev.Flags.String() + " "
	return s
}

type routeContext struct {
}

func newRouteContext() *routeContext {
	return &routeContext{}
}

func (ctx *routeContext) mustGetRouteTable() (ruleList []*RouteRuleV2) {
	interfaceList := MustGetNetDeviceList()
	interfaceIndexMap := map[int]*NetDevice{}
	for _, dev := range interfaceList {
		interfaceIndexMap[dev.GetIndex()] = dev
	}
	ruleList = []*RouteRuleV2{}

	tab, err := syscall.NetlinkRIB(syscall.RTM_GETROUTE, syscall.AF_UNSPEC)
	if err != nil {
		panic(err)
	}
	msgs, err := syscall.ParseNetlinkMessage(tab)
	if err != nil {
		panic(err)
	}
	for _, m := range msgs {
		if m.Header.Type == syscall.NLMSG_DONE {
			break
		} else if m.Header.Type != syscall.RTM_NEWROUTE {
			udwLog.Log("MustGetRouteTable", "unhandle msg type", m.Header.Type)
			continue
		}
		thisRule := getOneRouteRule(m, interfaceIndexMap)
		if thisRule == nil {
			continue
		}
		ruleList = append(ruleList, thisRule)
	}
	return ruleList
}

func (ctx *routeContext) mustGetDefaultRouteRule() (rule *RouteRuleV2) {
	ruleList := ctx.mustGetRouteTable()
	var lastRule *RouteRuleV2
	for _, rule := range ruleList {
		if !(rule.typ == syscall.RTN_UNICAST && IPNetGetOnes(rule.GetDstIpNet()) == 0) {
			continue
		}
		lastRule = rule
		if lastRule.GetGatewayIp() != nil {
			return lastRule
		}
	}
	return lastRule
}

func getOneRouteRule(m syscall.NetlinkMessage, interfaceIndexMap map[int]*NetDevice) *RouteRuleV2 {
	rt := (*routeInfoInMemory)(unsafe.Pointer(&m.Data[0]))
	routeInfo := RouteRuleV2{}
	attrs, err := syscall.ParseNetlinkRouteAttr(&m)
	if err != nil {
		panic(err)
	}
	if rt.Family == syscall.AF_INET {
		routeInfo.ipVersion = 4
	} else if rt.Family == syscall.AF_INET6 {
		routeInfo.ipVersion = 6
	} else {
		udwLog.Log("MustGetRouteTable", "unhandle rt.Family", rt.Family)
		return nil
	}
	routeInfo.flags = rt.Flags
	routeInfo.typ = rt.Type
	for _, attr := range attrs {
		switch attr.Attr.Type {
		case syscall.RTA_DST:
			routeInfo.dstIpNet = net.IPNet{
				IP:   net.IP(attr.Value),
				Mask: net.CIDRMask(int(rt.DstLen), len(attr.Value)*8),
			}
		case syscall.RTA_SRC:
			routeInfo.srcIpNet = net.IPNet{
				IP:   net.IP(attr.Value),
				Mask: net.CIDRMask(int(rt.SrcLen), len(attr.Value)*8),
			}
		case syscall.RTA_GATEWAY:
			routeInfo.gatewayIp = net.IP(attr.Value)
		case syscall.RTA_PREFSRC:
			routeInfo.prefSrcIp = net.IP(attr.Value)
		case syscall.RTA_IIF:
			index := *(*uint32)(unsafe.Pointer(&attr.Value[0]))
			if index > 0 {
				dev, ok := interfaceIndexMap[int(index)]
				if !ok {
					panic(fmt.Errorf("[getOneRouteRule] inputIface index [%d] can not found", index))
				}
				routeInfo.inputIface = dev
			}
		case syscall.RTA_OIF:
			index := *(*uint32)(unsafe.Pointer(&attr.Value[0]))
			if index > 0 {
				dev, ok := interfaceIndexMap[int(index)]
				if !ok {
					panic(fmt.Errorf("[getOneRouteRule] outputIface index [%d] can not found", index))
				}
				routeInfo.outputIface = dev
			}
		case syscall.RTA_PRIORITY:
			routeInfo.priority = *(*uint32)(unsafe.Pointer(&attr.Value[0]))
		case syscall.RTA_METRICS:

		case syscall.RTA_TABLE:
			routeInfo.tableId = *(*uint32)(unsafe.Pointer(&attr.Value[0]))
		case syscall.RTA_CACHEINFO:
		case syscall.RTA_UNSPEC:
		case rTA_PREF:

		default:
			udwLog.Log("MustGetRouteTable", "unhandle attr.Attr.Type", attr.Attr.Type, fmt.Sprintf("%#v", attr.Value))
		}
	}
	getZeroIpNet := func(len uint8) net.IPNet {
		if routeInfo.ipVersion == 4 {
			return net.IPNet{
				IP:   net.IPv4zero,
				Mask: net.CIDRMask(int(len), 32),
			}
		} else if routeInfo.ipVersion == 6 {
			return net.IPNet{
				IP:   net.IPv4zero,
				Mask: net.CIDRMask(int(len), 128),
			}
		} else {
			panic(fmt.Errorf("[MustGetRouteTable] unexpect ipversion %d", routeInfo.ipVersion))
		}
	}
	if IPNetIsZero(routeInfo.dstIpNet) {
		routeInfo.dstIpNet = getZeroIpNet(rt.DstLen)
	}
	if IPNetIsZero(routeInfo.srcIpNet) {
		routeInfo.srcIpNet = getZeroIpNet(rt.SrcLen)
	}

	return &routeInfo
}

const rTA_PREF = 20

type routeInfoInMemory struct {
	Family byte
	DstLen byte
	SrcLen byte
	TOS    byte

	Table    byte
	Protocol byte
	Scope    byte
	Type     byte

	Flags uint32
}

var gMustRouteSetCmdRoutePath = ""
var gMustRouteSetCmdIpPath = ""
var gMustRouteSetCmdOnce sync.Once

func mustRouteSetCmdInit() {
	gMustRouteSetCmdOnce.Do(func() {
		path, err := exec.LookPath("ip")
		if err == nil {
			gMustRouteSetCmdIpPath = path
			return
		}
		path, err = exec.LookPath("route")
		if err == nil {
			gMustRouteSetCmdRoutePath = path
			return
		}
	})
}

func (ctx *routeContext) mustRouteSet(ipNet net.IPNet, gateWayIp net.IP) {
	mustRouteSetCmdInit()
	if gMustRouteSetCmdIpPath != "" {
		udwCmd.CmdSlice([]string{"ip", "route", "add", ipNet.String(), "via", gateWayIp.String()}).
			MustCombinedOutputAndNotExitStatusCheck()
		return
	}
	if gMustRouteSetCmdRoutePath != "" {
		ones, bits := ipNet.Mask.Size()
		netOrHostS := "-net"
		if ones == bits {
			netOrHostS = "-host"
		}
		udwCmd.CmdSlice([]string{"route", "-n", "add", netOrHostS, ipNet.String(), "gateway", gateWayIp.String()}).
			MustCombinedOutputAndNotExitStatusCheck()
		return
	}
	panic(`Can not finud "ip" or "route", you can install them or check your $PATH`)
}
func (ctx *routeContext) mustRouteDelete(ipNet net.IPNet) {
	mustRouteSetCmdInit()
	if gMustRouteSetCmdIpPath != "" {
		udwCmd.MustCombinedOutputAndNotExitStatusCheck("ip route del " + ipNet.String())
		return
	}
	if gMustRouteSetCmdRoutePath != "" {
		ones, bits := ipNet.Mask.Size()
		if ones == bits {
			udwCmd.MustCombinedOutputAndNotExitStatusCheck("route -n delete -host " + ipNet.String())
		} else {
			udwCmd.MustCombinedOutputAndNotExitStatusCheck("route -n delete -net " + ipNet.String())
		}
		return
	}
	panic(`Can not finud "ip" or "route", you can install them or check your $PATH`)
}
