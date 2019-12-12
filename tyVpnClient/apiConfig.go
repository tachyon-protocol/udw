package tyVpnClient

import (
	"sync"
	"github.com/tachyon-protocol/udw/udwJson"
	"github.com/tachyon-protocol/udw/tyTls"
)

func SetConfig(config *Config){
	gConfigLocker.Lock()
	gConfig = *config
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

func ParseAndVerifyConfigS(configS string) (config *Config,errMsg string){
	err:=udwJson.UnmarshalFromString(configS,&config)
	if err!=nil{
		return nil,"svzvkntygd "+err.Error()
	}
	if config.ServerIp==""{
		return nil,"vzm3basqtz"
	}
	if config.ServerChk==""{
		return nil,"jypzbufjf2"
	}
	if tyTls.IsChkValid(config.ServerChk)==false{
		return nil,"3kvmjwhcdy"
	}
	return config,""
}