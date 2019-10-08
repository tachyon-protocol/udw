package udwTime

import (
	"github.com/tachyon-protocol/udw/udwStrconv"
	"sync"
	"time"
)

func GetDefaultTimeZone() *time.Location {
	return GetUtc8Zone()
}

func GetUtcD8Zone() *time.Location {
	return time.FixedZone("UTC-8", -8*60*60)
}

func GetUtcD7Zone() *time.Location {
	return time.FixedZone("UTC-7", -7*60*60)
}

func GetUtcD5Zone() *time.Location {
	return time.FixedZone("UTC-5", -5*60*60)
}
func GetUtcD4Zone() *time.Location {
	return time.FixedZone(`UTC-4`, -4*60*60)
}
func GetUtc2Zone() *time.Location {
	return time.FixedZone("UTC2", 2*60*60)
}

var gUtc8Zone *time.Location
var gUtc8ZoneOnce sync.Once

func GetUtc8Zone() *time.Location {
	gUtc8ZoneOnce.Do(func() {
		gUtc8Zone = time.FixedZone("UTC8", 8*60*60)
	})
	return gUtc8Zone
}

func FormatTimeZone(tz *time.Location) string {
	if tz == time.UTC {
		return "UTC"
	}
	name, offset := time.Now().In(tz).Zone()
	s := name + `(`
	if offset == 0 {
		return s + "00:00)"
	} else if offset > 0 {
		s += "+"
	} else {
		offset = -offset
		s += "-"
	}
	tzHour := offset / 3600
	tzMin := (offset - tzHour*3600) / 60
	s += udwStrconv.FormatIntPaddingWithZeroPre(tzHour, 2) + ":" + udwStrconv.FormatIntPaddingWithZeroPre(tzMin, 2) + ")"
	return s
}

func GetFixedZoneByTzSecond(locOffset int) *time.Location {
	if locOffset == 0 {
		return time.UTC
	}
	locationCacheMapLocker.Lock()
	if locationCacheMap == nil {
		locationCacheMap = map[int]*time.Location{}
	}
	timeZone := locationCacheMap[locOffset]
	if timeZone == nil {
		timeZone = time.FixedZone("", locOffset)
	}
	if len(locationCacheMap) >= 1024 {
		for k := range locationCacheMap {
			delete(locationCacheMap, k)
		}
	}
	locationCacheMap[locOffset] = timeZone
	locationCacheMapLocker.Unlock()
	return timeZone
}

var locationCacheMapLocker sync.Mutex
var locationCacheMap map[int]*time.Location
