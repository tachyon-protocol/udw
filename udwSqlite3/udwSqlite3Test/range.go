package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwMap"
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTest"
)

func TestMustGetRange() {
	MustRunTestDb(func(db *udwSqlite3.Db) {
		db.MustEmptyDatabase()

		ret := db.MustGetRange(udwSqlite3.GetRangeReq{
			K1: "test",
		})
		udwTest.Equal(len(ret), 0)

		db.MustSet("test", "3", "v3")
		db.MustSet("test", "1", "v1")
		db.MustSet("test", "2", "v2")

		ret = db.MustGetRange(udwSqlite3.GetRangeReq{
			K1: "test",
		})
		udwTest.Equal(len(ret), 3)
		udwTest.Equal(ret[0].Key, "1")
		udwTest.Equal(ret[0].Value, "v1")
		udwTest.Equal(ret[1].Key, "2")
		udwTest.Equal(ret[1].Value, "v2")
		udwTest.Equal(ret[2].Key, "3")
		udwTest.Equal(ret[2].Value, "v3")

		sList := db.MustGetRangeValueList(udwSqlite3.GetRangeReq{
			K1: "test",
		})
		udwTest.Equal(len(sList), 3)
		udwTest.Equal(sList, []string{"v1", "v2", "v3"})

		sList = db.MustGetRangeKeyList(udwSqlite3.GetRangeReq{
			K1: "test",
		})
		udwTest.Equal(sList, []string{"1", "2", "3"})

		keyValuePairList := db.MustMulitGet("test", []string{"1", "2"})
		udwTest.Equal(keyValuePairList, []udwMap.KeyValuePair{
			{Key: "1", Value: "v1"},
			{Key: "2", Value: "v2"},
		})

		keyValuePairList = db.MustMulitGet("test", []string{"1", "-1"})
		udwTest.Equal(keyValuePairList, []udwMap.KeyValuePair{
			{Key: "1", Value: "v1"},
		})

		keyValuePairList = db.MustMulitGet("test", []string{"0", "-1"})
		udwTest.Equal(keyValuePairList, nil)
	})
	TestGetRange2()
}

func TestGetRange2() {
	db := udwSqlite3.MustNewMemoryDb()
	defer db.Close()
	ret := db.MustGetRangeValueList(udwSqlite3.GetRangeReq{
		K1: "test",
	})
	udwTest.Equal(len(ret), 0)
	keyList := db.MustGetRangeKeyList(udwSqlite3.GetRangeReq{
		K1: "test",
	})
	udwTest.Equal(len(keyList), 0)
	db.MustMulitSet("test", []udwMap.KeyValuePair{
		{"k1", "v1"},
		{"k2", "v2"},
	})
	m := db.MustGetRangeKeyMap(udwSqlite3.GetRangeReq{
		K1: "test",
	})
	udwTest.Equal(len(m), 2)
	m2 := db.MustGetRangeToMap(udwSqlite3.GetRangeReq{
		K1: "test",
	})
	udwTest.Equal(len(m2), 2)

	db.MustMulitSet("test", []udwMap.KeyValuePair{
		{"k1", "v1"},
		{"k2", "v2"},
		{"k3", "v3"},
	})

	keyList = db.MustGetRangeKeyList(udwSqlite3.GetRangeReq{
		K1:       "test",
		MinValue: "k2",
	})
	udwTest.Equal(len(keyList), 2)
	udwTest.Equal(keyList[0], "k2")
	udwTest.Equal(keyList[1], "k3")

	keyList = db.MustGetRangeKeyList(udwSqlite3.GetRangeReq{
		K1:                 "test",
		MinValueNotInclude: "k2",
	})
	udwTest.Equal(len(keyList), 1)
	udwTest.Equal(keyList[0], "k3")

	keyList = db.MustGetRangeKeyList(udwSqlite3.GetRangeReq{
		K1:       "test",
		MaxValue: "k2",
	})
	udwTest.Equal(len(keyList), 2)
	udwTest.Equal(keyList[0], "k1")
	udwTest.Equal(keyList[1], "k2")

	keyList = db.MustGetRangeKeyList(udwSqlite3.GetRangeReq{
		K1:                 "test",
		MaxValueNotInclude: "k2",
	})
	udwTest.Equal(len(keyList), 1)
	udwTest.Equal(keyList[0], "k1")

	keyList = db.MustGetRangeKeyList(udwSqlite3.GetRangeReq{
		K1:     "test",
		Prefix: "k2",
	})
	udwTest.Equal(len(keyList), 1)
	udwTest.Equal(keyList[0], "k2")

	keyList = db.MustGetRangeKeyList(udwSqlite3.GetRangeReq{
		K1:          "test",
		IsDescOrder: true,
	})
	udwTest.Equal(len(keyList), 3)
	udwTest.Equal(keyList[0], "k3")
	udwTest.Equal(keyList[1], "k2")
	udwTest.Equal(keyList[2], "k1")
}
