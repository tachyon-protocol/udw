package udwGoParser

import (
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoReader"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal"
)

func mustReadGoString(r *udwGoReader.Reader) (output []byte) {
	return udwGoTypeMarshal.MustReadGoString(r)
}

func mustReadGoChar(r *udwGoReader.Reader) []byte {
	return udwGoTypeMarshal.MustReadGoChar(r)
}
