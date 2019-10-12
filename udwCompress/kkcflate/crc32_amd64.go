//+build !noasm
//+build !appengine
//+build amd64

package kkcflate

import (
	"github.com/tachyon-protocol/udw/udwCompress/kkccpuid"
)

func crc32sse(a []byte) uint32

func crc32sseAll(a []byte, dst []uint32)

func matchLenSSE4(a, b []byte, max int) int

func histogram(b []byte, h []int32)

func init() {
	useSSE42 = kkccpuid.CPU.SSE42()
}
