package udwStrconv

import (
	"fmt"
	"strconv"
	"strings"
)

func AtoIDefault0(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func ParseIntDefault0(s string) int {
	return AtoIDefault0(s)
}

func ParseInt64Default0(s string) int64 {
	i64, _ := strconv.ParseInt(s, 10, 64)
	return i64
}

func ParseIntWithDefaultNumber(s string, defaultNumber int) int {
	if s == "" {
		return defaultNumber
	}
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultNumber
	}
	return int(i64)
}

func FormatInt(i int) string {
	return strconv.Itoa(i)
}

func FormatIntPadding(i int, paddingToWidth int) string {
	return FormatIntPaddingWithChars(i, paddingToWidth, " ")
}

func FormatIntPaddingWithZeroPre(i int, paddingToWidth int) string {
	return FormatIntPaddingWithChars(i, paddingToWidth, "0")
}

func FormatIntPaddingWithChars(i int, paddingToWidth int, paddingStr string) string {
	s := strconv.Itoa(i)
	if len(s) < paddingToWidth {
		s = strings.Repeat(paddingStr, paddingToWidth-len(s)) + s
	}
	return s
}

func FormatInt64(i int64) string {
	return strconv.FormatInt(i, 10)
}

func FormatUint64(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func FormatUint64Padding(i uint64) string {
	s := strconv.FormatUint(i, 10)
	if len(s) < 20 {
		s = strings.Repeat("0", 20-len(s)) + s
	}
	return s
}

func MustParseInt(f string) int {
	i, err := strconv.Atoi(f)
	if err != nil {
		panic(err)
	}
	return i
}

func ParseInt(f string) (int, error) {
	return strconv.Atoi(f)
}

func MustParseInt64(f string) int64 {
	i, err := strconv.ParseInt(f, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func MustParseUint64(f string) uint64 {
	i, err := strconv.ParseUint(f, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func FormatBool(f bool) string {
	return strconv.FormatBool(f)
}

func FormatBoolFalseEmpty(f bool) string {
	if f == true {
		return "true"
	} else {
		return ""
	}
}

func FormatFloatPrec0(f float64) string {
	return strconv.FormatFloat(f, 'f', 0, 64)
}

func FormatFloatPrec2(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func FormatFloatPrec4(f float64) string {
	return strconv.FormatFloat(f, 'f', 4, 64)
}

func FormatFloatPrec6(f float64) string {
	return strconv.FormatFloat(f, 'f', 6, 64)
}

func FormatFloatPrec8(f float64) string {
	return strconv.FormatFloat(f, 'f', 8, 64)
}

func FormatFloatPercentPadding(f float64) string {
	s := FormatFloatPrec2(f*100) + "%"
	if len(s) < 7 {
		s = strings.Repeat(" ", 7-len(s)) + s
	}
	return s
}

func FormatFloatPercentPaddingPrec4(f float64) string {
	s := FormatFloatPrec4(f*100) + "%"
	if len(s) < 9 {
		s = strings.Repeat(" ", 9-len(s)) + s
	}
	return s
}

func FormatFloatPercentPrec2(f float64) string {
	return FormatFloatPrec2(f*100) + "%"
}

func FormatFloatPercentPrec4(f float64) string {
	return FormatFloatPrec4(f*100) + "%"
}

func ParseFloat64(f string) (float64, error) {
	return strconv.ParseFloat(f, 64)
}

func MustParseFloat64(f string) float64 {
	out, err := strconv.ParseFloat(f, 64)
	if err != nil {
		panic(err)
	}
	return out
}

func ParseFloat64Default0(s string) float64 {
	out, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return out
}

func ParseBoolDefaultFalse(s string) bool {
	out, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return out
}

func MustParseBool(f string) bool {
	out, err := strconv.ParseBool(f)
	if err != nil {
		panic(err)
	}
	return out
}

func InterfaceToString(m interface{}) string {
	s := fmt.Sprint(m)
	return s
}

func intBeyond1000ToScientificNotationWithNum(in int, num int) string {
	if in > -1000 && in < 1000 {
		return FormatInt(in)
	}
	format := "%." + FormatInt(num) + "e"
	return fmt.Sprintf(format, float64(in))
}
func IntBeyond1000ToScientificNotationDefault(in int) string {
	return intBeyond1000ToScientificNotationWithNum(in, 2)
}

func Uint8ToBool0AsFalse(b uint8) bool {
	if b == 0 {
		return false
	} else {
		return true
	}
}

func GetPercent(success int, total int) string {
	if total == 0 {
		return "0%"
	}
	return FormatFloatPercentPrec2(float64(success) / float64(total))
}
