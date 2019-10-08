// +build darwin,ios darwin,macAppStore

package udwNet

import (
	"fmt"
	"strings"
)

func HasNetworkOnDevice() (hasNetwork bool, err error) {
	list, err := GetNetDeviceList()
	if err != nil {
		return false, fmt.Errorf("[HasNetworkOnDevice] %s", err.Error())
	}
	isEthernetOn := false
	isCellularOn := false
	for _, dev := range list {
		if !dev.CanConnect() {
			continue
		}
		if dev.GetName() == "pdp_ip0" {
			isCellularOn = true
		}
		if strings.HasPrefix(dev.GetName(), "en") {
			isEthernetOn = true
		}
	}
	return isEthernetOn || isCellularOn, nil
}
