// +build linux,android,!browser,!chromeExtension,!js,!amazon

package udwPlatform

func GetCurrentUdwOs() string {
	return "android"
}
