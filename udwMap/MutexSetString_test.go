package udwMap

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestMutexSetString(ot *testing.T) {
	ss := MutexSetString{}
	udwTest.Ok(ss.Has("1") == false)
	ss.Set("1")
	udwTest.Ok(ss.Len() == 1)
	udwTest.Ok(ss.Has("1") == true)
	ss.Clear()
	udwTest.Ok(ss.Has("1") == false)
	udwTest.Ok(ss.Len() == 0)
}
