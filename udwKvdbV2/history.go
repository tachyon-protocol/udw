package udwKvdbV2

import (
	"github.com/tachyon-protocol/udw/udwBytes/udwBytesEncode"
)

type KvdbWriteHistory struct {
	IsSnapshot bool
	SetMap     map[string]string
}

func (this *KvdbWriteHistory) recordCmdSet(db *KvdbV2, k string, v string) {
	if this.IsSnapshot {
		return
	}
	if this.SetMap == nil {
		this.SetMap = map[string]string{}
	}
	this.SetMap[k] = v
	if len(db.data) <= this.getLength() {
		this.switchToSnapshot()
	}
}

func (this *KvdbWriteHistory) switchToSnapshot() {
	this.IsSnapshot = true
	this.SetMap = nil
}

func (this *KvdbWriteHistory) toZero() {
	this.IsSnapshot = false
	this.SetMap = nil
}

func (obj KvdbWriteHistory) isZero() bool {
	return !obj.IsSnapshot && len(obj.SetMap) == 0
}

func (obj KvdbWriteHistory) getLength() int {
	return len(obj.SetMap)
}

func marshalMap(SetMap map[string]string) []byte {
	return udwBytesEncode.MapStringStringMarshal(SetMap)
}

func unmarshalMap(data []byte) (SetMap map[string]string, ok bool) {
	return udwBytesEncode.MapStringStringUnmarshal(data)
}
