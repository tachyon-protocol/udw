package udwStrings

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestAsciiToLower(t *testing.T) {
	for input, expect := range map[string]string{
		"AbC":                "abc",
		"ÀàÂâ":               "ÀàÂâ",
		"�":                  "�",
		string([]byte{0x89}): string([]byte{0x89}),
	} {

		udwTest.Ok(AsciiToLower(input) == expect)
	}
}
