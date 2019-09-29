package udwErr

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestErrMsgList(ot *testing.T) {
	errmsgList := ErrmsgList{}
	udwTest.Equal(errmsgList.GetErrMsg(), "")
	errmsgList.AddErrMsg("test")
	udwTest.Equal(errmsgList.GetErrMsg(), "test")
	errmsgList.AddErrMsg("test2")
	udwTest.Equal(errmsgList.GetErrMsg(), "test\ntest2")
}
