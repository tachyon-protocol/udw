package udwQueryOnlyUrl_test

import (
	"github.com/tachyon-protocol/udw/udwQueryOnlyUrl"
	"github.com/tachyon-protocol/udw/udwQueryOnlyUrl/udwQueryOnlyUrlFuzzLib"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestA(ot *testing.T) {
	udwQueryOnlyUrlFuzzLib.Fuzz([]byte("abc://?a=2&b=d"))
	udwQueryOnlyUrlFuzzLib.Fuzz([]byte("abc://?a=2&b=%4"))
	udwQueryOnlyUrlFuzzLib.Fuzz([]byte("abc://?a=2&b=%45"))
	udwQueryOnlyUrlFuzzLib.Fuzz([]byte("abc://?a=2&b=%450"))
	udwTest.Equal(udwQueryOnlyUrl.ParseQueryUrlObj("abc://?a=2&b=d").String(), "abc://?a=2&b=d")
}
