// +build windows

package udwWindowsRegistry

import "syscall"

const (
	_REG_OPTION_NON_VOLATILE = 0

	_REG_CREATED_NEW_KEY     = 1
	_REG_OPENED_EXISTING_KEY = 2

	_ERROR_NO_MORE_ITEMS syscall.Errno = 259
)
