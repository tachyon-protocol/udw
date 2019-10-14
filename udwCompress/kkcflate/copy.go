package kkcflate

func forwardCopy(mem []byte, dst, src, n int) {
	if dst <= src {
		copy(mem[dst:dst+n], mem[src:src+n])
		return
	}
	for {
		if dst >= src+n {
			copy(mem[dst:dst+n], mem[src:src+n])
			return
		}

		k := dst - src
		copy(mem[dst:dst+k], mem[src:src+k])
		n -= k
		dst += k
	}
}
