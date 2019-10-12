package udwGoParser

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestGoSourceRemoveComment(ot *testing.T) {
	content := []byte("// abc\n/*  abc\n*/\npackage abc\n\nfunc b(){\n}")
	outContent := goSourceRemoveComment(content, nil)
	udwTest.Equal(string(outContent), "      \n       \n  \npackage abc\n\nfunc b(){\n}")
}
