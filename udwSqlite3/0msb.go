package udwSqlite3

import (
	"github.com/tachyon-protocol/udw/udwMap"
	"github.com/tachyon-protocol/udw/udwTask"
	"sync"
)

type Msb struct {
	db        *Db
	setMapMap map[string]map[string]string
	locker    sync.Mutex
	tasker    *udwTask.LimitMaxThreadTaskManager
}

func (msc *Msb) Set(k1 string, k2 string, v string) {
	msc.locker.Lock()
	if msc.setMapMap == nil {
		msc.setMapMap = map[string]map[string]string{}
	}
	m1 := msc.setMapMap[k1]
	if m1 == nil {
		m1 = map[string]string{}
		msc.setMapMap[k1] = m1
	}
	m1[k2] = v
	if len(m1) > maxMulitSetNum {
		delete(msc.setMapMap, k1)
		msc.locker.Unlock()
		msc.mulitSetMap(k1, m1)
		return
	}
	msc.locker.Unlock()
}

func (msc *Msb) MustFlush() {
	msc.locker.Lock()
	thisSetMM := msc.setMapMap
	if len(thisSetMM) > 0 {
		msc.setMapMap = nil
	}
	msc.locker.Unlock()
	for k1, m1 := range thisSetMM {
		msc.mulitSetMap(k1, m1)
	}
	thisSetMM = nil
	msc.tasker.WaitAndNotClose()
}

func (msc *Msb) MustClose() {
	msc.MustFlush()
	msc.tasker.Close()
}

func (msc *Msb) mulitSetMap(k1 string, m1 map[string]string) {
	msc.tasker.AddFunc(func() {
		msc.db.MustMulitSetMap(k1, m1)
	})
}

const maxMulitSetNum = 10000

func (db *Db) NewMsb(threadNum int) *Msb {
	return &Msb{
		db:     db,
		tasker: udwTask.New(threadNum),
	}
}

func (db *Db) MustMulitSetMap(k1 string, m map[string]string) {
	pairList := []udwMap.KeyValuePair{}
	for k, v := range m {
		pairList = append(pairList, udwMap.KeyValuePair{
			Key:   k,
			Value: v,
		})
	}
	db.MustMulitSet(k1, pairList)
}
