//+build linux,android,amazon,!browser,!chromeExtension,!js

package udwPlatform

func GetCurrentUdwOs() string {
	return UdwOsAmazon
}
