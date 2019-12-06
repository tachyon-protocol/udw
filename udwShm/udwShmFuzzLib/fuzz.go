package udwShmFuzzLib

import (
	"github.com/tachyon-protocol/udw/udwShm"
	"github.com/tachyon-protocol/udw/udwShm/udwShmTest"
	"github.com/tachyon-protocol/udw/udwTest"
)

func Fuzz(data []byte) int {
	var msg udwShmTest.MsgV5
	errMsg := udwShm.Unmarshal(data, &msg)
	if errMsg != "" {
		return 0
	}
	data2, errMsg := udwShm.Marshal(msg)
	udwTest.Ok(errMsg == ``)

	var msg3 udwShmTest.MsgV5
	errMsg = udwShm.Unmarshal(data2, &msg3)
	udwTest.Ok(errMsg == ``, errMsg)
	udwShmTest.MsgV5Equal(msg, msg3)

	var v5Sub_simpe udwShmTest.MsgV5_Simple
	errMsg = udwShm.Unmarshal(data2, &v5Sub_simpe)
	udwShmTest.MsgV5Equal_Simple(msg3, v5Sub_simpe)

	{
		b2 := data2
		var v5Sub udwShmTest.MsgV5
		errMsg = udwShm.Unmarshal(b2, v5Sub)
		udwTest.Ok(errMsg != "")

		var v5Sub2 *udwShmTest.MsgV5
		errMsg = udwShm.Unmarshal(b2, v5Sub2)
		udwTest.Ok(errMsg != "")

		errMsg = udwShm.Unmarshal(b2, &v5Sub2)
		udwTest.Ok(errMsg == "")

		var v5Sub3 **udwShmTest.MsgV5
		errMsg = udwShm.Unmarshal(b2, v5Sub3)
		udwTest.Ok(errMsg != "")

		var v5Sub4 *udwShmTest.MsgV5 = &udwShmTest.MsgV5{}
		errMsg = udwShm.Unmarshal(b2, v5Sub4)
		udwTest.Ok(errMsg == "")
	}
	return 1
}
