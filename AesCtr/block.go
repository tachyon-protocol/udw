package AesCtr

func encryptBlockGo(xk []uint32, dst, src []byte) {
	var s0, s1, s2, s3, t0, t1, t2, t3 uint32

	s0 = uint32(src[0])<<24 | uint32(src[1])<<16 | uint32(src[2])<<8 | uint32(src[3])
	s1 = uint32(src[4])<<24 | uint32(src[5])<<16 | uint32(src[6])<<8 | uint32(src[7])
	s2 = uint32(src[8])<<24 | uint32(src[9])<<16 | uint32(src[10])<<8 | uint32(src[11])
	s3 = uint32(src[12])<<24 | uint32(src[13])<<16 | uint32(src[14])<<8 | uint32(src[15])

	s0 ^= xk[0]
	s1 ^= xk[1]
	s2 ^= xk[2]
	s3 ^= xk[3]

	nr := len(xk)/4 - 2
	k := 4
	for r := 0; r < nr; r++ {
		t0 = xk[k+0] ^ te0[uint8(s0>>24)] ^ te1[uint8(s1>>16)] ^ te2[uint8(s2>>8)] ^ te3[uint8(s3)]
		t1 = xk[k+1] ^ te0[uint8(s1>>24)] ^ te1[uint8(s2>>16)] ^ te2[uint8(s3>>8)] ^ te3[uint8(s0)]
		t2 = xk[k+2] ^ te0[uint8(s2>>24)] ^ te1[uint8(s3>>16)] ^ te2[uint8(s0>>8)] ^ te3[uint8(s1)]
		t3 = xk[k+3] ^ te0[uint8(s3>>24)] ^ te1[uint8(s0>>16)] ^ te2[uint8(s1>>8)] ^ te3[uint8(s2)]
		k += 4
		s0, s1, s2, s3 = t0, t1, t2, t3
	}

	s0 = uint32(sbox0[t0>>24])<<24 | uint32(sbox0[t1>>16&0xff])<<16 | uint32(sbox0[t2>>8&0xff])<<8 | uint32(sbox0[t3&0xff])
	s1 = uint32(sbox0[t1>>24])<<24 | uint32(sbox0[t2>>16&0xff])<<16 | uint32(sbox0[t3>>8&0xff])<<8 | uint32(sbox0[t0&0xff])
	s2 = uint32(sbox0[t2>>24])<<24 | uint32(sbox0[t3>>16&0xff])<<16 | uint32(sbox0[t0>>8&0xff])<<8 | uint32(sbox0[t1&0xff])
	s3 = uint32(sbox0[t3>>24])<<24 | uint32(sbox0[t0>>16&0xff])<<16 | uint32(sbox0[t1>>8&0xff])<<8 | uint32(sbox0[t2&0xff])

	s0 ^= xk[k+0]
	s1 ^= xk[k+1]
	s2 ^= xk[k+2]
	s3 ^= xk[k+3]

	dst[0], dst[1], dst[2], dst[3] = byte(s0>>24), byte(s0>>16), byte(s0>>8), byte(s0)
	dst[4], dst[5], dst[6], dst[7] = byte(s1>>24), byte(s1>>16), byte(s1>>8), byte(s1)
	dst[8], dst[9], dst[10], dst[11] = byte(s2>>24), byte(s2>>16), byte(s2>>8), byte(s2)
	dst[12], dst[13], dst[14], dst[15] = byte(s3>>24), byte(s3>>16), byte(s3>>8), byte(s3)
}

func decryptBlockGo(xk []uint32, dst, src []byte) {
	var s0, s1, s2, s3, t0, t1, t2, t3 uint32

	s0 = uint32(src[0])<<24 | uint32(src[1])<<16 | uint32(src[2])<<8 | uint32(src[3])
	s1 = uint32(src[4])<<24 | uint32(src[5])<<16 | uint32(src[6])<<8 | uint32(src[7])
	s2 = uint32(src[8])<<24 | uint32(src[9])<<16 | uint32(src[10])<<8 | uint32(src[11])
	s3 = uint32(src[12])<<24 | uint32(src[13])<<16 | uint32(src[14])<<8 | uint32(src[15])

	s0 ^= xk[0]
	s1 ^= xk[1]
	s2 ^= xk[2]
	s3 ^= xk[3]

	nr := len(xk)/4 - 2
	k := 4
	for r := 0; r < nr; r++ {
		t0 = xk[k+0] ^ td0[uint8(s0>>24)] ^ td1[uint8(s3>>16)] ^ td2[uint8(s2>>8)] ^ td3[uint8(s1)]
		t1 = xk[k+1] ^ td0[uint8(s1>>24)] ^ td1[uint8(s0>>16)] ^ td2[uint8(s3>>8)] ^ td3[uint8(s2)]
		t2 = xk[k+2] ^ td0[uint8(s2>>24)] ^ td1[uint8(s1>>16)] ^ td2[uint8(s0>>8)] ^ td3[uint8(s3)]
		t3 = xk[k+3] ^ td0[uint8(s3>>24)] ^ td1[uint8(s2>>16)] ^ td2[uint8(s1>>8)] ^ td3[uint8(s0)]
		k += 4
		s0, s1, s2, s3 = t0, t1, t2, t3
	}

	s0 = uint32(sbox1[t0>>24])<<24 | uint32(sbox1[t3>>16&0xff])<<16 | uint32(sbox1[t2>>8&0xff])<<8 | uint32(sbox1[t1&0xff])
	s1 = uint32(sbox1[t1>>24])<<24 | uint32(sbox1[t0>>16&0xff])<<16 | uint32(sbox1[t3>>8&0xff])<<8 | uint32(sbox1[t2&0xff])
	s2 = uint32(sbox1[t2>>24])<<24 | uint32(sbox1[t1>>16&0xff])<<16 | uint32(sbox1[t0>>8&0xff])<<8 | uint32(sbox1[t3&0xff])
	s3 = uint32(sbox1[t3>>24])<<24 | uint32(sbox1[t2>>16&0xff])<<16 | uint32(sbox1[t1>>8&0xff])<<8 | uint32(sbox1[t0&0xff])

	s0 ^= xk[k+0]
	s1 ^= xk[k+1]
	s2 ^= xk[k+2]
	s3 ^= xk[k+3]

	dst[0], dst[1], dst[2], dst[3] = byte(s0>>24), byte(s0>>16), byte(s0>>8), byte(s0)
	dst[4], dst[5], dst[6], dst[7] = byte(s1>>24), byte(s1>>16), byte(s1>>8), byte(s1)
	dst[8], dst[9], dst[10], dst[11] = byte(s2>>24), byte(s2>>16), byte(s2>>8), byte(s2)
	dst[12], dst[13], dst[14], dst[15] = byte(s3>>24), byte(s3>>16), byte(s3>>8), byte(s3)
}

func subw(w uint32) uint32 {
	return uint32(sbox0[w>>24])<<24 |
		uint32(sbox0[w>>16&0xff])<<16 |
		uint32(sbox0[w>>8&0xff])<<8 |
		uint32(sbox0[w&0xff])
}

func rotw(w uint32) uint32 { return w<<8 | w>>24 }

func expandKeyGo(key []byte, enc, dec []uint32) {

	var i int
	nk := len(key) / 4
	for i = 0; i < nk; i++ {
		enc[i] = uint32(key[4*i])<<24 | uint32(key[4*i+1])<<16 | uint32(key[4*i+2])<<8 | uint32(key[4*i+3])
	}
	for ; i < len(enc); i++ {
		t := enc[i-1]
		if i%nk == 0 {
			t = subw(rotw(t)) ^ (uint32(powx[i/nk-1]) << 24)
		} else if nk > 6 && i%nk == 4 {
			t = subw(t)
		}
		enc[i] = enc[i-nk] ^ t
	}

	if dec == nil {
		return
	}
	n := len(enc)
	for i := 0; i < n; i += 4 {
		ei := n - i - 4
		for j := 0; j < 4; j++ {
			x := enc[ei+j]
			if i > 0 && i+4 < n {
				x = td0[sbox0[x>>24]] ^ td1[sbox0[x>>16&0xff]] ^ td2[sbox0[x>>8&0xff]] ^ td3[sbox0[x&0xff]]
			}
			dec[i+j] = x
		}
	}
}
