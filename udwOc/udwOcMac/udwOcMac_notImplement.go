// +build !darwin ios

package udwOcMac

func ShowLocalNotification(context string) {
}

func OpenUrlWithDefaultBrowser(url string) {
}
func CheckIsRunningApp(bundleId string) bool {
	return false
}
func SetApplicationDidBecomeActiveCallback(fn func()) {
}

func GetProcessIdByName(processName string) int {
	return 0
}
