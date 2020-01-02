package udwSortedMap

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestStringToFloat64Map_Keys(t *testing.T) {
	m := NewStringToFloat64Map()
	m.Set("A", 1)
	m.Set("B", 2)
	m.Set("C", 3)
	keys := m.Keys(AscSortByValue)
	udwTest.Ok(keys[0] == "A")
	udwTest.Ok(keys[1] == "B")
	udwTest.Ok(keys[2] == "C")
	keys = m.Keys(DescSortByValue)
	udwTest.Ok(keys[0] == "C")
	udwTest.Ok(keys[1] == "B")
	udwTest.Ok(keys[2] == "A")
}
