package udwJsonLib

import (
	"strconv"
	"sync"
	"time"
)

func ReadJsonTime(ctx *Context) time.Time {
	pos := ctx.readerPos
for1:
	for {
		if pos >= len(ctx.readerData) {
			panic(`[ReadJsonTime] need a " but get EOF pos ` + strconv.Itoa(pos))
		}
		c := ctx.readerData[pos]
		pos++
		switch c {
		case ' ', '\t', '\n':
		case '"':
			break for1
		case 'n':
			ctx.readerPos = pos - 1
			MustReadJsonNull(ctx)
			return time.Time{}
		default:
			panic(`[ReadJsonTime] need a " but get [` + string(c) + `] pos ` + strconv.Itoa(pos))
		}
	}
	ctx.readerPos = pos

	inBuf := ctx.readerData[ctx.readerPos:]
	if len(inBuf) < 21 ||
		isDigest(inBuf[0]) == false ||
		isDigest(inBuf[1]) == false ||
		isDigest(inBuf[2]) == false ||
		isDigest(inBuf[3]) == false ||
		inBuf[4] != '-' ||
		isDigest(inBuf[5]) == false ||
		isDigest(inBuf[6]) == false ||
		inBuf[7] != '-' ||
		isDigest(inBuf[8]) == false ||
		isDigest(inBuf[9]) == false ||
		inBuf[10] != 'T' ||
		isDigest(inBuf[11]) == false ||
		isDigest(inBuf[12]) == false ||
		inBuf[13] != ':' ||
		isDigest(inBuf[14]) == false ||
		isDigest(inBuf[15]) == false ||
		inBuf[16] != ':' ||
		isDigest(inBuf[17]) == false ||
		isDigest(inBuf[18]) == false {
		panic("timeParseRFC3339Nano format error 1 " + ComfirmGolangJsonLibBugTag)
	}
	year := int(inBuf[0])*1000 + int(inBuf[1])*100 +
		int(inBuf[2])*10 + int(inBuf[3]) - 53328

	month := int(inBuf[5])*10 + int(inBuf[6]) - 528

	day := int(inBuf[8])*10 + int(inBuf[9]) - 528

	hour := int(inBuf[11])*10 + int(inBuf[12]) - 528

	minute := int(inBuf[14])*10 + int(inBuf[15]) - 528

	second := int(inBuf[17])*10 + int(inBuf[18]) - 528

	var ns int
	rPos := 19
	if inBuf[19] == '.' {
		i := 20
		ns = 0
		for {
			if i >= len(inBuf) || i >= 29 || isDigest(inBuf[i]) == false {
				break
			}
			ns = ns*10 + int(inBuf[i]-'0')
			i++
		}
		i = i - 1
		if i == 19 {
			targetS := ""
			if len(inBuf) < 29 {
				targetS = string(inBuf)
			} else {
				targetS = string(inBuf[:29])
			}
			panic("timeParseRFC3339Nano format error 2 " + targetS)
		}
		scaleDigits := 28 - i
		for i := 0; i < scaleDigits; i++ {
			ns *= 10
		}

		rPos = i + 1
	}
	var timeZone *time.Location
	if inBuf[rPos] == 'Z' {
		if len(inBuf) < rPos+2 ||
			inBuf[rPos+1] != '"' {
			panic("timeParseRFC3339Nano format error 3")
		}
		ctx.readerPos = ctx.readerPos + rPos + 1 + 1
		timeZone = time.UTC
	} else {
		if len(inBuf) < rPos+7 ||
			inBuf[rPos+6] != '"' ||
			inBuf[rPos+3] != ':' ||
			isDigest(inBuf[rPos+1]) == false ||
			isDigest(inBuf[rPos+2]) == false ||
			isDigest(inBuf[rPos+4]) == false ||
			isDigest(inBuf[rPos+5]) == false {
			panic("timeParseRFC3339Nano format error 4 " + ComfirmGolangJsonLibBugTag)
		}

		locHour := int(inBuf[rPos+1])*10 + int(inBuf[rPos+2]) - 528
		locMin := int(inBuf[rPos+4])*10 + int(inBuf[rPos+5]) - 528
		locOffset := locHour*3600 + locMin*60
		if inBuf[rPos] == '+' {
		} else if inBuf[rPos] == '-' {
			locOffset = -locOffset
		} else {
			panic("timeParseRFC3339Nano format error 5")
		}
		if locOffset >= 24*60*60 || locOffset <= -24*60*60 {
			panic("timeParseRFC3339Nano format error 6 " + ComfirmGolangJsonLibBugTag)
		}
		if locOffset == 0 {
			timeZone = time.UTC
		} else {
			locationCacheMapLocker.Lock()
			if locationCacheMap == nil {
				locationCacheMap = map[int]*time.Location{}
			}
			timeZone = locationCacheMap[locOffset]
			if timeZone == nil {
				timeZone = time.FixedZone("", locOffset)
			}
			locationCacheMap[locOffset] = timeZone
			locationCacheMapLocker.Unlock()
		}
		ctx.readerPos = ctx.readerPos + rPos + 6 + 1
	}

	if month > 12 || month == 0 || day > 31 || day == 0 || hour > 23 || minute > 59 || second > 59 {
		panic("timeParseRFC3339Nano format error 7")
	}
	return time.Date(year, time.Month(month), day, hour, minute, second, ns, timeZone)
}

func isDigest(b byte) bool {
	return b >= '0' && b <= '9'
}

var locationCacheMapLocker sync.Mutex
var locationCacheMap map[int]*time.Location
