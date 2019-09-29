package udwJsonLib_test

import (
	"github.com/tachyon-protocol/udw/udwJsonLib"
	"github.com/tachyon-protocol/udw/udwTest"
	"github.com/tachyon-protocol/udw/udwTime"
	"testing"
	"time"
)

func TestWriteJsonTime(ot *testing.T) {
	for _, cas := range []struct {
		t            time.Time
		shouldResult string
	}{
		{time.Date(2008, 9, 17, 20, 4, 26, 1, time.UTC), `"2008-09-17T20:04:26.000000001Z"`},
		{time.Date(1994, 9, 17, 20, 4, 26, 0, time.FixedZone("EST", -18000)), `"1994-09-17T20:04:26-05:00"`},
		{time.Date(2000, 12, 26, 1, 15, 6, 123456789, time.FixedZone("OTO", 15600)), `"2000-12-26T01:15:06.123456789+04:20"`},
		{time.Date(2000, 12, 26, 1, 15, 6, 0, udwTime.GetBeijingZone()), `"2000-12-26T01:15:06+08:00"`},
	} {
		ctx := &udwJsonLib.Context{}
		udwJsonLib.WriteJsonTime(ctx, cas.t)
		udwTest.Equal(string(ctx.WriterBytes()), cas.shouldResult)
		newT, err := time.Parse(`"`+time.RFC3339Nano+`"`, string(ctx.WriterBytes()))
		udwTest.Equal(err, nil)
		udwTest.Ok(newT.Equal(cas.t))
	}
}
