// +build amd64,!gccgo,!appengine

package cipherhw

func hasAESNI() bool

func AESGCMSupport() bool {
	return hasAESNI()
}
