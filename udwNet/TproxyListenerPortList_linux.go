// +build linux,!js

package udwNet

import (
	"strconv"
	"strings"
)

func TproxyListenerPortWorkAround(addr string) []string {
	if strings.HasPrefix(addr, ":") {
		part := strings.Split(addr, ":")
		i, err := strconv.Atoi(part[1])
		if err != nil {
			panic(err)
		}
		return MustGetCurrentIpWithPortList(uint16(i))
	}
	return []string{addr}
}
