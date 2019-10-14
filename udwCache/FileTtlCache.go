package udwCache

import (
	"encoding/hex"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwProjectPath"
	"github.com/tachyon-protocol/udw/udwTime"
	"path/filepath"
	"time"
)

type ttlCacheEntryV2 struct {
	Value   []byte
	Timeout time.Time
}

func getFileTtlCachePath(key string) string {
	return filepath.Join(udwProjectPath.MustGetProjectPath(), "tmp/FileTtlCache2", hex.EncodeToString([]byte(key)))
}

func FileTtlCache(key string, f func() (b []byte, ttl time.Duration, err error)) (b []byte, err error) {
	entry := ttlCacheEntryV2{}
	cacheFilePath := getFileTtlCachePath(key)
	now := time.Now()
	func() {
		content, err := udwFile.ReadFile(cacheFilePath)
		if err != nil {
			return
		}
		m, ok := unmarshalTtlCacheEntryV2(content)
		if ok == false {
			return
		}
		entry = m
	}()
	if entry.Timeout.After(now) {
		return entry.Value, nil
	}
	b, ttl, err := f()
	if err != nil {
		return nil, err
	}
	entry.Value = b
	entry.Timeout = now.Add(ttl)
	err = udwFile.MkdirForFile(cacheFilePath)
	if err != nil {
		return nil, err
	}
	err = udwFile.WriteFile(cacheFilePath, marshalTtlCacheEntryV2(entry))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func marshalTtlCacheEntryV2(m ttlCacheEntryV2) []byte {
	_buf := udwBytes.BufWriter{}
	_buf.WriteUvarint(uint64(len(m.Value)))
	_buf.Write_(m.Value)
	uTime := udwTime.DbTimeGetUint64FromObjOrMax(m.Timeout)
	_buf.WriteLittleEndUint64(uTime)
	return _buf.GetBytes()
}

func unmarshalTtlCacheEntryV2(b []byte) (m ttlCacheEntryV2, ok bool) {
	_buf := udwBytes.NewBufReaderWithOk(b)
	l := int(_buf.ReadUvarint())
	m.Value = _buf.ReadSliceBySize(l)
	uTime := _buf.ReadLittleEndUint64()
	m.Timeout = udwTime.DbTimeGetObjFromUint64(uTime)
	if _buf.IsOk() == false {
		return m, false
	}
	return m, true
}
