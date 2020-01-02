// +build !amd64,!s390x gccgo appengine

package cipherhw

func AESGCMSupport() bool {
	return false
}
