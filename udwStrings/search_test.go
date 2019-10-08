package udwStrings

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestIsInSliceBSearch(t *testing.T) {
	ss := []string{"1", "3", "4", "6"}
	for _, one := range ss {
		udwTest.Ok(IsInSliceBSearch(ss, one))
	}
	for _, one := range []string{"0", "2", "5", "9"} {
		udwTest.Ok(!IsInSliceBSearch(ss, one))
	}
}
