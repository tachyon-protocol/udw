package udwTest

import (
	"testing"
)

func TestEqual(ot *testing.T) {
	Equal(true, true)
	Equal([]byte{1}, []byte{1})
	AssertPanic(func() {
		Equal(true, false)
	})
	Ok(true)
	AssertPanic(func() {
		Ok(false)
	})
	AssertPanic(func() {
		Equal(int64(1), int(1))
	})
	Equal(map[string]interface{}{"a": 1}, map[string]interface{}{"a": 1})
	Equal(nil, nil)
	Equal([]string{}, []string{})
}

func TestV2(ot *testing.T) {
	msg := AssertPanic(func() {
		Ok(false)
	})
	Equal(msg, "ok fail")
	Ok(true)

	Equal([]byte{1, 2}, []byte{1, 2})

	AssertPanic(func() {
		Equal([]byte{1, 2}, []byte{1, 3})
	})

	Equal(nil, (*testing.T)(nil))
	Equal(nil, []string(nil))
	vv := []string(nil)
	Equal(nil, vv)
	Equal(nil, getStringSliceNil())
	Equal(getStringSliceNil(), nil)

	flag := 1
	msg = AssertPanic(func() {
		AssertPanic(func() {
			flag = 2
		})

		flag = 3
	})
	Equal(flag, 2)
	Equal(msg, "should panic")
}

func getStringSliceNil() []string {
	return nil
}

func TestEqualSpeed(ot *testing.T) {
	return

	Benchmark(func() {
		num := int(1e7)
		BenchmarkSetNum(num)

		for i := 0; i < num; i++ {
			Equal(i, i)
		}
	})
}
