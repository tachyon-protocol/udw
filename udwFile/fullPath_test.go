package udwFile

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"path/filepath"
	"testing"
)

func TestMustGetFullPath(ot *testing.T) {
	wd := MustGetWd()
	udwTest.Equal(MustGetFullPath("/tmp"), "/tmp")
	udwTest.Equal(MustGetFullPath("tmp"), filepath.Join(wd, "tmp"))
	udwTest.Equal(MustGetFullPath("."), wd)
	udwTest.Equal(MustGetFullPath(".."), filepath.Dir(wd))

}
