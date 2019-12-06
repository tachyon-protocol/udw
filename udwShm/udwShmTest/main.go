package udwShmTest

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwJson"
	"github.com/tachyon-protocol/udw/udwTest"
	"math"
	"reflect"
	"strconv"
	"time"
)

type MsgV5 struct {
	Data            []byte
	Name            string
	ByteArray       [8]byte
	IntArray        [8]int
	Float64Slice    []float64
	IntPtr          *int
	IntPtr2         *int
	StringMsgV4Map  map[string]MsgV4
	IntStringPtrMap map[int]*string
	MapPtrKey       map[string]int
	V3              MsgV3
	V3List          []MsgV3
	UseRecord       []Record
	V4              MsgV4
	V4List          []MsgV4
	PtrV4List       *[]MsgV4
	V4PtrList       []*MsgV4
	TestTime        time.Time
	IsRoot          bool
	U32             uint32
	PtrTime         *time.Time
	TimeDuration    time.Duration
	Skip1           int `json:"-"`
	Skip2           int `json:"-,omitempty"`
	skip3           int
	Float32Slice    []float32
	IntListListList [][][]int
	F32             float32
}

type MsgV3 struct {
	Id int
}

type MsgV4 MsgV3

type Inter interface{}

type Record struct {
	Id   int64
	Name string
	F    float64
}

type MsgV5_Simple struct {
	StringMsgV4Map map[string]MsgV4
}

func MsgV5Equal(a MsgV5, b MsgV5) {
	udwTest.Equal(a.Data, b.Data)
	udwTest.Ok(bytes.Equal(a.Data, b.Data))
	udwTest.Ok(a.Name == b.Name)
	udwTest.Ok(a.ByteArray == b.ByteArray)
	udwTest.Ok(a.IntArray == b.IntArray, a.IntArray, b.IntArray)
	udwTest.Ok(len(a.Float64Slice) == len(b.Float64Slice))
	for i := 0; i < len(a.Float64Slice); i++ {
		af, bf := a.Float64Slice[i], b.Float64Slice[i]
		udwTest.Ok(math.Float64bits(af) == math.Float64bits(bf))
	}
	if a.IntPtr == nil {
		udwTest.Ok(b.IntPtr == nil)
	} else {
		udwTest.Ok(*a.IntPtr == *b.IntPtr, a.IntPtr, b.IntPtr)
	}
	if a.IntPtr2 == nil {
		udwTest.Ok(b.IntPtr2 == nil)
	} else {
		udwTest.Ok(*a.IntPtr2 == *b.IntPtr2, a.IntPtr2, b.IntPtr2)
	}
	if len(a.StringMsgV4Map) != 0 || len(b.StringMsgV4Map) != 0 {
		udwTest.Ok(reflect.DeepEqual(a.StringMsgV4Map, b.StringMsgV4Map))
	}
	udwTest.Ok(len(a.IntStringPtrMap) == len(b.IntStringPtrMap))
	if len(a.IntStringPtrMap) > 0 {
		for i := range a.IntStringPtrMap {
			if stringPtrToString(a.IntStringPtrMap[i]) != stringPtrToString(b.IntStringPtrMap[i]) {
				panic("fail " + strconv.Itoa(i))
			}
		}
	}
	if len(a.MapPtrKey) != 0 || len(b.MapPtrKey) != 0 {
		udwTest.Ok(reflect.DeepEqual(a.MapPtrKey, b.MapPtrKey))
	}
	udwTest.Ok(a.V3.Id == b.V3.Id, a.V3.Id, b.V3.Id)
	if len(a.V3List) != 0 || len(b.V3List) != 0 {
		udwTest.Ok(reflect.DeepEqual(a.V3List, b.V3List))
	} else {
		udwTest.Ok(len(a.V3List) == 0 && len(b.V3List) == 0)
	}
	if len(a.UseRecord) != 0 || len(b.UseRecord) != 0 {
		udwTest.Ok(len(a.UseRecord) == len(b.UseRecord), a.UseRecord, b.UseRecord)
		for i := 0; i < len(a.UseRecord); i++ {
			av := a.UseRecord[i]
			bv := b.UseRecord[i]
			udwTest.Ok(av.Id == bv.Id)
			udwTest.Ok(av.Name == bv.Name)
			udwTest.Ok(math.Float64bits(av.F) == math.Float64bits(bv.F))
		}
	} else {
		udwTest.Ok(len(a.UseRecord) == 0 && len(b.UseRecord) == 0)
	}
	udwTest.Ok(a.V4.Id == b.V4.Id)
	if len(a.V4List) != 0 || len(a.V4List) != 0 {
		udwTest.Ok(len(a.V4List) == len(b.V4List))
		for i := 0; i < len(a.V4List); i++ {
			udwTest.Ok(a.V4List[i].Id == b.V4List[i].Id, i)
		}
	}
	if a.PtrV4List != nil && len(*a.PtrV4List) > 0 {
		av := *a.PtrV4List
		bv := *b.PtrV4List
		udwTest.Ok(len(av) == len(bv))
		for i := 0; i < len(av); i++ {
			udwTest.Ok(av[i].Id == bv[i].Id, i, av[i].Id, bv[i].Id)
		}
	} else {
		av := udwJson.MustMarshalIndentToString(a.PtrV4List)
		bv := udwJson.MustMarshalIndentToString(b.PtrV4List)
		udwTest.Ok(b.PtrV4List == nil || len(*b.PtrV4List) == 0, av, bv)
	}
	if len(a.V4PtrList) != 0 || len(b.V4PtrList) != 0 {
		for idx := 0; idx < len(a.V4PtrList); idx++ {
			av := a.V4PtrList[idx]
			bv := b.V4PtrList[idx]
			if av == nil {
				udwTest.Ok(bv == nil)
			} else {
				fmt.Println(av.Id, bv, idx)
				udwTest.Ok(av.Id == bv.Id)
			}
		}
	}
	udwTest.Ok(a.TestTime.Equal(b.TestTime))
	udwTest.Ok(a.IsRoot == b.IsRoot)
	udwTest.Ok(a.U32 == b.U32)
	if a.PtrTime != nil {
		udwTest.Ok(a.PtrTime.Equal(*b.PtrTime))
	} else {
		udwTest.Ok(b.PtrTime == nil)
	}
	udwTest.Ok(a.TimeDuration == b.TimeDuration)
	udwTest.Ok(a.Skip1 == b.Skip1)
	udwTest.Ok(a.Skip2 == b.Skip2)
	udwTest.Ok(a.skip3 == b.skip3)
	udwTest.Ok(len(a.Float32Slice) == len(b.Float32Slice))
	for i := range a.Float32Slice {
		udwTest.Ok(a.Float32Slice[i] == b.Float32Slice[i])
	}
}

func MsgV5Equal_Simple(a MsgV5, b MsgV5_Simple) {
	if len(a.StringMsgV4Map) != 0 || len(b.StringMsgV4Map) != 0 {
		udwTest.Ok(reflect.DeepEqual(a.StringMsgV4Map, b.StringMsgV4Map))
	}
}

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
