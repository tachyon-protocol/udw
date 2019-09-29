package udwTypeTransform

import "testing"
import (
	"github.com/tachyon-protocol/udw/udwTest"
)

type StringTranT1 struct {
	T2 StringTranT2
}
type StringTranT2 string

func TestStringTransformSubType(ot *testing.T) {
	in := &StringTranT1{
		T2: "6",
	}
	err := StringTransformSubType(in, map[string]map[string]string{
		"github.com/tachyon-protocol/udw/udwTypeTransform.StringTranT2": {
			"6": "Fire",
		},
	})
	udwTest.Equal(err, nil)
	udwTest.Equal(in.T2, StringTranT2("Fire"))
}
