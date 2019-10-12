package kkcflate

import (
	"fmt"
	"io"
	"math"
)

const (
	NoCompression      = 0
	BestSpeed          = 1
	BestCompression    = 9
	DefaultCompression = -1

	HuffmanOnly         = -2
	ConstantCompression = HuffmanOnly

	logWindowSize    = 15
	windowSize       = 1 << logWindowSize
	windowMask       = windowSize - 1
	logMaxOffsetSize = 15
	minMatchLength   = 4
	maxMatchLength   = 258
	minOffsetSize    = 1

	maxFlateBlockTokens = 1 << 14
	maxStoreBlockSize   = 65535
	hashBits            = 17
	hashSize            = 1 << hashBits
	hashMask            = (1 << hashBits) - 1
	hashShift           = (hashBits + minMatchLength - 1) / minMatchLength
	maxHashOffset       = 1 << 24

	skipNever = math.MaxInt32
)

var useSSE42 bool

type compressionLevel struct {
	good, lazy, nice, chain, fastSkipHashing, level int
}

var levels = []compressionLevel{
	{},

	{0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 2},
	{0, 0, 0, 0, 0, 3},
	{0, 0, 0, 0, 0, 4},

	{6, 0, 12, 8, 12, 5},
	{8, 0, 24, 16, 16, 6},

	{8, 8, 24, 16, skipNever, 7},
	{10, 16, 24, 64, skipNever, 8},
	{32, 258, 258, 4096, skipNever, 9},
}

type compressor struct {
	compressionLevel

	w          *huffmanBitWriter
	bulkHasher func([]byte, []uint32)

	fill func(*compressor, []byte) int
	step func(*compressor)
	sync bool

	chainHead  int
	hashHead   [hashSize]uint32
	hashPrev   [windowSize]uint32
	hashOffset int

	index         int
	window        []byte
	windowEnd     int
	blockStart    int
	byteAvailable bool

	tokens tokens

	length         int
	offset         int
	hash           uint32
	maxInsertIndex int
	err            error
	ii             uint16

	snap      snappyEnc
	hashMatch [maxMatchLength + minMatchLength]uint32
}

func (d *compressor) fillDeflate(b []byte) int {
	if d.index >= 2*windowSize-(minMatchLength+maxMatchLength) {

		copy(d.window[:], d.window[windowSize:2*windowSize])
		d.index -= windowSize
		d.windowEnd -= windowSize
		if d.blockStart >= windowSize {
			d.blockStart -= windowSize
		} else {
			d.blockStart = math.MaxInt32
		}
		d.hashOffset += windowSize
		if d.hashOffset > maxHashOffset {
			delta := d.hashOffset - 1
			d.hashOffset -= delta
			d.chainHead -= delta
			for i, v := range d.hashPrev {
				if int(v) > delta {
					d.hashPrev[i] = uint32(int(v) - delta)
				} else {
					d.hashPrev[i] = 0
				}
			}
			for i, v := range d.hashHead {
				if int(v) > delta {
					d.hashHead[i] = uint32(int(v) - delta)
				} else {
					d.hashHead[i] = 0
				}
			}
		}
	}
	n := copy(d.window[d.windowEnd:], b)
	d.windowEnd += n
	return n
}

func (d *compressor) writeBlock(tok tokens, index int, eof bool) error {
	if index > 0 || eof {
		var window []byte
		if d.blockStart <= index {
			window = d.window[d.blockStart:index]
		}
		d.blockStart = index
		d.w.writeBlock(tok.tokens[:tok.n], eof, window)
		return d.w.err
	}
	return nil
}

func (d *compressor) writeBlockSkip(tok tokens, index int, eof bool) error {
	if index > 0 || eof {
		if d.blockStart <= index {
			window := d.window[d.blockStart:index]

			if int(tok.n) > len(window)-int(tok.n>>6) {
				d.w.writeBlockHuff(eof, window)
			} else {

				d.w.writeBlockDynamic(tok.tokens[:tok.n], eof, window)
			}
		} else {
			d.w.writeBlock(tok.tokens[:tok.n], eof, nil)
		}
		d.blockStart = index
		return d.w.err
	}
	return nil
}

func (d *compressor) fillWindow(b []byte) {

	switch d.compressionLevel.level {
	case 0, 1, 2:
		return
	}

	if len(b) > windowSize {
		b = b[len(b)-windowSize:]
	}

	n := copy(d.window[d.windowEnd:], b)

	loops := (n + 256 - minMatchLength) / 256
	for j := 0; j < loops; j++ {
		startindex := j * 256
		end := startindex + 256 + minMatchLength - 1
		if end > n {
			end = n
		}
		tocheck := d.window[startindex:end]
		dstSize := len(tocheck) - minMatchLength + 1

		if dstSize <= 0 {
			continue
		}

		dst := d.hashMatch[:dstSize]
		d.bulkHasher(tocheck, dst)
		var newH uint32
		for i, val := range dst {
			di := i + startindex
			newH = val & hashMask

			d.hashPrev[di&windowMask] = d.hashHead[newH]

			d.hashHead[newH] = uint32(di + d.hashOffset)
		}
		d.hash = newH
	}

	d.windowEnd += n
	d.index = n
}

func (d *compressor) findMatch(pos int, prevHead int, prevLength int, lookahead int) (length, offset int, ok bool) {
	minMatchLook := maxMatchLength
	if lookahead < minMatchLook {
		minMatchLook = lookahead
	}

	win := d.window[0 : pos+minMatchLook]

	nice := len(win) - pos
	if d.nice < nice {
		nice = d.nice
	}

	tries := d.chain
	length = prevLength
	if length >= d.good {
		tries >>= 2
	}

	wEnd := win[pos+length]
	wPos := win[pos:]
	minIndex := pos - windowSize

	for i := prevHead; tries > 0; tries-- {
		if wEnd == win[i+length] {
			n := matchLen(win[i:], wPos, minMatchLook)

			if n > length && (n > minMatchLength || pos-i <= 4096) {
				length = n
				offset = pos - i
				ok = true
				if n >= nice {

					break
				}
				wEnd = win[pos+n]
			}
		}
		if i == minIndex {

			break
		}
		i = int(d.hashPrev[i&windowMask]) - d.hashOffset
		if i < minIndex || i < 0 {
			break
		}
	}
	return
}

func (d *compressor) findMatchSSE(pos int, prevHead int, prevLength int, lookahead int) (length, offset int, ok bool) {
	minMatchLook := maxMatchLength
	if lookahead < minMatchLook {
		minMatchLook = lookahead
	}

	win := d.window[0 : pos+minMatchLook]

	nice := len(win) - pos
	if d.nice < nice {
		nice = d.nice
	}

	tries := d.chain
	length = prevLength
	if length >= d.good {
		tries >>= 2
	}

	wEnd := win[pos+length]
	wPos := win[pos:]
	minIndex := pos - windowSize

	for i := prevHead; tries > 0; tries-- {
		if wEnd == win[i+length] {
			n := matchLenSSE4(win[i:], wPos, minMatchLook)

			if n > length && (n > minMatchLength || pos-i <= 4096) {
				length = n
				offset = pos - i
				ok = true
				if n >= nice {

					break
				}
				wEnd = win[pos+n]
			}
		}
		if i == minIndex {

			break
		}
		i = int(d.hashPrev[i&windowMask]) - d.hashOffset
		if i < minIndex || i < 0 {
			break
		}
	}
	return
}

func (d *compressor) writeStoredBlock(buf []byte) error {
	if d.w.writeStoredHeader(len(buf), false); d.w.err != nil {
		return d.w.err
	}
	d.w.writeBytes(buf)
	return d.w.err
}

const hashmul = 0x1e35a7bd

func hash4(b []byte) uint32 {
	return ((uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24) * hashmul) >> (32 - hashBits)
}

func bulkHash4(b []byte, dst []uint32) {
	if len(b) < minMatchLength {
		return
	}
	hb := uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
	dst[0] = (hb * hashmul) >> (32 - hashBits)
	end := len(b) - minMatchLength + 1
	for i := 1; i < end; i++ {
		hb = (hb << 8) | uint32(b[i+3])
		dst[i] = (hb * hashmul) >> (32 - hashBits)
	}
}

func matchLen(a, b []byte, max int) int {
	a = a[:max]
	b = b[:len(a)]
	for i, av := range a {
		if b[i] != av {
			return i
		}
	}
	return max
}

func (d *compressor) initDeflate() {
	d.window = make([]byte, 2*windowSize)
	d.hashOffset = 1
	d.length = minMatchLength - 1
	d.offset = 0
	d.byteAvailable = false
	d.index = 0
	d.hash = 0
	d.chainHead = -1
	d.bulkHasher = bulkHash4
	if useSSE42 {
		d.bulkHasher = crc32sseAll
	}
}

func (d *compressor) deflate() {

	const sanity = false

	if d.windowEnd-d.index < minMatchLength+maxMatchLength && !d.sync {
		return
	}

	d.maxInsertIndex = d.windowEnd - (minMatchLength - 1)
	if d.index < d.maxInsertIndex {
		d.hash = hash4(d.window[d.index : d.index+minMatchLength])
	}

	for {
		if sanity && d.index > d.windowEnd {
			panic("index > windowEnd")
		}
		lookahead := d.windowEnd - d.index
		if lookahead < minMatchLength+maxMatchLength {
			if !d.sync {
				return
			}
			if sanity && d.index > d.windowEnd {
				panic("index > windowEnd")
			}
			if lookahead == 0 {
				if d.tokens.n > 0 {
					if d.err = d.writeBlockSkip(d.tokens, d.index, false); d.err != nil {
						return
					}
					d.tokens.n = 0
				}
				return
			}
		}
		if d.index < d.maxInsertIndex {

			d.hash = hash4(d.window[d.index : d.index+minMatchLength])
			ch := d.hashHead[d.hash&hashMask]
			d.chainHead = int(ch)
			d.hashPrev[d.index&windowMask] = ch
			d.hashHead[d.hash&hashMask] = uint32(d.index + d.hashOffset)
		}
		d.length = minMatchLength - 1
		d.offset = 0
		minIndex := d.index - windowSize
		if minIndex < 0 {
			minIndex = 0
		}

		if d.chainHead-d.hashOffset >= minIndex && lookahead > minMatchLength-1 {
			if newLength, newOffset, ok := d.findMatch(d.index, d.chainHead-d.hashOffset, minMatchLength-1, lookahead); ok {
				d.length = newLength
				d.offset = newOffset
			}
		}
		if d.length >= minMatchLength {
			d.ii = 0

			d.tokens.tokens[d.tokens.n] = matchToken(uint32(d.length-3), uint32(d.offset-minOffsetSize))
			d.tokens.n++

			if d.length <= d.fastSkipHashing {
				var newIndex int
				newIndex = d.index + d.length

				end := newIndex
				if end > d.maxInsertIndex {
					end = d.maxInsertIndex
				}
				end += minMatchLength - 1
				startindex := d.index + 1
				if startindex > d.maxInsertIndex {
					startindex = d.maxInsertIndex
				}
				tocheck := d.window[startindex:end]
				dstSize := len(tocheck) - minMatchLength + 1
				if dstSize > 0 {
					dst := d.hashMatch[:dstSize]
					bulkHash4(tocheck, dst)
					var newH uint32
					for i, val := range dst {
						di := i + startindex
						newH = val & hashMask

						d.hashPrev[di&windowMask] = d.hashHead[newH]

						d.hashHead[newH] = uint32(di + d.hashOffset)
					}
					d.hash = newH
				}
				d.index = newIndex
			} else {

				d.index += d.length
				if d.index < d.maxInsertIndex {
					d.hash = hash4(d.window[d.index : d.index+minMatchLength])
				}
			}
			if d.tokens.n == maxFlateBlockTokens {

				if d.err = d.writeBlockSkip(d.tokens, d.index, false); d.err != nil {
					return
				}
				d.tokens.n = 0
			}
		} else {
			d.ii++
			end := d.index + int(d.ii>>uint(d.fastSkipHashing)) + 1
			if end > d.windowEnd {
				end = d.windowEnd
			}
			for i := d.index; i < end; i++ {
				d.tokens.tokens[d.tokens.n] = literalToken(uint32(d.window[i]))
				d.tokens.n++
				if d.tokens.n == maxFlateBlockTokens {
					if d.err = d.writeBlockSkip(d.tokens, i+1, false); d.err != nil {
						return
					}
					d.tokens.n = 0
				}
			}
			d.index = end
		}
	}
}

func (d *compressor) deflateLazy() {

	const sanity = false

	if d.windowEnd-d.index < minMatchLength+maxMatchLength && !d.sync {
		return
	}

	d.maxInsertIndex = d.windowEnd - (minMatchLength - 1)
	if d.index < d.maxInsertIndex {
		d.hash = hash4(d.window[d.index : d.index+minMatchLength])
	}

	for {
		if sanity && d.index > d.windowEnd {
			panic("index > windowEnd")
		}
		lookahead := d.windowEnd - d.index
		if lookahead < minMatchLength+maxMatchLength {
			if !d.sync {
				return
			}
			if sanity && d.index > d.windowEnd {
				panic("index > windowEnd")
			}
			if lookahead == 0 {

				if d.byteAvailable {

					d.tokens.tokens[d.tokens.n] = literalToken(uint32(d.window[d.index-1]))
					d.tokens.n++
					d.byteAvailable = false
				}
				if d.tokens.n > 0 {
					if d.err = d.writeBlock(d.tokens, d.index, false); d.err != nil {
						return
					}
					d.tokens.n = 0
				}
				return
			}
		}
		if d.index < d.maxInsertIndex {

			d.hash = hash4(d.window[d.index : d.index+minMatchLength])
			ch := d.hashHead[d.hash&hashMask]
			d.chainHead = int(ch)
			d.hashPrev[d.index&windowMask] = ch
			d.hashHead[d.hash&hashMask] = uint32(d.index + d.hashOffset)
		}
		prevLength := d.length
		prevOffset := d.offset
		d.length = minMatchLength - 1
		d.offset = 0
		minIndex := d.index - windowSize
		if minIndex < 0 {
			minIndex = 0
		}

		if d.chainHead-d.hashOffset >= minIndex && lookahead > prevLength && prevLength < d.lazy {
			if newLength, newOffset, ok := d.findMatch(d.index, d.chainHead-d.hashOffset, minMatchLength-1, lookahead); ok {
				d.length = newLength
				d.offset = newOffset
			}
		}
		if prevLength >= minMatchLength && d.length <= prevLength {

			d.tokens.tokens[d.tokens.n] = matchToken(uint32(prevLength-3), uint32(prevOffset-minOffsetSize))
			d.tokens.n++

			var newIndex int
			newIndex = d.index + prevLength - 1

			end := newIndex
			if end > d.maxInsertIndex {
				end = d.maxInsertIndex
			}
			end += minMatchLength - 1
			startindex := d.index + 1
			if startindex > d.maxInsertIndex {
				startindex = d.maxInsertIndex
			}
			tocheck := d.window[startindex:end]
			dstSize := len(tocheck) - minMatchLength + 1
			if dstSize > 0 {
				dst := d.hashMatch[:dstSize]
				bulkHash4(tocheck, dst)
				var newH uint32
				for i, val := range dst {
					di := i + startindex
					newH = val & hashMask

					d.hashPrev[di&windowMask] = d.hashHead[newH]

					d.hashHead[newH] = uint32(di + d.hashOffset)
				}
				d.hash = newH
			}

			d.index = newIndex
			d.byteAvailable = false
			d.length = minMatchLength - 1
			if d.tokens.n == maxFlateBlockTokens {

				if d.err = d.writeBlock(d.tokens, d.index, false); d.err != nil {
					return
				}
				d.tokens.n = 0
			}
		} else {

			if d.length >= minMatchLength {
				d.ii = 0
			}

			if d.byteAvailable {
				d.ii++
				d.tokens.tokens[d.tokens.n] = literalToken(uint32(d.window[d.index-1]))
				d.tokens.n++
				if d.tokens.n == maxFlateBlockTokens {
					if d.err = d.writeBlock(d.tokens, d.index, false); d.err != nil {
						return
					}
					d.tokens.n = 0
				}
				d.index++

				if d.ii > 31 {
					n := int(d.ii >> 5)
					for j := 0; j < n; j++ {
						if d.index >= d.windowEnd-1 {
							break
						}

						d.tokens.tokens[d.tokens.n] = literalToken(uint32(d.window[d.index-1]))
						d.tokens.n++
						if d.tokens.n == maxFlateBlockTokens {
							if d.err = d.writeBlock(d.tokens, d.index, false); d.err != nil {
								return
							}
							d.tokens.n = 0
						}
						d.index++
					}

					d.tokens.tokens[d.tokens.n] = literalToken(uint32(d.window[d.index-1]))
					d.tokens.n++
					d.byteAvailable = false

					if d.tokens.n == maxFlateBlockTokens {
						if d.err = d.writeBlock(d.tokens, d.index, false); d.err != nil {
							return
						}
						d.tokens.n = 0
					}
				}
			} else {
				d.index++
				d.byteAvailable = true
			}
		}
	}
}

func (d *compressor) deflateSSE() {

	const sanity = false

	if d.windowEnd-d.index < minMatchLength+maxMatchLength && !d.sync {
		return
	}

	d.maxInsertIndex = d.windowEnd - (minMatchLength - 1)
	if d.index < d.maxInsertIndex {
		d.hash = crc32sse(d.window[d.index:d.index+minMatchLength]) & hashMask
	}

	for {
		if sanity && d.index > d.windowEnd {
			panic("index > windowEnd")
		}
		lookahead := d.windowEnd - d.index
		if lookahead < minMatchLength+maxMatchLength {
			if !d.sync {
				return
			}
			if sanity && d.index > d.windowEnd {
				panic("index > windowEnd")
			}
			if lookahead == 0 {
				if d.tokens.n > 0 {
					if d.err = d.writeBlockSkip(d.tokens, d.index, false); d.err != nil {
						return
					}
					d.tokens.n = 0
				}
				return
			}
		}
		if d.index < d.maxInsertIndex {

			d.hash = crc32sse(d.window[d.index:d.index+minMatchLength]) & hashMask
			ch := d.hashHead[d.hash]
			d.chainHead = int(ch)
			d.hashPrev[d.index&windowMask] = ch
			d.hashHead[d.hash] = uint32(d.index + d.hashOffset)
		}
		d.length = minMatchLength - 1
		d.offset = 0
		minIndex := d.index - windowSize
		if minIndex < 0 {
			minIndex = 0
		}

		if d.chainHead-d.hashOffset >= minIndex && lookahead > minMatchLength-1 {
			if newLength, newOffset, ok := d.findMatchSSE(d.index, d.chainHead-d.hashOffset, minMatchLength-1, lookahead); ok {
				d.length = newLength
				d.offset = newOffset
			}
		}
		if d.length >= minMatchLength {
			d.ii = 0

			d.tokens.tokens[d.tokens.n] = matchToken(uint32(d.length-3), uint32(d.offset-minOffsetSize))
			d.tokens.n++

			if d.length <= d.fastSkipHashing {
				var newIndex int
				newIndex = d.index + d.length

				end := newIndex
				if end > d.maxInsertIndex {
					end = d.maxInsertIndex
				}
				end += minMatchLength - 1
				startindex := d.index + 1
				if startindex > d.maxInsertIndex {
					startindex = d.maxInsertIndex
				}
				tocheck := d.window[startindex:end]
				dstSize := len(tocheck) - minMatchLength + 1
				if dstSize > 0 {
					dst := d.hashMatch[:dstSize]

					crc32sseAll(tocheck, dst)
					var newH uint32
					for i, val := range dst {
						di := i + startindex
						newH = val & hashMask

						d.hashPrev[di&windowMask] = d.hashHead[newH]

						d.hashHead[newH] = uint32(di + d.hashOffset)
					}
					d.hash = newH
				}
				d.index = newIndex
			} else {

				d.index += d.length
				if d.index < d.maxInsertIndex {
					d.hash = crc32sse(d.window[d.index:d.index+minMatchLength]) & hashMask
				}
			}
			if d.tokens.n == maxFlateBlockTokens {

				if d.err = d.writeBlockSkip(d.tokens, d.index, false); d.err != nil {
					return
				}
				d.tokens.n = 0
			}
		} else {
			d.ii++
			end := d.index + int(d.ii>>5) + 1
			if end > d.windowEnd {
				end = d.windowEnd
			}
			for i := d.index; i < end; i++ {
				d.tokens.tokens[d.tokens.n] = literalToken(uint32(d.window[i]))
				d.tokens.n++
				if d.tokens.n == maxFlateBlockTokens {
					if d.err = d.writeBlockSkip(d.tokens, i+1, false); d.err != nil {
						return
					}
					d.tokens.n = 0
				}
			}
			d.index = end
		}
	}
}

func (d *compressor) deflateLazySSE() {

	const sanity = false

	if d.windowEnd-d.index < minMatchLength+maxMatchLength && !d.sync {
		return
	}

	d.maxInsertIndex = d.windowEnd - (minMatchLength - 1)
	if d.index < d.maxInsertIndex {
		d.hash = crc32sse(d.window[d.index:d.index+minMatchLength]) & hashMask
	}

	for {
		if sanity && d.index > d.windowEnd {
			panic("index > windowEnd")
		}
		lookahead := d.windowEnd - d.index
		if lookahead < minMatchLength+maxMatchLength {
			if !d.sync {
				return
			}
			if sanity && d.index > d.windowEnd {
				panic("index > windowEnd")
			}
			if lookahead == 0 {

				if d.byteAvailable {

					d.tokens.tokens[d.tokens.n] = literalToken(uint32(d.window[d.index-1]))
					d.tokens.n++
					d.byteAvailable = false
				}
				if d.tokens.n > 0 {
					if d.err = d.writeBlock(d.tokens, d.index, false); d.err != nil {
						return
					}
					d.tokens.n = 0
				}
				return
			}
		}
		if d.index < d.maxInsertIndex {

			d.hash = crc32sse(d.window[d.index:d.index+minMatchLength]) & hashMask
			ch := d.hashHead[d.hash]
			d.chainHead = int(ch)
			d.hashPrev[d.index&windowMask] = ch
			d.hashHead[d.hash] = uint32(d.index + d.hashOffset)
		}
		prevLength := d.length
		prevOffset := d.offset
		d.length = minMatchLength - 1
		d.offset = 0
		minIndex := d.index - windowSize
		if minIndex < 0 {
			minIndex = 0
		}

		if d.chainHead-d.hashOffset >= minIndex && lookahead > prevLength && prevLength < d.lazy {
			if newLength, newOffset, ok := d.findMatchSSE(d.index, d.chainHead-d.hashOffset, minMatchLength-1, lookahead); ok {
				d.length = newLength
				d.offset = newOffset
			}
		}
		if prevLength >= minMatchLength && d.length <= prevLength {

			d.tokens.tokens[d.tokens.n] = matchToken(uint32(prevLength-3), uint32(prevOffset-minOffsetSize))
			d.tokens.n++

			var newIndex int
			newIndex = d.index + prevLength - 1

			end := newIndex
			if end > d.maxInsertIndex {
				end = d.maxInsertIndex
			}
			end += minMatchLength - 1
			startindex := d.index + 1
			if startindex > d.maxInsertIndex {
				startindex = d.maxInsertIndex
			}
			tocheck := d.window[startindex:end]
			dstSize := len(tocheck) - minMatchLength + 1
			if dstSize > 0 {
				dst := d.hashMatch[:dstSize]
				crc32sseAll(tocheck, dst)
				var newH uint32
				for i, val := range dst {
					di := i + startindex
					newH = val & hashMask

					d.hashPrev[di&windowMask] = d.hashHead[newH]

					d.hashHead[newH] = uint32(di + d.hashOffset)
				}
				d.hash = newH
			}

			d.index = newIndex
			d.byteAvailable = false
			d.length = minMatchLength - 1
			if d.tokens.n == maxFlateBlockTokens {

				if d.err = d.writeBlock(d.tokens, d.index, false); d.err != nil {
					return
				}
				d.tokens.n = 0
			}
		} else {

			if d.length >= minMatchLength {
				d.ii = 0
			}

			if d.byteAvailable {
				d.ii++
				d.tokens.tokens[d.tokens.n] = literalToken(uint32(d.window[d.index-1]))
				d.tokens.n++
				if d.tokens.n == maxFlateBlockTokens {
					if d.err = d.writeBlock(d.tokens, d.index, false); d.err != nil {
						return
					}
					d.tokens.n = 0
				}
				d.index++

				if d.ii > 31 {
					n := int(d.ii >> 6)
					for j := 0; j < n; j++ {
						if d.index >= d.windowEnd-1 {
							break
						}

						d.tokens.tokens[d.tokens.n] = literalToken(uint32(d.window[d.index-1]))
						d.tokens.n++
						if d.tokens.n == maxFlateBlockTokens {
							if d.err = d.writeBlock(d.tokens, d.index, false); d.err != nil {
								return
							}
							d.tokens.n = 0
						}
						d.index++
					}

					d.tokens.tokens[d.tokens.n] = literalToken(uint32(d.window[d.index-1]))
					d.tokens.n++
					d.byteAvailable = false

					if d.tokens.n == maxFlateBlockTokens {
						if d.err = d.writeBlock(d.tokens, d.index, false); d.err != nil {
							return
						}
						d.tokens.n = 0
					}
				}
			} else {
				d.index++
				d.byteAvailable = true
			}
		}
	}
}

func (d *compressor) store() {
	if d.windowEnd > 0 && (d.windowEnd == maxStoreBlockSize || d.sync) {
		d.err = d.writeStoredBlock(d.window[:d.windowEnd])
		d.windowEnd = 0
	}
}

func (d *compressor) fillBlock(b []byte) int {
	n := copy(d.window[d.windowEnd:], b)
	d.windowEnd += n
	return n
}

func (d *compressor) storeHuff() {
	if d.windowEnd < len(d.window) && !d.sync || d.windowEnd == 0 {
		return
	}
	d.w.writeBlockHuff(false, d.window[:d.windowEnd])
	d.err = d.w.err
	d.windowEnd = 0
}

func (d *compressor) storeSnappy() {

	if d.windowEnd < maxStoreBlockSize {
		if !d.sync {
			return
		}

		if d.windowEnd < 128 {
			if d.windowEnd == 0 {
				return
			}
			if d.windowEnd <= 32 {
				d.err = d.writeStoredBlock(d.window[:d.windowEnd])
				d.tokens.n = 0
				d.windowEnd = 0
			} else {
				d.w.writeBlockHuff(false, d.window[:d.windowEnd])
				d.err = d.w.err
			}
			d.tokens.n = 0
			d.windowEnd = 0
			d.snap.Reset()
			return
		}
	}

	d.snap.Encode(&d.tokens, d.window[:d.windowEnd])

	if int(d.tokens.n) == d.windowEnd {
		d.err = d.writeStoredBlock(d.window[:d.windowEnd])

	} else if int(d.tokens.n) > d.windowEnd-(d.windowEnd>>4) {
		d.w.writeBlockHuff(false, d.window[:d.windowEnd])
		d.err = d.w.err
	} else {
		d.w.writeBlockDynamic(d.tokens.tokens[:d.tokens.n], false, d.window[:d.windowEnd])
		d.err = d.w.err
	}
	d.tokens.n = 0
	d.windowEnd = 0
}

func (d *compressor) write(b []byte) (n int, err error) {
	if d.err != nil {
		return 0, d.err
	}
	n = len(b)
	for len(b) > 0 {
		d.step(d)
		b = b[d.fill(d, b):]
		if d.err != nil {
			return 0, d.err
		}
	}
	return n, d.err
}

func (d *compressor) syncFlush() error {
	d.sync = true
	if d.err != nil {
		return d.err
	}
	d.step(d)
	if d.err == nil {
		d.w.writeStoredHeader(0, false)
		d.w.flush()
		d.err = d.w.err
	}
	d.sync = false
	return d.err
}

func (d *compressor) init(w io.Writer, level int) (err error) {
	d.w = newHuffmanBitWriter(w)

	switch {
	case level == NoCompression:
		d.window = make([]byte, maxStoreBlockSize)
		d.fill = (*compressor).fillBlock
		d.step = (*compressor).store
	case level == ConstantCompression:
		d.window = make([]byte, maxStoreBlockSize)
		d.fill = (*compressor).fillBlock
		d.step = (*compressor).storeHuff
	case level >= 1 && level <= 4:
		d.snap = newSnappy(level)
		d.window = make([]byte, maxStoreBlockSize)
		d.fill = (*compressor).fillBlock
		d.step = (*compressor).storeSnappy
	case level == DefaultCompression:
		level = 5
		fallthrough
	case 5 <= level && level <= 9:
		d.compressionLevel = levels[level]
		d.initDeflate()
		d.fill = (*compressor).fillDeflate
		if d.fastSkipHashing == skipNever {
			if useSSE42 {
				d.step = (*compressor).deflateLazySSE
			} else {
				d.step = (*compressor).deflateLazy
			}
		} else {
			if useSSE42 {
				d.step = (*compressor).deflateSSE
			} else {
				d.step = (*compressor).deflate

			}
		}
	default:
		return fmt.Errorf("flate: invalid compression level %d: want value in range [-2, 9]", level)
	}
	return nil
}

func (d *compressor) reset(w io.Writer) {
	d.w.reset(w)
	d.sync = false
	d.err = nil

	if d.snap != nil {
		d.snap.Reset()
		d.windowEnd = 0
		d.tokens.n = 0
		return
	}
	switch d.compressionLevel.chain {
	case 0:

		d.windowEnd = 0
	default:
		d.chainHead = -1
		for i := range d.hashHead {
			d.hashHead[i] = 0
		}
		for i := range d.hashPrev {
			d.hashPrev[i] = 0
		}
		d.hashOffset = 1
		d.index, d.windowEnd = 0, 0
		d.blockStart, d.byteAvailable = 0, false
		d.tokens.n = 0
		d.length = minMatchLength - 1
		d.offset = 0
		d.hash = 0
		d.ii = 0
		d.maxInsertIndex = 0
	}
}

func (d *compressor) close() error {
	if d.err != nil {
		return d.err
	}
	d.sync = true
	d.step(d)
	if d.err != nil {
		return d.err
	}
	if d.w.writeStoredHeader(0, true); d.w.err != nil {
		return d.w.err
	}
	d.w.flush()
	return d.w.err
}

func NewWriter(w io.Writer, level int) (*Writer, error) {
	var dw Writer
	if err := dw.d.init(w, level); err != nil {
		return nil, err
	}
	return &dw, nil
}

func NewWriterDict(w io.Writer, level int, dict []byte) (*Writer, error) {
	dw := &dictWriter{w}
	zw, err := NewWriter(dw, level)
	if err != nil {
		return nil, err
	}
	zw.d.fillWindow(dict)
	zw.dict = append(zw.dict, dict...)
	return zw, err
}

type dictWriter struct {
	w io.Writer
}

func (w *dictWriter) Write(b []byte) (n int, err error) {
	return w.w.Write(b)
}

type Writer struct {
	d    compressor
	dict []byte
}

func (w *Writer) Write(data []byte) (n int, err error) {
	return w.d.write(data)
}

func (w *Writer) Flush() error {

	return w.d.syncFlush()
}

func (w *Writer) Close() error {
	return w.d.close()
}

func (w *Writer) Reset(dst io.Writer) {
	if dw, ok := w.d.w.writer.(*dictWriter); ok {

		dw.w = dst
		w.d.reset(dw)
		w.d.fillWindow(w.dict)
	} else {

		w.d.reset(dst)
	}
}

func (w *Writer) ResetDict(dst io.Writer, dict []byte) {
	w.dict = dict
	w.d.reset(dst)
	w.d.fillWindow(w.dict)
}
