package udwCache

import (
	"encoding/json"
	"github.com/tachyon-protocol/udw/udwCryptoMd5"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwProjectPath"
)

func MustKvdbGetBytes(s string) (b []byte) {
	key := udwCryptoMd5.Md5Hex([]byte(s))
	content, err := udwFile.ReadFile(udwProjectPath.MustPathInProject("tmp/kvdb/" + key))
	if err != nil {
		if udwFile.ErrorIsFileNotFound(err) {
			return nil
		}
		panic(err)
	}
	return content
}

func MustKvdbSetBytes(s string, b []byte) {
	key := udwCryptoMd5.Md5Hex([]byte(s))
	udwFile.MustWriteFileWithMkdir(udwProjectPath.MustPathInProject("tmp/kvdb/"+key), b)
}

func MustKvdbGet(s string, obj interface{}) bool {
	b := MustKvdbGetBytes(s)
	if b == nil {
		return false
	}
	err := json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
	return true
}

func MustKvdbSet(s string, obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	MustKvdbSetBytes(s, b)
}
