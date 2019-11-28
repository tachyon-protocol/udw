package udwGoTypeMarshal

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"reflect"
	"testing"
	"time"
)

func TestIsObjHasInitAlloc(t *testing.T) {
	testFn := func(obj interface{}, expectAlloc bool) {
		isAlloc := IsObjHasInitAlloc(reflect.ValueOf(obj), reflect.TypeOf(obj))
		if isAlloc != expectAlloc {
			panic("fail")
		}
	}
	testFn("1", false)
	testFn(map[string]string{}, false)
	testFn(struct {
		M map[string]string
	}{
		M: map[string]string{"1": "1"},
	}, true)
	testFn(time.Now(), true)
}

type T1_1 struct {
	F1 T1
	F2 [8]byte
	F3 float64
	F4 *int
	F5 *[]*int
	F6 map[string]string
}
type T1 struct {
	*T2
}

type T2 struct {
	*T3
	T2_F1 int
}
type T3 struct{}

func TestMustWriteObjectToMainPackage(t *testing.T) {
	outS := MustWriteObjectToMainPackage(1)
	udwTest.Equal(outS, "1")
	outS = MustWriteObjectToMainPackage(&T3{})
	udwTest.Equal(outS, "&udwGoTypeMarshal.T3{\n}")
	outS = MustWriteObjectToMainPackage(&T2{})
	udwTest.Equal(outS, "&udwGoTypeMarshal.T2{\n}")
	outS = MustWriteObjectToMainPackage(&T1{})
	udwTest.Equal(outS, "&udwGoTypeMarshal.T1{\n}")
	outS = MustWriteObjectToMainPackage(&T1_1{})
	udwTest.Equal(outS, "&udwGoTypeMarshal.T1_1{\n}")
	outS = MustWriteObjectToMainPackage(&T1_1{F3: 1.5})
	udwTest.Equal(outS, "&udwGoTypeMarshal.T1_1{\nF3:1.5,\n}")
	outS = MustWriteObjectToMainPackage(&T1_1{F1: T1{T2: &T2{
		T2_F1: 1,
	}}})
	udwTest.Equal(outS, `&udwGoTypeMarshal.T1_1{
F1:udwGoTypeMarshal.T1{
T2:&udwGoTypeMarshal.T2{
T2_F1:1,
},
},
}`)
	outS = MustWriteObjectToMainPackage(&T1_1{F4: func() *int { _a := 3; return &_a }()})
	udwTest.Equal(outS, "&udwGoTypeMarshal.T1_1{\nF4:udwGoTypeMarshalLib.PtrInt(3),\n}")
	outS = MustWriteObjectToMainPackage(&T1_1{F5: &[]*int{}})
	udwTest.Equal(outS, "&udwGoTypeMarshal.T1_1{\n}")
	outS = MustWriteObjectToMainPackage(&T1_1{F5: &[]*int{func() *int { _a := 5; return &_a }()}})
	udwTest.Equal(outS, "&udwGoTypeMarshal.T1_1{\nF5:&[]*int{\nudwGoTypeMarshalLib.PtrInt(5),\n},\n}")
}
