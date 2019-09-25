package udwDebug

import (
	"bytes"
	"encoding/hex"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"strings"
)

func HexDumpWithAddrOffset(data []byte, addrOffset uint64) string {
	outBuf := bytes.Buffer{}
	used := 0
	rightCharsBuf := make([]byte, 16)
	for i := range data {
		if used == 0 {
			outBuf.WriteString(udwStrconv.FormatUint64HexPadding8(addrOffset))
			outBuf.WriteString("  ")
		}
		outBuf.WriteString(hex.EncodeToString([]byte{data[i]}))
		outBuf.WriteByte(' ')
		rightCharsBuf[i%16] = asciiToChar(data[i])

		used++
		addrOffset++
		if used == 8 {
			outBuf.WriteByte(' ')
		}
		if used == 16 {
			outBuf.WriteString(" |")
			outBuf.Write(rightCharsBuf)
			outBuf.WriteString("|\n")
			used = 0
		}
	}
	if used != 0 {
		toAddSpace := (16-used)*3 + 1
		if used <= 7 {
			toAddSpace += 1
		}
		outBuf.WriteString(strings.Repeat(" ", toAddSpace))

		outBuf.WriteString("|")
		outBuf.Write(rightCharsBuf[:used])
		outBuf.WriteString("|\n")
	}
	return outBuf.String()
}

func asciiToChar(b byte) byte {
	if b < 32 || b > 126 {
		return '.'
	}
	return b
}
