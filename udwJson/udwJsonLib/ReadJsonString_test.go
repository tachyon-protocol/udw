package udwJsonLib_test

import (
	"github.com/tachyon-protocol/udw/udwJson"
	"github.com/tachyon-protocol/udw/udwJson/udwJsonLib"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestReadJsonString1(ot *testing.T) {
	inByte := []byte(`"00\ufffd0000"`)
	ctx := udwJsonLib.NewContextFromBuffer(inByte)
	out := udwJsonLib.ReadJsonString(ctx)

	s := ""
	udwJson.MustUnmarshal(inByte, &s)
	udwTest.Equal(out, s)
}
func TestReadJsonString2(t *testing.T) {
	inByte := []byte("\"00\xcf0000\"")
	ctx := udwJsonLib.NewContextFromBuffer(inByte)
	udwTest.AssertPanicWithErrorMessage(func() {
		udwJsonLib.ReadJsonString(ctx)
	}, "invalid utf8 code point")
}

func TestReadJsonString3(t *testing.T) {
	inByte := []byte("\"\\\"0\x82k\"")
	ctx := udwJsonLib.NewContextFromBuffer(inByte)
	udwTest.AssertPanicWithErrorMessage(func() {
		udwJsonLib.ReadJsonString(ctx)
	}, "invalid utf8 code point")
}

func TestReadJsonString4(t *testing.T) {
	inByte := []byte("\"\\u0080&\"")
	ctx := udwJsonLib.NewContextFromBuffer(inByte)
	s := udwJsonLib.ReadJsonString(ctx)

	ctx2 := &udwJsonLib.Context{}
	udwJsonLib.WriteJsonString(ctx2, s)
	udwTest.Equal(string(ctx2.WriterBytes()), "\"\\u0026\"")
}

func TestReadFloat64(t *testing.T) {
	inByte := []byte("-0")
	ctx := udwJsonLib.NewContextFromBuffer(inByte)
	f := udwJsonLib.ReadJsonFloat64(ctx)

	ctx2 := &udwJsonLib.Context{}
	udwJsonLib.WriteJsonFloat64(ctx2, f)
	udwTest.Equal(string(ctx2.WriterBytes()), "0")
	udwTest.Equal(f == 0, true)

}

func TestReadFloat32(t *testing.T) {
	inByte := []byte("-0")
	ctx := udwJsonLib.NewContextFromBuffer(inByte)
	f := udwJsonLib.ReadJsonFloat32(ctx)

	ctx2 := &udwJsonLib.Context{}
	udwJsonLib.WriteJsonFloat32(ctx2, f)
	udwTest.Equal(string(ctx2.WriterBytes()), "0")

}
