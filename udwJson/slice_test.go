package udwJson

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestMustMarshalStringArgumentToString(t *testing.T) {
	udwTest.Equal(MustMarshalStringArgumentToString("1", "2"), `["1","2"]`)
	udwTest.Equal(MustUnmarshalFromStringToStringSlice(`["1","2"]`), []string{"1", "2"})

	udwTest.BenchmarkWithRepeatNum(1e4, func() {
		MustUnmarshalFromStringToStringSlice(`["1","2"]`)
	})

	udwTest.BenchmarkWithRepeatNum(1e4, func() {
		MustMarshalStringArgumentToString("1", "2")
	})
}

func TestMustMarshalFirstKey(t *testing.T) {
	udwTest.Equal(MustMarshalKeyPrefix("1"), `["1`)
	udwTest.Equal(MustMarshalKeyPrefix("1", "2"), `["1","2`)
}

func TestPrefix(t *testing.T) {
	udwTest.Equal(MustMarshalKeyPrefix(), "")
	udwTest.Equal(MustMarshalKeyPrefix(""), `["`)
}

func TestMustUnmarshalFromStringToStringSliceByIndex(t *testing.T) {
	udwTest.Equal(MustUnmarshalFromStringToStringSliceByIndex(`["1","2"]`, 1), "2")
	udwTest.Equal(MustUnmarshalFromStringToStringSliceByIndex(`["1","2"]`, 0), "1")
}
