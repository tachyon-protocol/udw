package kkcflate

import (
	"math"
)

type hcode struct {
	code, len uint16
}

type huffmanEncoder struct {
	codes     []hcode
	freqcache []literalNode
	bitCount  [17]int32
}

type literalNode struct {
	literal uint16
	freq    int32
}

type levelInfo struct {
	level int32

	lastFreq int32

	nextCharFreq int32

	nextPairFreq int32

	needed int32
}

func (h *hcode) set(code uint16, length uint16) {
	h.len = length
	h.code = code
}

func maxNode() literalNode { return literalNode{math.MaxUint16, math.MaxInt32} }

func newHuffmanEncoder(size int) *huffmanEncoder {
	return &huffmanEncoder{codes: make([]hcode, size)}
}

func generateFixedLiteralEncoding() *huffmanEncoder {
	h := newHuffmanEncoder(maxNumLit)
	codes := h.codes
	var ch uint16
	for ch = 0; ch < maxNumLit; ch++ {
		var bits uint16
		var size uint16
		switch {
		case ch < 144:

			bits = ch + 48
			size = 8
			break
		case ch < 256:

			bits = ch + 400 - 144
			size = 9
			break
		case ch < 280:

			bits = ch - 256
			size = 7
			break
		default:

			bits = ch + 192 - 280
			size = 8
		}
		codes[ch] = hcode{code: reverseBits(bits, byte(size)), len: size}
	}
	return h
}

func generateFixedOffsetEncoding() *huffmanEncoder {
	h := newHuffmanEncoder(30)
	codes := h.codes
	for ch := range codes {
		codes[ch] = hcode{code: reverseBits(uint16(ch), 5), len: 5}
	}
	return h
}

var fixedLiteralEncoding *huffmanEncoder = generateFixedLiteralEncoding()
var fixedOffsetEncoding *huffmanEncoder = generateFixedOffsetEncoding()

func (h *huffmanEncoder) bitLength(freq []int32) int {
	var total int
	for i, f := range freq {
		if f != 0 {
			total += int(f) * int(h.codes[i].len)
		}
	}
	return total
}

const maxBitsLimit = 16

func (h *huffmanEncoder) bitCounts(list []literalNode, maxBits int32) []int32 {
	if maxBits >= maxBitsLimit {
		panic("flate: maxBits too large")
	}
	n := int32(len(list))
	list = list[0 : n+1]
	list[n] = maxNode()

	if maxBits > n-1 {
		maxBits = n - 1
	}

	var levels [maxBitsLimit]levelInfo

	var leafCounts [maxBitsLimit][maxBitsLimit]int32

	for level := int32(1); level <= maxBits; level++ {

		levels[level] = levelInfo{
			level:        level,
			lastFreq:     list[1].freq,
			nextCharFreq: list[2].freq,
			nextPairFreq: list[0].freq + list[1].freq,
		}
		leafCounts[level][level] = 2
		if level == 1 {
			levels[level].nextPairFreq = math.MaxInt32
		}
	}

	levels[maxBits].needed = 2*n - 4

	level := maxBits
	for {
		l := &levels[level]
		if l.nextPairFreq == math.MaxInt32 && l.nextCharFreq == math.MaxInt32 {

			l.needed = 0
			levels[level+1].nextPairFreq = math.MaxInt32
			level++
			continue
		}

		prevFreq := l.lastFreq
		if l.nextCharFreq < l.nextPairFreq {

			n := leafCounts[level][level] + 1
			l.lastFreq = l.nextCharFreq

			leafCounts[level][level] = n
			l.nextCharFreq = list[n].freq
		} else {

			l.lastFreq = l.nextPairFreq

			copy(leafCounts[level][:level], leafCounts[level-1][:level])
			levels[l.level-1].needed = 2
		}

		if l.needed--; l.needed == 0 {

			if l.level == maxBits {

				break
			}
			levels[l.level+1].nextPairFreq = prevFreq + l.lastFreq
			level++
		} else {

			for levels[level-1].needed > 0 {
				level--
			}
		}
	}

	if leafCounts[maxBits][maxBits] != n {
		panic("leafCounts[maxBits][maxBits] != n")
	}

	bitCount := h.bitCount[:maxBits+1]
	bits := 1
	counts := &leafCounts[maxBits]
	for level := maxBits; level > 0; level-- {

		bitCount[bits] = counts[level] - counts[level-1]
		bits++
	}
	return bitCount
}

func (h *huffmanEncoder) assignEncodingAndSize(bitCount []int32, list []literalNode) {
	code := uint16(0)
	for n, bits := range bitCount {
		code <<= 1
		if n == 0 || bits == 0 {
			continue
		}

		chunk := list[len(list)-int(bits):]

		sortbyLiteral(chunk)

		for _, node := range chunk {
			h.codes[node.literal] = hcode{code: reverseBits(code, uint8(n)), len: uint16(n)}
			code++
		}
		list = list[0 : len(list)-int(bits)]
	}
}

func (h *huffmanEncoder) generate(freq []int32, maxBits int32) {
	if h.freqcache == nil {

		h.freqcache = make([]literalNode, maxNumLit+1)
	}
	list := h.freqcache[:len(freq)+1]

	count := 0

	for i, f := range freq {
		if f != 0 {
			list[count] = literalNode{uint16(i), f}
			count++
		} else {
			list[count] = literalNode{}
			h.codes[i].len = 0
		}
	}
	list[len(freq)] = literalNode{}

	list = list[:count]
	if count <= 2 {

		for i, node := range list {

			h.codes[node.literal].set(uint16(i), 1)
		}
		return
	}
	sortbyFreq(list)

	bitCount := h.bitCounts(list, maxBits)

	h.assignEncodingAndSize(bitCount, list)
}
