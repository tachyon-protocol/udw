package udwTypeTransform

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"reflect"
	"testing"
)

func TestManager(ot *testing.T) {
	Int := 0
	ArrMapStringInt := []map[string]int{}
	type T1 struct {
		A int
		B string
	}
	ArrStruct := []T1{}
	StringSlice := []string{}
	mapStringString := map[string]string{}
	testCaseTable := []struct {
		in  interface{}
		out interface{}
		exp interface{}
	}{
		{1, &Int, 1},
		{int64(1), &Int, 1},
		{
			[]map[string]string{
				{
					"a": "1",
				},
				{
					"b": "1",
				},
			},
			&ArrMapStringInt,
			[]map[string]int{
				{
					"a": 1,
				},
				{
					"b": 1,
				},
			},
		},
		{
			[]map[string]string{
				{
					"A": "1",
					"B": "abc",
					"C": "abd",
				},
				{
					"A": "",
					"B": "",
					"C": "abd",
				},
			},
			&ArrStruct,
			[]T1{
				{
					A: 1,
					B: "abc",
				},
				{
					A: 0,
					B: "",
				},
			},
		},
		{
			[]interface{}{
				"1",
				"2",
			},
			&StringSlice,
			[]string{
				"1",
				"2",
			},
		},
		{
			struct {
				A int
				a int
			}{
				A: 1,
				a: 1,
			},
			&mapStringString,
			map[string]string{
				"A": "1",
			},
		},
	}
	for i, testCase := range testCaseTable {
		err := Transform(testCase.in, testCase.out)
		udwTest.Equal(err, nil, "fail at %d", i)
		udwTest.Equal(reflect.ValueOf(testCase.out).Elem().Interface(), testCase.exp, "fail at %d", i)
	}
}

func TestMapToStruct(t *testing.T) {
	in := map[string]string{
		"a": "1",
		"B": "2",
	}
	type tOut struct {
		A string `typeTransform:"a"`
		B string
	}
	var out tOut
	MustTransform(in, &out)
	udwTest.Equal(out.B, "2")
	udwTest.Equal(out.A, "1")
}
