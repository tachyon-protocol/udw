// +build windows

package udwTapTun

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwWindowsRegistry"
)

type TagReg struct {
	Guid string
}

type PanelReg struct {
	Name string
	Guid string
}

const tapRegPrefix = `LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Class\{4D36E972-E325-11CE-BFC1-08002BE10318}`

func MustGetTapReg() []string {
	outArray := []string{}
	nameList := udwWindowsRegistry.MustGetDirectoryOrFileNameListOneLevel(
		tapRegPrefix)
	for _, name := range nameList {
		val, err := udwWindowsRegistry.GetStringByPath(
			tapRegPrefix + `\` + name + `\ComponentId`)
		if err != nil {
			if udwWindowsRegistry.IsErrorNotExist(err) || udwWindowsRegistry.IsErrorAccessDenied(err) {
				fmt.Println("[MustGetTapReg]", err)
				continue
			} else {
				panic(err)
			}
		}
		if val != "tap0901" {
			continue
		}
		guid, err := udwWindowsRegistry.GetStringByPath(
			`LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Class\{4D36E972-E325-11CE-BFC1-08002BE10318}\` + name + `\NetCfgInstanceId`)
		if err != nil {
			panic(err)
		}
		outArray = append(outArray, guid)
	}
	return outArray
}

const panelRegPrefix = `LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Network\{4D36E972-E325-11CE-BFC1-08002BE10318}`

func MustGetPanelReg() []PanelReg {
	outArray := []PanelReg{}
	nameList := udwWindowsRegistry.MustGetDirectoryOrFileNameListOneLevel(
		panelRegPrefix)
	for _, name := range nameList {
		val, err := udwWindowsRegistry.GetStringByPath(
			panelRegPrefix + `\` + name + `\Connection\Name`)
		if err != nil {
			if udwWindowsRegistry.IsErrorNotExist(err) || udwWindowsRegistry.IsErrorAccessDenied(err) {
				fmt.Println("[MustGetPanelReg]", err)
				continue
			} else {
				panic(err)
			}
		}
		outArray = append(outArray, PanelReg{
			Name: val,
			Guid: name,
		})
	}
	return outArray
}

func GetDeviceGuidAndActualName() (guid, actualName string) {
	tapRegList := MustGetTapReg()
	if len(tapRegList) == 0 {
		panic("can not found tap guid")
	}
	PanelRegList := MustGetPanelReg()
	for _, guid := range tapRegList {
		for _, panelReg := range PanelRegList {
			if panelReg.Guid == guid {
				return guid, panelReg.Name
			}
		}
	}
	panic("can not GetDeviceGuidAndActualName()")
}
