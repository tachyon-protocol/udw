// +build !darwin,!linux,!windows,!js

package udwNet

type routeContext struct {
}

func newRouteContext() *routeContext {
	panic("Not Implement")
}

func (ctx *routeContext) mustRouteSet(IpNetString string, gatewayString string) {
	panic("Not Implement")
}

func (ctx *routeContext) mustRouteDelete(IpNetString string) {
	panic("Not Implement")
}

func (ctx *routeContext) mustGetRouteTable() (ruleList []*RouteRuleV2) {
	panic("Not Implement")
}

func (ctx *routeContext) mustGetDefaultRouteRule() (rule *RouteRuleV2) {
	panic("Not Implement")
}
