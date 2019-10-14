package kkcflate

func emitLiteral(dst *tokens, lit []byte) {
	ol := int(dst.n)
	for i, v := range lit {
		dst.tokens[(i+ol)&maxStoreBlockSize] = token(v)
	}
	dst.n += uint16(len(lit))
}

func emitCopy(dst *tokens, offset, length int) {
	dst.tokens[dst.n] = matchToken(uint32(length-3), uint32(offset-minOffsetSize))
	dst.n++
}

type snappyEnc interface {
	Encode(dst *tokens, src []byte)
	Reset()
}

func newSnappy(level int) snappyEnc {
	switch level {
	case 1:
		return &snappyL1{}
	case 2:
		return &snappyL2{snappyGen: snappyGen{cur: maxStoreBlockSize, prev: make([]byte, 0, maxStoreBlockSize)}}
	case 3:
		return &snappyL3{snappyGen: snappyGen{cur: maxStoreBlockSize, prev: make([]byte, 0, maxStoreBlockSize)}}
	case 4:
		return &snappyL4{snappyL3{snappyGen: snappyGen{cur: maxStoreBlockSize, prev: make([]byte, 0, maxStoreBlockSize)}}}
	default:
		panic("invalid level specified")
	}
}

const (
	tableBits       = 14
	tableSize       = 1 << tableBits
	tableMask       = tableSize - 1
	tableShift      = 32 - tableBits
	baseMatchOffset = 1
	baseMatchLength = 3
	maxMatchOffset  = 1 << 15
)

func load32(b []byte, i int) uint32 {
	b = b[i : i+4 : len(b)]
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func load64(b []byte, i int) uint64 {
	b = b[i : i+8 : len(b)]
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func hash(u uint32) uint32 {
	return (u * 0x1e35a7bd) >> tableShift
}

type snappyL1 struct{}

func (e *snappyL1) Reset() {}

func (e *snappyL1) Encode(dst *tokens, src []byte) {
	const (
		inputMargin            = 16 - 1
		minNonLiteralBlockSize = 1 + 1 + inputMargin
	)

	if len(src) < minNonLiteralBlockSize {

		dst.n = uint16(len(src))
		return
	}

	var table [tableSize]uint16

	sLimit := len(src) - inputMargin

	nextEmit := 0

	s := 1
	nextHash := hash(load32(src, s))

	for {

		skip := 32

		nextS := s
		candidate := 0
		for {
			s = nextS
			bytesBetweenHashLookups := skip >> 5
			nextS = s + bytesBetweenHashLookups
			skip += bytesBetweenHashLookups
			if nextS > sLimit {
				goto emitRemainder
			}
			candidate = int(table[nextHash&tableMask])
			table[nextHash&tableMask] = uint16(s)
			nextHash = hash(load32(src, nextS))
			if s-candidate <= maxMatchOffset && load32(src, s) == load32(src, candidate) {
				break
			}
		}

		emitLiteral(dst, src[nextEmit:s])

		for {

			base := s

			s += 4
			s1 := base + maxMatchLength
			if s1 > len(src) {
				s1 = len(src)
			}
			a := src[s:s1]
			b := src[candidate+4:]
			b = b[:len(a)]
			l := len(a)
			for i := range a {
				if a[i] != b[i] {
					l = i
					break
				}
			}
			s += l

			dst.tokens[dst.n] = matchToken(uint32(s-base-baseMatchLength), uint32(base-candidate-baseMatchOffset))
			dst.n++
			nextEmit = s
			if s >= sLimit {
				goto emitRemainder
			}

			x := load64(src, s-1)
			prevHash := hash(uint32(x >> 0))
			table[prevHash&tableMask] = uint16(s - 1)
			currHash := hash(uint32(x >> 8))
			candidate = int(table[currHash&tableMask])
			table[currHash&tableMask] = uint16(s)
			if s-candidate > maxMatchOffset || uint32(x>>8) != load32(src, candidate) {
				nextHash = hash(uint32(x >> 16))
				s++
				break
			}
		}
	}

emitRemainder:
	if nextEmit < len(src) {
		emitLiteral(dst, src[nextEmit:])
	}
}

type tableEntry struct {
	val    uint32
	offset int32
}

func load3232(b []byte, i int32) uint32 {
	b = b[i : i+4 : len(b)]
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func load6432(b []byte, i int32) uint64 {
	b = b[i : i+8 : len(b)]
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

type snappyGen struct {
	prev []byte
	cur  int32
}

type snappyL2 struct {
	snappyGen
	table [tableSize]tableEntry
}

func (e *snappyL2) Encode(dst *tokens, src []byte) {
	const (
		inputMargin            = 8 - 1
		minNonLiteralBlockSize = 1 + 1 + inputMargin
	)

	if e.cur > 1<<30 {
		for i := range e.table {
			e.table[i] = tableEntry{}
		}
		e.cur = maxStoreBlockSize
	}

	if len(src) < minNonLiteralBlockSize {

		dst.n = uint16(len(src))
		e.cur += maxStoreBlockSize
		e.prev = e.prev[:0]
		return
	}

	sLimit := int32(len(src) - inputMargin)

	nextEmit := int32(0)
	s := int32(0)
	cv := load3232(src, s)
	nextHash := hash(cv)

	for {

		skip := int32(32)

		nextS := s
		var candidate tableEntry
		for {
			s = nextS
			bytesBetweenHashLookups := skip >> 5
			nextS = s + bytesBetweenHashLookups
			skip += bytesBetweenHashLookups
			if nextS > sLimit {
				goto emitRemainder
			}
			candidate = e.table[nextHash&tableMask]
			now := load3232(src, nextS)
			e.table[nextHash&tableMask] = tableEntry{offset: s + e.cur, val: cv}
			nextHash = hash(now)

			offset := s - (candidate.offset - e.cur)
			if offset > maxMatchOffset || cv != candidate.val {

				cv = now
				continue
			}
			break
		}

		emitLiteral(dst, src[nextEmit:s])

		for {

			s += 4
			t := candidate.offset - e.cur + 4
			l := e.matchlen(s, t, src)

			dst.tokens[dst.n] = matchToken(uint32(l+4-baseMatchLength), uint32(s-t-baseMatchOffset))
			dst.n++
			s += l
			nextEmit = s
			if s >= sLimit {
				t += l

				if int(t+4) < len(src) && t > 0 {
					cv := load3232(src, t)
					e.table[hash(cv)&tableMask] = tableEntry{offset: t + e.cur, val: cv}
				}
				goto emitRemainder
			}

			x := load6432(src, s-1)
			prevHash := hash(uint32(x))
			e.table[prevHash&tableMask] = tableEntry{offset: e.cur + s - 1, val: uint32(x)}
			x >>= 8
			currHash := hash(uint32(x))
			candidate = e.table[currHash&tableMask]
			e.table[currHash&tableMask] = tableEntry{offset: e.cur + s, val: uint32(x)}

			offset := s - (candidate.offset - e.cur)
			if offset > maxMatchOffset || uint32(x) != candidate.val {
				cv = uint32(x >> 8)
				nextHash = hash(cv)
				s++
				break
			}
		}
	}

emitRemainder:
	if int(nextEmit) < len(src) {
		emitLiteral(dst, src[nextEmit:])
	}
	e.cur += int32(len(src))
	e.prev = e.prev[:len(src)]
	copy(e.prev, src)
}

type tableEntryPrev struct {
	Cur  tableEntry
	Prev tableEntry
}

type snappyL3 struct {
	snappyGen
	table [tableSize]tableEntryPrev
}

func (e *snappyL3) Encode(dst *tokens, src []byte) {
	const (
		inputMargin            = 8 - 1
		minNonLiteralBlockSize = 1 + 1 + inputMargin
	)

	if e.cur > 1<<30 {
		for i := range e.table {
			e.table[i] = tableEntryPrev{}
		}
		e.snappyGen = snappyGen{cur: maxStoreBlockSize, prev: e.prev[:0]}
	}

	if len(src) < minNonLiteralBlockSize {

		dst.n = uint16(len(src))
		e.cur += maxStoreBlockSize
		e.prev = e.prev[:0]
		return
	}

	sLimit := int32(len(src) - inputMargin)

	nextEmit := int32(0)
	s := int32(0)
	cv := load3232(src, s)
	nextHash := hash(cv)

	for {

		skip := int32(32)

		nextS := s
		var candidate tableEntry
		for {
			s = nextS
			bytesBetweenHashLookups := skip >> 5
			nextS = s + bytesBetweenHashLookups
			skip += bytesBetweenHashLookups
			if nextS > sLimit {
				goto emitRemainder
			}
			candidates := e.table[nextHash&tableMask]
			now := load3232(src, nextS)
			e.table[nextHash&tableMask] = tableEntryPrev{Prev: candidates.Cur, Cur: tableEntry{offset: s + e.cur, val: cv}}
			nextHash = hash(now)

			candidate = candidates.Cur
			if cv == candidate.val {
				offset := s - (candidate.offset - e.cur)
				if offset <= maxMatchOffset {
					break
				}
			} else {

				candidate = candidates.Prev
				if cv == candidate.val {
					offset := s - (candidate.offset - e.cur)
					if offset <= maxMatchOffset {
						break
					}
				}
			}
			cv = now
		}

		emitLiteral(dst, src[nextEmit:s])

		for {

			s += 4
			t := candidate.offset - e.cur + 4
			l := e.matchlen(s, t, src)

			dst.tokens[dst.n] = matchToken(uint32(l+4-baseMatchLength), uint32(s-t-baseMatchOffset))
			dst.n++
			s += l
			nextEmit = s
			if s >= sLimit {
				t += l

				if int(t+4) < len(src) && t > 0 {
					cv := load3232(src, t)
					nextHash = hash(cv)
					e.table[nextHash&tableMask] = tableEntryPrev{
						Prev: e.table[nextHash&tableMask].Cur,
						Cur:  tableEntry{offset: e.cur + t, val: cv},
					}
				}
				goto emitRemainder
			}

			x := load6432(src, s-3)
			prevHash := hash(uint32(x))
			e.table[prevHash&tableMask] = tableEntryPrev{
				Prev: e.table[prevHash&tableMask].Cur,
				Cur:  tableEntry{offset: e.cur + s - 3, val: uint32(x)},
			}
			x >>= 8
			prevHash = hash(uint32(x))

			e.table[prevHash&tableMask] = tableEntryPrev{
				Prev: e.table[prevHash&tableMask].Cur,
				Cur:  tableEntry{offset: e.cur + s - 2, val: uint32(x)},
			}
			x >>= 8
			prevHash = hash(uint32(x))

			e.table[prevHash&tableMask] = tableEntryPrev{
				Prev: e.table[prevHash&tableMask].Cur,
				Cur:  tableEntry{offset: e.cur + s - 1, val: uint32(x)},
			}
			x >>= 8
			currHash := hash(uint32(x))
			candidates := e.table[currHash&tableMask]
			cv = uint32(x)
			e.table[currHash&tableMask] = tableEntryPrev{
				Prev: candidates.Cur,
				Cur:  tableEntry{offset: s + e.cur, val: cv},
			}

			candidate = candidates.Cur
			if cv == candidate.val {
				offset := s - (candidate.offset - e.cur)
				if offset <= maxMatchOffset {
					continue
				}
			} else {

				candidate = candidates.Prev
				if cv == candidate.val {
					offset := s - (candidate.offset - e.cur)
					if offset <= maxMatchOffset {
						continue
					}
				}
			}
			cv = uint32(x >> 8)
			nextHash = hash(cv)
			s++
			break
		}
	}

emitRemainder:
	if int(nextEmit) < len(src) {
		emitLiteral(dst, src[nextEmit:])
	}
	e.cur += int32(len(src))
	e.prev = e.prev[:len(src)]
	copy(e.prev, src)
}

type snappyL4 struct {
	snappyL3
}

func (e *snappyL4) Encode(dst *tokens, src []byte) {
	const (
		inputMargin            = 8 - 3
		minNonLiteralBlockSize = 1 + 1 + inputMargin
		matchLenGood           = 12
	)

	if e.cur > 1<<30 {
		for i := range e.table {
			e.table[i] = tableEntryPrev{}
		}
		e.snappyGen = snappyGen{cur: maxStoreBlockSize, prev: e.prev[:0]}
	}

	if len(src) < minNonLiteralBlockSize {

		dst.n = uint16(len(src))
		e.cur += maxStoreBlockSize
		e.prev = e.prev[:0]
		return
	}

	sLimit := int32(len(src) - inputMargin)

	nextEmit := int32(0)
	s := int32(0)
	cv := load3232(src, s)
	nextHash := hash(cv)

	for {

		skip := int32(32)

		nextS := s
		var candidate tableEntry
		var candidateAlt tableEntry
		for {
			s = nextS
			bytesBetweenHashLookups := skip >> 5
			nextS = s + bytesBetweenHashLookups
			skip += bytesBetweenHashLookups
			if nextS > sLimit {
				goto emitRemainder
			}
			candidates := e.table[nextHash&tableMask]
			now := load3232(src, nextS)
			e.table[nextHash&tableMask] = tableEntryPrev{Prev: candidates.Cur, Cur: tableEntry{offset: s + e.cur, val: cv}}
			nextHash = hash(now)

			candidate = candidates.Cur
			if cv == candidate.val {
				offset := s - (candidate.offset - e.cur)
				if offset < maxMatchOffset {
					offset = s - (candidates.Prev.offset - e.cur)
					if cv == candidates.Prev.val && offset < maxMatchOffset {
						candidateAlt = candidates.Prev
					}
					break
				}
			} else {

				candidate = candidates.Prev
				if cv == candidate.val {
					offset := s - (candidate.offset - e.cur)
					if offset < maxMatchOffset {
						break
					}
				}
			}
			cv = now
		}

		emitLiteral(dst, src[nextEmit:s])

		for {

			s += 4
			t := candidate.offset - e.cur + 4
			l := e.matchlen(s, t, src)

			if l < matchLenGood-4 && candidateAlt.offset != 0 {
				t2 := candidateAlt.offset - e.cur + 4
				l2 := e.matchlen(s, t2, src)
				if l2 > l {
					l = l2
					t = t2
				}
			}

			dst.tokens[dst.n] = matchToken(uint32(l+4-baseMatchLength), uint32(s-t-baseMatchOffset))
			dst.n++
			s += l
			nextEmit = s
			if s >= sLimit {
				t += l

				if int(t+4) < len(src) && t > 0 {
					cv := load3232(src, t)
					nextHash = hash(cv)
					e.table[nextHash&tableMask] = tableEntryPrev{
						Prev: e.table[nextHash&tableMask].Cur,
						Cur:  tableEntry{offset: e.cur + t, val: cv},
					}
				}
				goto emitRemainder
			}

			x := load6432(src, s-3)
			prevHash := hash(uint32(x))
			e.table[prevHash&tableMask] = tableEntryPrev{
				Prev: e.table[prevHash&tableMask].Cur,
				Cur:  tableEntry{offset: e.cur + s - 3, val: uint32(x)},
			}
			x >>= 8
			prevHash = hash(uint32(x))

			e.table[prevHash&tableMask] = tableEntryPrev{
				Prev: e.table[prevHash&tableMask].Cur,
				Cur:  tableEntry{offset: e.cur + s - 2, val: uint32(x)},
			}
			x >>= 8
			prevHash = hash(uint32(x))

			e.table[prevHash&tableMask] = tableEntryPrev{
				Prev: e.table[prevHash&tableMask].Cur,
				Cur:  tableEntry{offset: e.cur + s - 1, val: uint32(x)},
			}
			x >>= 8
			currHash := hash(uint32(x))
			candidates := e.table[currHash&tableMask]
			cv = uint32(x)
			e.table[currHash&tableMask] = tableEntryPrev{
				Prev: candidates.Cur,
				Cur:  tableEntry{offset: s + e.cur, val: cv},
			}

			candidate = candidates.Cur
			candidateAlt = tableEntry{}
			if cv == candidate.val {
				offset := s - (candidate.offset - e.cur)
				if offset <= maxMatchOffset {
					offset = s - (candidates.Prev.offset - e.cur)
					if cv == candidates.Prev.val && offset <= maxMatchOffset {
						candidateAlt = candidates.Prev
					}
					continue
				}
			} else {

				candidate = candidates.Prev
				if cv == candidate.val {
					offset := s - (candidate.offset - e.cur)
					if offset <= maxMatchOffset {
						continue
					}
				}
			}
			cv = uint32(x >> 8)
			nextHash = hash(cv)
			s++
			break
		}
	}

emitRemainder:
	if int(nextEmit) < len(src) {
		emitLiteral(dst, src[nextEmit:])
	}
	e.cur += int32(len(src))
	e.prev = e.prev[:len(src)]
	copy(e.prev, src)
}

func (e *snappyGen) matchlen(s, t int32, src []byte) int32 {
	s1 := int(s) + maxMatchLength - 4
	if s1 > len(src) {
		s1 = len(src)
	}

	if t >= 0 {
		b := src[t:]
		a := src[s:s1]
		b = b[:len(a)]

		for i := range a {
			if a[i] != b[i] {
				return int32(i)
			}
		}
		return int32(len(a))
	}

	tp := int32(len(e.prev)) + t
	if tp < 0 {
		return 0
	}

	a := src[s:s1]
	b := e.prev[tp:]
	if len(b) > len(a) {
		b = b[:len(a)]
	}
	a = a[:len(b)]
	for i := range b {
		if a[i] != b[i] {
			return int32(i)
		}
	}

	n := int32(len(b))
	if int(s+n) == s1 {
		return n
	}

	a = src[s+n : s1]
	b = src[:len(a)]
	for i := range a {
		if a[i] != b[i] {
			return int32(i) + n
		}
	}
	return int32(len(a)) + n
}

func (e *snappyGen) Reset() {
	e.prev = e.prev[:0]
	e.cur += maxMatchOffset
}
