// +build darwin,!ios

package udwOcMac

import (
	_ "github.com/tachyon-protocol/udw/udwOc/udwOcFoundation"
)

func ShowLocalNotification(context string) {
	udwShowLocalNotification(context)
}

func OpenUrlWithDefaultBrowser(url string) {
	udwOpenUrlWithDefaultBrowser(url)
}

func CheckIsRunningApp(bundleId string) bool {
	return udwCheckIsRunningApp(bundleId)
}

var applicationActiveFn func()

func SetApplicationDidBecomeActiveCallback(fn func()) {
	applicationActiveFn = fn
}

func applicationDidBecomeActiveCall() {
	if applicationActiveFn != nil {
		applicationActiveFn()
	}
}

func GetProcessIdByName(processName string) int {
	return udwGetProcessIdByName(processName)
}
