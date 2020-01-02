// +build s390x,!gccgo,!appengine

package cipherhw

func hasHWSupport() bool

var hwSupport = hasHWSupport()

func AESGCMSupport() bool {
	return hwSupport
}
