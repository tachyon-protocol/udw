//+build !amd64 noasm appengine

package kkcflate

func init() {
	useSSE42 = false
}

func crc32sse(a []byte) uint32 {
	panic("no assembler")
}

func crc32sseAll(a []byte, dst []uint32) {
	panic("no assembler")
}

func matchLenSSE4(a, b []byte, max int) int {
	panic("no assembler")
	return 0
}

func histogram(b []byte, h []int32) {
	h = h[:256]
	for _, t := range b {
		h[t]++
	}
}
