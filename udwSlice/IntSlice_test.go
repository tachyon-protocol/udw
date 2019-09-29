package udwSlice

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestIntSliceRemoveAt(ot *testing.T) {
	s := []int{1, 2, 3}
	IntSliceRemoveAt(&s, 1)
	udwTest.Equal(s, []int{1, 3})
}
