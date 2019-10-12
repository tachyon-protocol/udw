// +build windows

package udwTapTun

import (
	"fmt"
	"testing"
)

func ctlCode(code int) int {
	return (((0x00000022) << 16) | ((0) << 14) | ((code) << 2) | (0))
}

func TestTapCode(t *testing.T) {

	fmt.Println(ctlCode(2))
	fmt.Println(ctlCode(10))
	fmt.Println(ctlCode(6))
	fmt.Println(ctlCode(3))
	fmt.Println(ctlCode(5))
	fmt.Println("----------")
	fmt.Println(ctlCode(7))
}
