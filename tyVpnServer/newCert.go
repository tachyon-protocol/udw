package tyVpnServer

import (
	"github.com/tachyon-protocol/udw/udwJson"
	"github.com/tachyon-protocol/udw/tyTls"
	"crypto/tls"
	"github.com/tachyon-protocol/udw/udwRand/udwRandNewId"
)

type serverStorageInfo struct{
	ServerTlsCert *tls.Certificate `json:"-"`
	ServerTlsCertPem string `json:",omitempty"`
	ServerTlsPkPem string `json:",omitempty"`
	ServerChk string `json:",omitempty"`
	SelfTKey string `json:",omitempty"`
}

func getServerStorageInfo() (info serverStorageInfo){
	udwJson.ReadFile("/usr/local/etc/tachyonServer.json",&info)
	hasChange:=false
	newCertFn:=func(){
		resp:=tyTls.MustNewTlsCert(false)
		info.ServerTlsCertPem = resp.CertPem
		info.ServerTlsPkPem = resp.PkPem
		hasChange = true
	}
	if info.ServerTlsPkPem=="" || info.ServerTlsCertPem==""{
		newCertFn()
	}
	if info.SelfTKey==""{
		info.SelfTKey = udwRandNewId.NewIdLen10()
		hasChange = true
	}
	var err error
	tlsCert,err := tls.X509KeyPair([]byte(info.ServerTlsCertPem),[]byte(info.ServerTlsPkPem))
	if err!=nil{
		newCertFn()
		tlsCert,err = tls.X509KeyPair([]byte(info.ServerTlsCertPem),[]byte(info.ServerTlsPkPem))
		if err!=nil{
			panic(err)
		}
	}
	info.ServerTlsCert = &tlsCert
	info.ServerChk = tyTls.MustHashChkFromTlsCert(info.ServerTlsCert)
	if hasChange{
		udwJson.WriteFile("/usr/local/etc/tachyonServer.json",info)
	}
	return info
}

