package udwShmTest

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwTest"
	"math"
	"strings"
	"time"
)

type Marshaler struct {
	Marshal   func(obj interface{}) (data []byte, errMsg string)
	Unmarshal func(data []byte, obj interface{}) (errMsg string)
}

var gMarshaler Marshaler

func SetMarshaler(Marshal func(obj interface{}) (data []byte, errMsg string), Unmarshal func(data []byte, obj interface{}) (errMsg string)) {
	gMarshaler = Marshaler{
		Marshal:   Marshal,
		Unmarshal: Unmarshal,
	}
}
func Marshal(obj interface{}) (data []byte, errMsg string) {
	return gMarshaler.Marshal(obj)
}
func Unmarshal(data []byte, obj interface{}) (errMsg string) {
	return gMarshaler.Unmarshal(data, obj)
}

func RunTest() {
	testMsgV5_1()
	testMsgV5_2()
	testMsgV5_3()
}

func testMsgV5_1() {
	v5 := MsgV5{
		Data:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		Name:         strings.Repeat(`name123`, 1000),
		ByteArray:    [8]byte{22, 9, 29, 33, 2},
		IntArray:     [8]int{1, 29993, 232, 92903},
		Float64Slice: []float64{292.32, 2993.2, 1.0002, 0.232},
		IntPtr:       nil,
		IntPtr2: func() *int {
			i := 1000102002
			return &i
		}(),
		StringMsgV4Map: map[string]MsgV4{
			"stri": {
				Id: 80088,
			},
			"stri2": MsgV4(MsgV3{
				Id: 8088,
			}),
		},
		IsRoot:       true,
		U32:          math.MaxUint32,
		TimeDuration: time.Second * 20,
		TestTime:     time.Now(),
		IntListListList: [][][]int{
			{{
				1,
			}},
		},
		Float32Slice: []float32{
			1.12,
		},
		V3: MsgV3{
			Id: 1234,
		},
	}
	b2, errMsg := Marshal(v5)
	udwErr.PanicIfErrorMsg(errMsg)
	var v5Sub MsgV5
	errMsg = Unmarshal(b2, &v5Sub)
	udwErr.PanicIfErrorMsg(errMsg)
	MsgV5Equal(v5, v5Sub)
	udwTest.Ok(*v5Sub.IntPtr2 == 1000102002)

	var v5Sub_simpe MsgV5_Simple
	errMsg = Unmarshal(b2, &v5Sub_simpe)
	MsgV5Equal_Simple(v5, v5Sub_simpe)

}

func testMsgV5_2() {
	v5 := MsgV5{
		V4PtrList: []*MsgV4{
			nil,
		},
	}
	b2, errMsg := Marshal(v5)
	udwErr.PanicIfErrorMsg(errMsg)

	var v5Sub MsgV5
	errMsg = Unmarshal(b2, &v5Sub)
	udwErr.PanicIfErrorMsg(errMsg)
	MsgV5Equal(v5, v5Sub)
}

func testMsgV5_3() {
	v5 := MsgV5{
		IsRoot: true,
	}
	b2, errMsg := Marshal(v5)
	udwErr.PanicIfErrorMsg(errMsg)
	var v5Sub MsgV5
	errMsg = Unmarshal(b2, v5Sub)
	udwTest.Ok(errMsg != "")

	var v5Sub2 *MsgV5
	errMsg = Unmarshal(b2, v5Sub2)
	udwTest.Ok(errMsg != "")

	errMsg = Unmarshal(b2, &v5Sub2)
	udwTest.Ok(errMsg == "")

	var v5Sub3 **MsgV5
	errMsg = Unmarshal(b2, v5Sub3)
	udwTest.Ok(errMsg != "")

	var v5Sub4 *MsgV5 = &MsgV5{}
	errMsg = Unmarshal(b2, v5Sub4)
	udwTest.Ok(errMsg == "")
}
