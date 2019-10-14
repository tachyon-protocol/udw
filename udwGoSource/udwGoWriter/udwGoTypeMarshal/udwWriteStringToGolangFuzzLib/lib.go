package udwWriteStringToGolangFuzzLib

import (
	"encoding/hex"
	"fmt"
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal"
	"github.com/tachyon-protocol/udw/udwHex"
	"github.com/tachyon-protocol/udw/udwTest"
	"go/format"
)

func Fuzz(data []byte) int {
	s := string(data)
	writeStringCheck(s)
	writeStringCheckASCII(s)
	return 1
}

func writeStringCheck(s string) {
	out := udwGoTypeMarshal.WriteStringToGolang(s)

	checkGoCanParseCorrect(out, s)
}

func writeStringCheckASCII(s string) {
	out := udwGoTypeMarshal.WriteStringToGolangASCII(s)
	for i := 0; i < len(out); i++ {
		b := out[i]
		if !(b >= 0x20 && b <= 0x7e) {
			panic(fmt.Errorf("[writeStringCheckASCII] fail at [%s] is not ACSII safe",
				udwHex.EncodeStringToString(s)))
		}
	}
	checkGoCanParseCorrect(out, s)
}

func checkGoCanParseCorrect(out string, s string) {
	outS := udwGoTypeMarshal.MustReadGoStringFromString(out)
	udwTest.Equal([]byte(outS), []byte(s))
	goFile := `package main

import (
	"encoding/hex"
	"os"
)
var a = ` + out + `

func main(){
	os.Stdout.Write([]byte(hex.EncodeToString([]byte(a))))
}
`
	_, err := format.Source([]byte(goFile))
	if err != nil {
		panic(fmt.Errorf("fail at %s err %s", out, err.Error()))
	}
	tmpFilePath := udwFile.NewTmpFilePathWithExt("go")
	defer udwFile.MustDelete(tmpFilePath)
	udwFile.MustWriteFileWithMkdir(tmpFilePath, []byte(goFile))
	out2 := udwCmd.CmdSlice([]string{"go", "run", tmpFilePath}).MustCombinedOutput()
	out3, err := hex.DecodeString(string(out2))
	udwTest.Equal(err, nil)
	udwTest.Equal([]byte(s), out3, out, []byte(out), []byte(s), []byte(out3))
}
