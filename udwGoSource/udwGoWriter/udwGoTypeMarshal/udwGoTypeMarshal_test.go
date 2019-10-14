package udwGoTypeMarshal

import (
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
