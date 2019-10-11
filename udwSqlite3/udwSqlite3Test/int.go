package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTest"
	"math"
)

func testInt() {
	db := udwSqlite3.MustNewMemoryDb()
	defer db.Close()
	v := db.MustGetInt64("test", "k1")
	udwTest.Equal(v, int64(0))
	db.MustSetInt64("test", "k1", 15)
	v = db.MustGetInt64("test", "k1")
	udwTest.Equal(v, int64(15))
	udwTest.Equal(db.MustGet("test", "k1"), "15")
	db.MustSet("test", "k2", "16")
	udwTest.Equal(db.MustGetInt64("test", "k2"), int64(16))
	for _, c := range []int64{
		0,
		-1000,
		math.MaxInt64,
		math.MinInt64,
	} {
		db.MustSetInt64("test", "k1", c)
		udwTest.Equal(db.MustGetInt64("test", "k1"), c)
	}
}
