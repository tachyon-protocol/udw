package udwStrconv

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseHexUint64(s string) (u uint64, errMsg string) {
	if strings.HasPrefix(s, "0x") == false {
		return 0, "ParseHexUint64 pvve7ddbf6 do not have prefix 0x"
	}
	u, err := strconv.ParseUint(s[2:], 16, 64)
	if err != nil {
		return 0, "ParseHexUint64 tx6xpvmgqf " + err.Error()
	}
	return u, ""
}

func MustParseInt0xHex(f string) int {
	if !strings.HasPrefix(f, "0x") {
		panic(fmt.Errorf("[MustParseInt0xHex] not 0x hex string[%s]", f))
	}
	i, err := strconv.ParseInt(f[2:], 16, 64)
	if err != nil {
		panic(err)
	}
	return int(i)
}

func MustParseIntHex(f string) int {
	i, err := strconv.ParseInt(f, 16, 64)
	if err != nil {
		panic(err)
	}
	return int(i)
}

func FormatUint64Hex(i uint64) string {
	return "0x" + strconv.FormatUint(i, 16)
}

func FormatUint64HexPadding8(i uint64) string {
	s := strconv.FormatUint(i, 16)
	if len(s) < 8 {
		s = strings.Repeat("0", 8-len(s)) + s
	}
	return s
}

func FormatUint64HexPadding(i uint64) string {
	s := strconv.FormatUint(i, 16)
	if len(s) < 16 {
		s = strings.Repeat("0", 16-len(s)) + s
	}
	return s
}

func MustParseUint64Hex(f string) uint64 {
	i, err := strconv.ParseUint(f, 16, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func FormatUint64HexPaddingWithZeroPrefix(i uint64, paddingToWidth int) string {
	s := strconv.FormatUint(i, 16)
	if len(s) < paddingToWidth {
		s = strings.Repeat("0", paddingToWidth-len(s)) + s
	}
	return s
}

func ParseUint64HexWithout0xOrZero(s string) (i uint64) {
	u, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0
	}
	return u
}
