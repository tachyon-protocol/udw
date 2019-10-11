package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwSync"
	"github.com/tachyon-protocol/udw/udwTest"
)

func TestMustGetRangeValueListCallback() {
	MustRunTestDb(func(db *udwSqlite3.Db) {
		db.MustEmptyDatabase()
		num := 0
		db.MustGetRangeValueListCallback(udwSqlite3.MustGetRangeRequest{
			K1: "test",
		}, func(v string) {
			num++
		})
		udwTest.Equal(num, 0)

		db.MustSet("test", "3", "v3")
		db.MustSet("test", "1", "v1")
		db.MustSet("test", "2", "v2")

		valueList := []string{}
		db.MustGetRangeValueListCallback(udwSqlite3.MustGetRangeRequest{
			K1: "test",
		}, func(v string) {
			valueList = append(valueList, v)
		})
		udwTest.Equal(valueList, []string{
			"v1", "v2", "v3",
		})
	})
	testRangeCallback2()
}

func testRangeCallback2() {
	db := udwSqlite3.MustNewMemoryDb()
	defer db.Close()
	num := udwSync.NewInt(0)
	db.MustGetRangeCallback(udwSqlite3.MustGetRangeRequest{
		K1: "test",
	}, func(key string, value string) bool {
		num.Add(1)
		return true
	})
	udwTest.Equal(num.Get(), 0)
	db.MustGetRangeKeyListCallback(udwSqlite3.MustGetRangeRequest{
		K1: "test",
	}, func(key string) bool {
		num.Add(1)
		return true
	})
	udwTest.Equal(num.Get(), 0)

	db.MustSet("test", "3", "v3")
	db.MustSet("test", "1", "v1")
	db.MustSet("test", "2", "v2")
	valueList := []string{}
	db.MustGetRangeCallback(udwSqlite3.MustGetRangeRequest{
		K1: "test",
	}, func(key string, value string) bool {
		valueList = append(valueList, value)
		return true
	})
	udwTest.Equal(valueList, []string{
		"v1", "v2", "v3",
	})
}
