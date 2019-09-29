package udwReflect_test

import (
	"github.com/tachyon-protocol/udw/udwReflect"
	"github.com/tachyon-protocol/udw/udwTest"
	"reflect"
	"testing"
)

type GetAllFieldT1 struct {
	GetAllFieldT3
	*GetAllFieldT4
	B int
}

type GetAllFieldT2 struct {
	A int
	B int
	C int
}

type GetAllFieldT3 struct {
	A int
	B int
	GetAllFieldT2
}
type GetAllFieldT4 struct {
	A int
	D int
}

type GetAllFieldT5 struct {
	GetAllFieldT6
	A int
}
type GetAllFieldT6 int

func TestStructGetAllField(ot *testing.T) {
	t1 := reflect.TypeOf(&GetAllFieldT1{})
	ret := udwReflect.StructGetAllField(t1)
	udwTest.Equal(len(ret), 7)
	udwTest.Equal(ret[0].Name, "GetAllFieldT3")
	udwTest.Equal(ret[1].Name, "GetAllFieldT4")
	udwTest.Equal(ret[2].Name, "B")
	udwTest.Equal(ret[2].Index, []int{2})
	udwTest.Equal(ret[3].Name, "A")
	udwTest.Equal(ret[3].Index, []int{0, 0})
	udwTest.Equal(ret[4].Name, "GetAllFieldT2")
	udwTest.Equal(ret[5].Name, "C")
	udwTest.Equal(ret[5].Index, []int{0, 2, 2})
	udwTest.Equal(ret[6].Name, "D")
	udwTest.Equal(ret[6].Index, []int{1, 1})

	ret = udwReflect.StructGetAllField(reflect.TypeOf(&GetAllFieldT5{}))
	udwTest.Equal(len(ret), 2)

}

func TestStructGetAllFieldMap(ot *testing.T) {
	t1 := reflect.TypeOf(&GetAllFieldT1{})
	ret := udwReflect.StructGetAllFieldMap(t1)
	udwTest.Equal(ret["A"].Index, []int{0, 0})
	udwTest.Equal(ret["B"].Index, []int{2})
	udwTest.Equal(ret["C"].Index, []int{0, 2, 2})
	udwTest.Equal(ret["D"].Index, []int{1, 1})
	udwTest.Equal(len(ret), 7)

	ret = udwReflect.StructGetAllFieldMap(reflect.TypeOf(&GetAllFieldT5{}))
	udwTest.Equal(ret["A"].Index, []int{1})
	udwTest.Equal(len(ret), 2)
}
