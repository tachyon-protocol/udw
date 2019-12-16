package tyVpnClient

import (
	"sync"
	"github.com/tachyon-protocol/udw/udwQueryOnlyUrl"
	"github.com/tachyon-protocol/udw/tyTls"
)

func SetConfig(config Config){
	gConfigLocker.Lock()
	gConfig = config
	gConfigLocker.Unlock()
}

var gConfig Config
var gConfigLocker sync.RWMutex

func getConfig() Config {
	gConfigLocker.RLock()
	thisConfig:=gConfig
	gConfigLocker.RUnlock()
	return thisConfig
}

func ParseAndVerifyConfigS(configS string) (config Config,errMsg string){
	obj:=udwQueryOnlyUrl.ParseQueryUrlObj(configS)
	if obj==nil{
		return config,"5e7mgtu5rh"
	}
	if obj.ProtocolName!="ty"{
		return config,"2ssf326m86"
	}
	ip:=obj.GetFirstValueByKey("ip")
	if ip!=""{
		config.ServerIp = ip
	}
	chk:=obj.GetFirstValueByKey("chk")
	if chk!=""{
		config.ServerChk = chk
		if tyTls.IsChkValid(chk)==false{
			return config,"w8j8y5jbvr"
		}
	}
	t:=obj.GetFirstValueByKey("t")
	if t!=""{
		config.ServerTKey = t
	}
	return config,""
}

func MarshalConfig(config Config) (configS string){
	obj:=udwQueryOnlyUrl.QueryUrlObj{
		ProtocolName: "ty",
	}
	if config.ServerIp!=""{
		obj.AddKv("ip",config.ServerIp)
	}
	if config.ServerTKey!=""{
		obj.AddKv("t",config.ServerTKey)
	}
	if config.ServerChk!=""{
		obj.AddKv("chk",config.ServerChk)
	}
	return obj.Marshal()
}