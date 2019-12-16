package tyVpnClient

import (
	"testing"
	"github.com/tachyon-protocol/udw/udwTest"
)

func TestA(ot *testing.T){
	u:=MarshalConfig(Config{
		ServerIp: "127.0.0.1",
		ServerTKey: "abc",
		ServerChk: "WyScbcpH-0JGk8m-BuZAcStIBseswEEklnuGr2PNEJDE",
	})
	udwTest.Equal(u,"ty://?ip=127.0.0.1&t=abc&chk=WyScbcpH-0JGk8m-BuZAcStIBseswEEklnuGr2PNEJDE")
	config2,errMsg:=ParseAndVerifyConfigS(u)
	udwTest.Equal(errMsg,"")
	udwTest.Equal(config2.ServerIp,"127.0.0.1")
	udwTest.Equal(config2.ServerTKey,"abc")
	udwTest.Equal(config2.ServerChk,"WyScbcpH-0JGk8m-BuZAcStIBseswEEklnuGr2PNEJDE")
}