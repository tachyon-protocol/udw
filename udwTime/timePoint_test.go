package udwTime

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
	"time"
)

func TestGetTimePointInBetween(ot *testing.T) {
	type tCas struct {
		t1  string
		t2  string
		dur time.Duration
		num int
	}
	for _, cas := range []tCas{
		{"2000-01-01T01:00:00.000000000", "2000-01-01T01:05:00.000000000", time.Minute * 5, 1},
		{"2000-01-01T01:00:01.000000000", "2000-01-01T01:05:00.000000000", time.Minute * 5, 1},
		{"2000-01-01T01:00:01.000000000", "2000-01-01T01:04:00.000000000", time.Minute * 5, 0},
		{"2000-01-01T01:04:59.999999999", "2000-01-01T01:05:00.000000000", time.Minute * 5, 1},
		{"2000-01-01T01:00:00.000000000", "2000-01-01T02:04:00.000000000", time.Minute * 5, 12},
		{"2000-01-01T01:00:00.000000000", "2000-01-01T01:00:01.000000000", time.Second, 1},
		{"2000-01-01T01:00:02.000000000", "2000-01-01T01:00:03.000000000", time.Second, 1},
	} {
		ret := GetTimePointInBetween(GetTimePointInBetweenRequest{
			StartTime: MustDbTimeGetObjFromString(cas.t1),
			EndTime:   MustDbTimeGetObjFromString(cas.t2),
			Duration:  cas.dur,
		})
		udwTest.Equal(ret, cas.num, cas.t1, cas.dur)
	}
}
