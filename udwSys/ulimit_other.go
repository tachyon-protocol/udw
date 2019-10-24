// +build !linux,!darwin,!windows

package udwSys

func SetCurrentMaxFileNum(limit uint64) (err error) {
	return GetErrPlatformNotSupport()
}
