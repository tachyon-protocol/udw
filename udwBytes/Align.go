package udwBytes

func AlignToBigger(x uint64, align uint64) uint64 {
	mod := x % align
	if mod == 0 {
		return x
	}
	return x + (align - mod)
}

func AlignToSmaller(x uint64, align uint64) uint64 {
	mod := x % align
	if mod == 0 {
		return x
	}
	return x - mod
}
