package udwSys

import (
	"errors"
)

func GetErrPlatformNotSupport() error {
	return errors.New("Platform Not Support")
}

const MaxMaxFileNum = 1048576

func MustSetCurrentMaxFileNum(limit uint64) {
	err := SetCurrentMaxFileNum(limit)
	if err != nil {
		panic(err)
	}
}

func MustSetCurrentMaxFileNumToMax() {
	MustSetCurrentMaxFileNum(MaxMaxFileNum)
}
