package udwTime

import (
	"sync"
	"time"
)

func NowFromDefaultNower() (out time.Time) {
	gDefaultNowerLocker.RLock()
	out = gDefaultNower()
	gDefaultNowerLocker.RUnlock()
	return out
}

func MysqlNowFromDefaultNower() string {
	return NowFromDefaultNower().Format(FormatMysql)
}

func MysqlUsNowFromDefaultNower() string {
	return MysqlUsFormat(NowFromDefaultNower())
}

func SetDefaultNowerCallback(cb func() time.Time) {
	gDefaultNowerLocker.Lock()
	gDefaultNower = cb
	gDefaultNowerLocker.Unlock()

}
func SetDefaultNowerToFixTimeString(s string) {
	SetDefaultNowerCallback(newFixedNower(MustParseAutoInDefault(s)))
}

func SetDefaultNowerToFixTime(t time.Time) {
	SetDefaultNowerCallback(func() time.Time {
		return t
	})
}

func SetDefaultNowerToRealTime() {
	SetDefaultNowerCallback(realTimeNow)
}
func SetDefaultNowerToOffset(offset time.Duration) {
	SetDefaultNowerCallback(func() time.Time {
		return time.Now().In(GetDefaultTimeZone()).Add(offset)
	})
}

type tDefaultNower func() time.Time

var gDefaultNowerLocker sync.RWMutex
var gDefaultNower tDefaultNower = realTimeNow

func newFixedNower(t time.Time) tDefaultNower {
	return func() time.Time {
		return t
	}
}

func realTimeNow() time.Time {
	return time.Now().In(GetDefaultTimeZone())
}
