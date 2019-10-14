package udwGoWriter

import "github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal"

func WriteStringToGolang(s string) string {
	return udwGoTypeMarshal.WriteStringToGolang(s)
}
