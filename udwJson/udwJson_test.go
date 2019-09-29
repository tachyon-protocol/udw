package udwJson

import (
	"encoding/json"
	"fmt"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestJsLiteralObjectToJson(t *testing.T) {
	for input, expect := range map[string]struct {
		output      string
		hasErr      bool
		invalidJson bool
	}{
		`{"a":"\'v"}`: {
			output: `{"a":"'v"}`,
		},
		`{"a":"\\'v"}`: {
			output: `{"a":"\\'v"}`,
		},
		`{"a":"'v"}\`: {
			hasErr:      true,
			invalidJson: true,
		},
		`{"a":"\'v"}\`: {
			hasErr:      true,
			invalidJson: true,
		},
		`{"a":"\'v"}\\`: {
			output:      `{"a":"'v"}\\`,
			invalidJson: true,
		},
		`{"a":"\\\'v"}`: {
			output: `{"a":"\\'v"}`,
		},
		`{"\"a\'":"\\\'v"}`: {
			output: `{"\"a'":"\\'v"}`,
		},
		`{"\b\'":"\\\'v"}`: {
			output: `{"\b'":"\\'v"}`,
		},
	} {
		output, err := JsLiteralObjectToJson([]byte(input))
		udwTest.Ok(expect.hasErr == (err != nil), input)
		if err != nil {
			continue
		}
		udwTest.Ok(string(output) == expect.output, string(output), expect.output)
		err = json.Unmarshal(output, &struct{}{})
		if expect.invalidJson {
			udwTest.Ok(err != nil)
			return
		}
		if err != nil {
			fmt.Println("input:", input, "expect output:", expect.output, "output:", string(output))
			udwErr.PanicIfError(err)
		}
	}
}
