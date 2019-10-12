package udwGoParser

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestMustParseGoFunc(t *testing.T) {
	fn := MustParseGoFunc("func a()")
	udwTest.Equal(fn.GetName(), "a")
	udwTest.Equal(len(fn.GetInParameter()), 0)
	udwTest.Equal(len(fn.GetOutParameter()), 0)

	fn = MustParseGoFunc("func abc()bool")
	udwTest.Equal(fn.GetName(), "abc")
	udwTest.Equal(len(fn.GetInParameter()), 0)
	udwTest.Equal(len(fn.GetOutParameter()), 1)
	para := fn.GetOutParameter()[0]
	udwTest.Equal(para.GetName(), "")
	udwTest.Equal(para.GetType().String(), "bool")

	fn = MustParseGoFunc("func abc(a int,b string,c int)(d bool,e int)")
	udwTest.Equal(fn.GetName(), "abc")
	udwTest.Equal(len(fn.GetInParameter()), 3)
	udwTest.Equal(len(fn.GetOutParameter()), 2)
}
