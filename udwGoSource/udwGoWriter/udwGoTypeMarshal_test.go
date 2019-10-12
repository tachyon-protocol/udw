package udwGoWriter

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"runtime"
	"testing"
)

func TestMustMarshal(ot *testing.T) {
	udwTest.Equal(MustWriteObjectToMainPackage("abc"), "`abc`")
	udwTest.Equal(MustWriteObjectToMainPackage(map[string]string{"abc": "abc"}), "map[string]string{\n`abc`:`abc`,\n}")
	udwTest.Equal(MustWriteObjectToMainPackage([]string{`abc`, `1`}), "[]string{\n`abc`,\n`1`,\n}")
	udwTest.Equal(MustWriteObjectToMainPackage([]byte(`abc`)), "[]byte(`abc`)")
	MustWriteObjectToMainPackage(runtime.MemStats{})
	MustWriteObjectToMainPackage(runtime.MemProfileRecord{})
}
