package udwKvdbV2

import (
	"github.com/tachyon-protocol/udw/udwBytes/udwBytesEncode"
	"github.com/tachyon-protocol/udw/udwCryptoAesCtrV7"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwLog"
	"github.com/tachyon-protocol/udw/udwRand/udwRandNewId"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type KvdbV2 struct {
	locker          sync.Mutex
	data            map[string]string
	history         KvdbWriteHistory
	wgAutoFlush     *sync.WaitGroup
	flushLocker     sync.Mutex
	writeOffset     int64
	incCount        int
	dataFileDir     string
	dataFileAbsPath string
	psk             *[32]byte
	outputLogFn     func(s string)
	loadOnce        sync.Once
}

type NewKvdbV2Request struct {
	DataDirPath string
	PskString   string
	OutputLogFn func(s string)
}

func NewKvdbV2(req NewKvdbV2Request) *KvdbV2 {
	req.DataDirPath = filepath.Clean(req.DataDirPath)
	if !filepath.IsAbs(req.DataDirPath) {
		req.DataDirPath = filepath.Join(udwFile.MustGetWd(), req.DataDirPath)
	}
	this := &KvdbV2{
		data:            map[string]string{},
		dataFileDir:     req.DataDirPath,
		dataFileAbsPath: filepath.Join(req.DataDirPath, `data`),
	}
	if req.PskString != `` {
		this.psk = udwCryptoAesCtrV7.Get32PskSha3FromString(req.PskString)
	}
	if req.OutputLogFn == nil {
		this.outputLogFn = func(s string) {
			udwLog.Log(`jvc3zu929u`, s)
		}
	}
	go this.loadFromDisk()

	return this
}

func (this *KvdbV2) Get(k string) (v string) {
	this.loadFromDisk()
	this.locker.Lock()
	v = this.data[k]
	this.locker.Unlock()
	return v
}

func (this *KvdbV2) Set(k string, v string) {
	this.loadFromDisk()
	this.locker.Lock()
	oldv, ok := this.data[k]
	if v == `` {
		if ok {
			delete(this.data, k)
			this.history.recordCmdSet(this, k, v)
			this.noticeNewWrite()
		}
	} else if !ok || oldv != v {
		this.data[k] = v
		this.history.recordCmdSet(this, k, v)
		this.noticeNewWrite()
	}
	this.locker.Unlock()
}

func (this *KvdbV2) Delete(k string) {
	this.Set(k, ``)
}

func (this *KvdbV2) Clear() {
	this.loadFromDisk()
	this.locker.Lock()
	if len(this.data) > 0 {
		this.data = map[string]string{}
		this.history.switchToSnapshot()
		this.noticeNewWrite()
	}
	this.locker.Unlock()
}

func (this *KvdbV2) CloneInnerMap() (m map[string]string) {
	this.loadFromDisk()
	this.locker.Lock()
	m = make(map[string]string, len(this.data))
	for k, v := range this.data {
		m[k] = v
	}
	this.locker.Unlock()
	return m
}

func (this *KvdbV2) GetLength() int {
	this.loadFromDisk()
	this.locker.Lock()
	length := len(this.data)
	this.locker.Unlock()
	return length
}

func (this *KvdbV2) moveHistory() (history KvdbWriteHistory) {
	this.locker.Lock()
	if this.history.IsSnapshot || len(this.data) <= this.incCount {
		history.IsSnapshot = true
		if len(this.data) > 0 {
			history.SetMap = make(map[string]string, len(this.data))
			for k, v := range this.data {
				history.SetMap[k] = v
			}
		}
	} else {
		history = this.history
	}
	this.history.toZero()
	this.locker.Unlock()
	return history
}

func (this *KvdbV2) rollbackHistory() {
	this.locker.Lock()
	this.history.switchToSnapshot()
	this.locker.Unlock()
}

func (this *KvdbV2) flushAndWaitDisk_NOLOCK() (err error, isEmpty bool) {
	history := this.moveHistory()
	if history.isZero() {
		return nil, true
	}
	length, err := this.writeHistoryToDisk(history)
	if err != nil {
		this.rollbackHistory()
		return err, false
	}

	if history.IsSnapshot {
		this.writeOffset = length
		this.incCount = 0
	} else {
		this.writeOffset += length
		this.incCount += history.getLength()
	}
	this.locker.Lock()
	isEmpty = this.history.isZero()
	this.locker.Unlock()
	return nil, isEmpty
}

func (this *KvdbV2) openFileWithMkdir(name string, flag int) (f *os.File, err error) {
	f, err = os.OpenFile(name, flag, 0777)
	if err == nil {
		return f, nil
	}
	err = os.MkdirAll(filepath.Dir(name), 0777)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(name, flag, 0777)
}

func (this *KvdbV2) writeHistoryToDisk(history KvdbWriteHistory) (length int64, err error) {
	var f *os.File
	var tmpName string
	if history.IsSnapshot {
		tmpName = this.dataFileAbsPath + `.` + udwRandNewId.NewId()
		f, err = this.openFileWithMkdir(tmpName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC)
	} else {
		f, err = this.openFileWithMkdir(this.dataFileAbsPath, os.O_WRONLY|os.O_CREATE|os.O_SYNC)
	}
	if err != nil {
		return 0, err
	}
	if !history.IsSnapshot {
		_, err = f.Seek(this.writeOffset, io.SeekStart)
		if err != nil {
			f.Close()
			return 0, err
		}
	}
	data := udwBytesEncode.MapStringStringMarshal(history.SetMap)
	data = cryptoBlock(this.psk, data)
	_, err = f.Write(data)
	if err != nil {
		f.Close()
		return 0, err
	}
	err = udwFile.FileSync(f)
	if err != nil {
		f.Close()
		return 0, err
	}
	err = f.Close()
	if err != nil {
		return 0, err
	}
	if history.IsSnapshot {
		err = os.Rename(tmpName, this.dataFileAbsPath)
		if err != nil {
			return 0, err
		}
	}

	return int64(len(data)), nil
}

func (this *KvdbV2) loadFromDisk() {
	this.loadOnce.Do(func() {
		bs, err := ioutil.ReadFile(this.dataFileAbsPath)
		if err != nil {
			if !os.IsNotExist(err) {
				this.outputLogFn(`wpn5bt5uqv load1: ` + err.Error())
			}
			return
		}
		list, err := decryptoBlockList(this.psk, bs)
		if err != nil {
			this.outputLogFn(`z9js29q5yq load2: ` + err.Error())
		}
		this.loadFromList(list)

		fileList, err := udwFile.ReadDirFileOneLevel(this.dataFileDir)
		if err != nil {
			this.outputLogFn(`v4a3a4gmbn ` + err.Error())
		}
		for _, fileName := range fileList {
			if strings.HasPrefix(fileName, `data.`) {
				err = os.Remove(fileName)
				if err != nil {
					this.outputLogFn(`m2n4ewk82d ` + fileName + `:` + err.Error())
				}
			}
		}
	})
}

func (this *KvdbV2) loadFromList(list []decryptItem) {
	for idx, one := range list {
		SetMap, ok := unmarshalMap(one.data)
		if !ok {
			this.outputLogFn(`6r7dxt4hrc load err`)
			break
		}
		for k, v := range SetMap {
			if v == `` {
				delete(this.data, k)
			} else {
				this.data[k] = v
			}
		}
		this.writeOffset += int64(one.lengthInFile)
		if idx > 0 {
			this.incCount += len(SetMap)
		}
	}
}

func (this *KvdbV2) autoFlush_NOLOCK() {
	for {
		err, isEmpty := this.flushAndWaitDisk_NOLOCK()
		if err != nil {
			this.outputLogFn(`4b9k4qn5xb autoFlush error: ` + err.Error())
			break
		}
		if isEmpty {
			break
		}
	}
}

func (this *KvdbV2) noticeNewWrite() {
	if this.wgAutoFlush == nil {
		this.wgAutoFlush = &sync.WaitGroup{}
		this.wgAutoFlush.Add(1)
		go func() {
			this.flushLocker.Lock()
			this.autoFlush_NOLOCK()
			this.flushLocker.Unlock()

			this.locker.Lock()
			this.wgAutoFlush.Done()
			this.wgAutoFlush = nil
			this.locker.Unlock()
		}()
	}
}

func (this *KvdbV2) FlushAndWaitDiskIO() (err error) {
	this.loadFromDisk()

	this.locker.Lock()
	wg := this.wgAutoFlush
	this.locker.Unlock()
	if wg != nil {
		wg.Wait()
	}

	this.flushLocker.Lock()
	err, _ = this.flushAndWaitDisk_NOLOCK()
	this.flushLocker.Unlock()
	return err
}
