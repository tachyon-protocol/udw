package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwDbI"
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTest"
)

func testZeroBug() {
	db := udwSqlite3.MustNewMemoryDb()
	defer db.Close()
	db.MustSet("k1test", "00000", "00000")
	udwTest.Equal(db.MustGet("k1test", "00000"), "00000")
	udwTest.Equal(db.MustGet("k1test", "0000"), "")
	udwTest.Equal(db.MustGet("k1test", "0"), "")
	udwTest.Equal(db.MustGet("k1test", "1"), "")
	num := 0
	db.IMustGetRangeCallback(udwDbI.GetRangeReq{
		K1: `k1test`,
	}, func(k string, v string) {
		udwTest.Equal(k, "00000")
		udwTest.Equal(v, "00000")
		num++
	})
	udwTest.Equal(num, 1)

	db.MustEmptyK1("k1test")
	db.MustSet("k1test", "0", "0")
	udwTest.Equal(db.MustGet("k1test", "0"), "0")
	udwTest.Equal(db.MustGet("k1test", "00000"), "")
}
