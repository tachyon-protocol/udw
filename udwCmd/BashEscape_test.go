package udwCmd_test

import (
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwRand"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestBashEscape(ot *testing.T) {
	udwTest.Equal(len(udwCmd.BashEscape("'")), 6)
	udwTest.Equal(len(udwCmd.BashEscape(udwCmd.BashEscape("'"))), 23)
	udwTest.Equal(len(udwCmd.BashEscape(udwCmd.BashEscape(udwCmd.BashEscape("'")))), 76)
	udwFile.MustDelete("/tmp/testFile")
	defer udwFile.MustDelete("/tmp/testFile")
	udwFile.MustMkdirForFile("/tmp/testFile/1.txt")
	for _, cas := range []string{
		" '",
		`abc 'abc'`,
		`abc '中文'`,
		`"!@#$%^&*()`,
		`\"'\'`,
		"\n\b\r\t",
		"\x01",
		"'",
		udwCmd.BashEscape("'"),
		udwCmd.BashEscape(udwCmd.BashEscape("'")),
	} {
		mustRunTestBashEscapeContent(cas)
	}

}

func RunTestBashEscapeFull() {
	for i := 1; i < 256; i++ {
		mustRunTestBashEscapeContent(string([]byte{byte(i)}))
	}
	for i := 0; i < 1000; i++ {
		bl := make([]byte, 3)
		for j := 0; j < 3; j++ {
			bl[j] = byte(udwRand.IntBetween(1, 255))
		}
		mustRunTestBashEscapeContent(string(bl))
	}
}

func mustRunTestBashEscapeContent(cas string) {
	cmd := `echo -n ` + udwCmd.BashEscape(cas) + ` > /tmp/testFile/1.txt`
	udwCmd.CmdBash(cmd).MustCombinedOutput()
	udwTest.Equal(string(udwFile.MustReadFile("/tmp/testFile/1.txt")), cas)
}
