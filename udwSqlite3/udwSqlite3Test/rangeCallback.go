package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwSync"
	"github.com/tachyon-protocol/udw/udwTest"
)

func TestRangeCallback2() {
	db := udwSqlite3.MustNewMemoryDb()
	defer db.Close()
	num := udwSync.NewInt(0)
	db.MustGetRangeCallback(udwSqlite3.GetRangeReq{
		K1: "test",
	}, func(key string, value string) {
		num.Add(1)
		return
	})
	udwTest.Equal(num.Get(), 0)
	db.MustGetRangeKeyListCallback(udwSqlite3.GetRangeReq{
		K1: "test",
	}, func(key string) {
		num.Add(1)
		return
	})
	udwTest.Equal(num.Get(), 0)

	db.MustSet("test", "3", "v3")
	db.MustSet("test", "1", "v1")
	db.MustSet("test", "2", "v2")
	valueList := []string{}
	db.MustGetRangeCallback(udwSqlite3.GetRangeReq{
		K1: "test",
	}, func(key string, value string) {
		valueList = append(valueList, value)
		return
	})
	udwTest.Equal(valueList, []string{
		"v1", "v2", "v3",
	})
}
