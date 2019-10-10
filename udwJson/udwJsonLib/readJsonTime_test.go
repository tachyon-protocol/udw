package udwJsonLib_test

import (
	"github.com/tachyon-protocol/udw/udwJson"
	"github.com/tachyon-protocol/udw/udwJson/udwJsonLib"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
	"time"
)

func TestReadJsonTime(ot *testing.T) {
	parseString := func(s string) time.Time {
		ctx := udwJsonLib.NewContextFromBuffer([]byte(s))
		return udwJsonLib.ReadJsonTime(ctx)
	}
	for _, cas := range []struct {
		t        time.Time
		inString string
	}{
		{time.Date(1994, 9, 17, 20, 4, 26, 0, time.FixedZone("", -18000)), `"1994-09-17T20:04:26-05:00"`},
		{time.Date(2017, 6, 12, 01, 53, 04, 270407360, time.UTC), `"2017-06-12T01:53:04.27040736Z"`},
		{time.Date(2017, 6, 12, 01, 53, 04, 270407300, time.UTC), `"2017-06-12T01:53:04.2704073Z"`},
		{time.Date(2017, 6, 12, 01, 53, 04, 270407000, time.UTC), `"2017-06-12T01:53:04.270407Z"`},
		{time.Date(2017, 6, 12, 01, 53, 04, 270400000, time.UTC), `"2017-06-12T01:53:04.27040Z"`},
		{time.Date(2017, 6, 12, 01, 53, 04, 270400000, time.UTC), `"2017-06-12T01:53:04.2704Z"`},
		{time.Date(2017, 6, 12, 01, 53, 04, 270000000, time.UTC), `"2017-06-12T01:53:04.270Z"`},
		{time.Date(2017, 6, 12, 01, 53, 04, 270000000, time.UTC), `"2017-06-12T01:53:04.27Z"`},
		{time.Date(2017, 6, 12, 01, 53, 04, 200000000, time.UTC), `"2017-06-12T01:53:04.20Z"`},
		{time.Date(2017, 6, 12, 01, 53, 04, 200000000, time.UTC), `"2017-06-12T01:53:04.2Z"`},
		{time.Date(2017, 6, 12, 01, 53, 04, 0, time.UTC), `"2017-06-12T01:53:04Z"`},

		{time.Date(2008, 9, 17, 20, 4, 26, 1, time.UTC), `"2008-09-17T20:04:26.000000001Z"`},
		{time.Date(2000, 12, 26, 1, 15, 6, 123456789, time.FixedZone("", 15600)), `"2000-12-26T01:15:06.123456789+04:20"`},
		{time.Date(2000, 12, 26, 1, 15, 6, 0, time.FixedZone("", 8*3600)), `"2000-12-26T01:15:06+08:00"`},
		{time.Time{}, `"0001-01-01T00:00:00Z"`},
		{time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC), `"0000-01-01T00:00:00Z"`},
		{time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC), `"9999-12-31T23:59:59Z"`},
		{time.Date(1994, 9, 17, 20, 4, 26, 0, time.FixedZone("", -18000)), `"1994-09-17T20:04:26-05:00","abc"`},
	} {
		t := parseString(cas.inString)
		udwTest.Equal(t, cas.t)
		outS := udwJson.MustMarshalToString(t)
		t2 := parseString(outS)
		udwTest.Equal(t2, cas.t)
	}
	for _, ts := range []string{
		`"1994-09-17T20:04:26Z1"`,
		`"2008-09-17T20:04:26.000000001Z1"`,
		`"2000-12-26T01:15:06.123456789+04:201"`,
		`"199:-09-17T20:04:26Z"`,
		`"1994-0:-17T20:04:26Z"`,
		`"1994-09-1:T20:04:26Z"`,
		`"1994-09-17T2::04:26Z"`,
		`"1994-09-17T20:0::26Z"`,
		`"1994-09-17T20:04:2:Z"`,
		`"1994-09-17T20:04:26F"`,
		`"1994-09-17T20:04:26.00000010:Z"`,
		`"1994-09-17T20:04:60Z"`,
		`"1994-09-17T20:60:26Z"`,
		`"1994-09-17T24:04:26Z"`,
		`"1994-09-32T20:04:26Z"`,
		`"1994-13-17T20:04:26Z"`,
		`"1994-00-17T20:04:26Z"`,
		`"1994-09-00T20:04:26Z"`,
		`"0000-00-00T00:00:00Z"`,
		`"1994-09-17T20:04:26-24:00"`,
		`"1994-09-17T20:04:26+24:00"`,
	} {
		udwTest.AssertPanic(func() {
			parseString(ts)
		})
	}

}
