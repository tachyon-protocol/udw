package udwRandFast

import (
	"testing"
)

func TestMustReadV3(ot *testing.T) {
	DoTestRead(MustRead)
	DoTestSpeed(MustRead, "MustReadV3")
}
