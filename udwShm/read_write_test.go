package udwShm

import (
	"github.com/tachyon-protocol/udw/udwShm/udwShmTest"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestMarshal(t *testing.T) {
	data, errMsg := Marshal(1)
	udwTest.Ok(errMsg == ``, errMsg)
	var i int
	Unmarshal(data, &i)
	udwTest.Ok(i == 1)
}

func TestMarshal17(t *testing.T) {
	var L0 uint = 10
	var L1 = &L0
	var L2 = &L1
	var L3 = &L2
	data, errMsg := Marshal(L3)
	udwTest.Ok(errMsg == ``, errMsg)
	var After ************uint
	errMsg = Unmarshal(data, &After)
	udwTest.Ok(errMsg == `` && ************After == 10, errMsg)
}

func TestRun1(t *testing.T) {
	udwShmTest.SetMarshaler(Marshal, Unmarshal)
	udwShmTest.RunTest()
}
