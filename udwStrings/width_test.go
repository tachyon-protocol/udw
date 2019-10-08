package udwStrings

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestGetWidth(t *testing.T) {
	for _, cas := range []struct {
		input  string
		output int
	}{
		{
			input:  "魑魅魍魉",
			output: 2 * 4,
		},
		{
			input:  "中文。123",
			output: 2*3 + 3,
		},
		{
			input:  "�����",
			output: 5,
		},
	} {
		udwTest.Ok(cas.output == GetMonospaceWidth(cas.input), cas.input)

	}
}
