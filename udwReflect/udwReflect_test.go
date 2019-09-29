package udwReflect_test

import (
	"github.com/tachyon-protocol/udw/udwReflect"
	"github.com/tachyon-protocol/udw/udwTest"
	"reflect"
	"testing"
)

type ta struct {
}

func TestGetFullName(ot *testing.T) {
	name := udwReflect.GetTypeFullName(reflect.TypeOf(""))
	udwTest.Equal(name, "string")

	name = udwReflect.GetTypeFullName(reflect.TypeOf(1))
	udwTest.Equal(name, "int")

	name = udwReflect.GetTypeFullName(reflect.TypeOf(&ta{}))
	udwTest.Equal(name, "github.com/tachyon-protocol/udw/udwReflect_test.ta")

	name = udwReflect.GetTypeFullName(reflect.TypeOf([]string{}))
	udwTest.Equal(name, "")

}
