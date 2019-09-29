// +build !js

package udwNet

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwErr"
	"io"
	"net"
	"os"
)

func RouteTableDebugPrint() {
	RouteTableDebugPrintToWriter(os.Stdout)
}

func RouteTableDebugPrintToWriter(w io.Writer) {
	udwErr.PanicToErrorAndLog(func() {
		fmt.Fprintln(w, "MustGetNetDeviceList()======================")
		interfaceList := MustGetNetDeviceList()
		for _, dev := range interfaceList {
			fmt.Fprintln(w, dev.String())
		}
		fmt.Fprintln(w, "MustGetRouteTable()======================")
		routeList := MustGetRouteTable()
		for _, route := range routeList {
			fmt.Fprintln(w, route.String())
		}
		fmt.Fprintln(w, "MustGetDefaultRouteRuleInterfaceFirstIpAddr()======================")
		fmt.Fprintln(w, MustGetDefaultRouteRuleInterfaceFirstIpAddr())
		fmt.Fprintln(w, "MustGetDefaultGatewayIP()======================")
		fmt.Fprintln(w, MustGetDefaultGatewayIP())
	})
}

func (rule *RouteRuleV2) GetDstIpNet() net.IPNet {
	return rule.dstIpNet
}
func (rule *RouteRuleV2) GetGatewayIp() net.IP {
	return rule.gatewayIp
}
func (rule *RouteRuleV2) GetOutInterface() *NetDevice {
	return rule.outputIface
}

func MustGetRouteTable() (ruleList []*RouteRuleV2) {
	return NewRouteContext().MustGetRouteTable()
}

func MustGetDefaultRouteRule() (rule *RouteRuleV2) {
	return NewRouteContext().MustGetDefaultRouteRule()
}

func MustGetDefaultGatewayIP() net.IP {
	rule := MustGetDefaultRouteRule()
	return rule.GetGatewayIp()
}

func MustGetDefaultRouteRuleInterfaceFirstIpAddr() net.IP {
	rule := MustGetDefaultRouteRule()
	if rule == nil {

		panic("can not get default route rule, ENETUNREACH")
	}
	iface := rule.GetOutInterface()

	return iface.GetFirstIpv4IP()
}

func MustRouteSetWithString(IpNetString string, gatewayString string) {
	NewRouteContext().MustRouteSetWithString(IpNetString, gatewayString)
}

func MustRouteDeleteWithString(IpNetString string) {
	NewRouteContext().MustRouteDeleteWithString(IpNetString)
}

type RouteContext struct {
	rc *routeContext
}

func NewRouteContext() RouteContext {
	return RouteContext{rc: newRouteContext()}
}
func (ctx RouteContext) MustGetRouteTable() (ruleList []*RouteRuleV2) {
	return ctx.rc.mustGetRouteTable()
}
func (ctx RouteContext) MustGetDefaultRouteRule() (rule *RouteRuleV2) {
	return ctx.rc.mustGetDefaultRouteRule()
}
func (ctx RouteContext) MustRouteSet(ipNet net.IPNet, gateWayIp net.IP) {
	ctx.rc.mustRouteSet(ipNet, gateWayIp)
}
func (ctx RouteContext) MustRouteDelete(ipNet net.IPNet) {
	ctx.rc.mustRouteDelete(ipNet)
}
func (ctx RouteContext) MustRouteSetWithString(IpNetString string, gatewayString string) {
	_, ipnet, err := net.ParseCIDR(IpNetString)
	if err != nil {
		panic(fmt.Errorf("[MustRouteSetWithString] net.ParseCIDR IpNetString [%s] fail %s", IpNetString, err.Error()))
	}
	gateWayIp := net.ParseIP(gatewayString)
	if gateWayIp == nil {
		panic(fmt.Errorf("[MustRouteSetWithString] can not parse gatewayIp [%s]", gatewayString))
	}
	ctx.rc.mustRouteSet(*ipnet, gateWayIp)
}
func (ctx RouteContext) MustRouteDeleteWithString(IpNetString string) {
	_, ipnet, err := net.ParseCIDR(IpNetString)
	if err != nil {
		panic(fmt.Errorf("[MustRouteDeleteWithString] net.ParseCIDR IpNetString [%s] fail %s", IpNetString, err.Error()))
	}
	ctx.rc.mustRouteDelete(*ipnet)
}
