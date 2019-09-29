package udwCmd_test

import (
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestBashEscape(ot *testing.T) {
	udwFile.MustDelete("testFile")
	defer udwFile.MustDelete("testFile")
	udwFile.MustMkdirForFile("testFile/1.txt")
	for _, cas := range []string{
		"'",
		" '",
		`abc 'abc'`,
		`abc '中文'`,
		`"!@#$%^&*()`,
		`\"'\'`,
		"\n\b\r\t",
		"\x01",
	} {
		cmd := `echo -n ` + udwCmd.BashEscape(cas) + ` > testFile/1.txt`
		udwCmd.MustRunInBash(cmd)
		udwTest.Equal(string(udwFile.MustReadFile("testFile/1.txt")), cas)
	}

}
