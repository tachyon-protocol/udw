package udwPlatform

const UdwOsAndroid = "android"
const UdwOsIos = "ios"
const UdwOsMac = "mac"

const UdwOsLinux = "linux"
const UdwOsWindows = "windows"

const UdwOsJs = "js"

const UdwOsBrowser = "browser"

const UdwOsChromeExtension = "chromeExtension"
const UdwOsWeb = "web"
const UdwOsAmazon = `amazon`
const UdwOsWindowsuwp = `windowsuwp`

func UdwOsIsAndroidOrIos(udwOs string) bool {
	return udwOs == UdwOsAndroid || udwOs == UdwOsIos
}

func UdwOsIsWindowsOrMac(udwOs string) bool {
	return udwOs == UdwOsWindows || udwOs == UdwOsMac
}
