package udwEncodeBinary

import (
	"testing"
	"bytes"
	"github.com/tachyon-protocol/udwTest"
)

type TestT struct{
	S string
}

func TestMarshal(ot *testing.T){
	{
		b,errMsg:=Marshal("1")
		udwTest.Ok(errMsg=="",errMsg)
		udwTest.Ok(bytes.Equal([]byte{1,0x31},b),b)

		var s string
		errMsg = Unmarshal([]byte{1,0x31},&s)
		udwTest.Ok(errMsg=="",errMsg)
		udwTest.Ok(s=="1",s)
	}
	{
		obj:=TestT{}
		obj.S = "2"
		b,errMsg:=Marshal(obj)
		udwTest.Ok(errMsg=="",errMsg)
		udwTest.Ok(bytes.Equal([]byte{1,0x32},b),b)

		var obj2 TestT
		errMsg = Unmarshal([]byte{1,0x33},&obj2)
		udwTest.Ok(errMsg=="",errMsg)
		udwTest.Ok(obj2.S=="3",obj2.S)
	}

}