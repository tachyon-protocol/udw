package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwTask"
	"github.com/tachyon-protocol/udw/udwTest"
	"strconv"
	"strings"
)

func TestNoSql() {
	MustRunTestDb(func(db *udwSqlite3.Db) {
		db.MustDeleteK1("test")
		db.MustDelete("test", "\x00")
		retS, err := db.Get("test", "2")
		udwTest.Equal(err, nil)
		udwTest.Equal(retS, "")

		db.MustDelete("test", "\x00")

		db.MustSet("test", "1", "2")
		udwTest.Equal(db.MustGet("test", "1"), "2")

		retS, err = db.Get("test", "2")
		udwTest.Equal(err, nil)
		udwTest.Equal(retS, "")

		db.MustSet("test", "1", "3")
		udwTest.Equal(db.MustGet("test", "1"), "3")

		db.MustSet("test", "\x00", "\x00\x00\x00")
		udwTest.Equal(db.MustGet("test", "\x00"), "\x00\x00\x00")

		db.MustDelete("test", "\x00")

		db.MustMulitSet("test", []udwSqlite3.KeyValuePair{
			{"1", "v1"},
			{"2", "v2"},
			{"3", "v3"},
		})
		udwTest.Equal(db.MustGet("test", "1"), "v1")
		udwTest.Equal(db.MustGet("test", "2"), "v2")
		udwTest.Equal(db.MustGet("test", "3"), "v3")

		tname := strings.Repeat("0", 100)
		db.MustSet(tname, "1", "2")
		udwTest.Equal(db.MustGet(tname, "1"), "2")
		db.MustDeleteK1(tname)
		udwTest.Equal(db.MustGet(tname, "1"), "")
		db.MustDeleteK1(tname)
		udwTest.Equal(db.MustGet(tname, "1"), "")
		tname = strings.Repeat("0", 1000)
		db.MustSet(tname, "1", "2")
		udwTest.Equal(db.MustGet(tname, "1"), "2")
		udwTest.Equal(db.MustExist(tname, "1"), true)

		db.MustSet(tname, "1", "")
		udwTest.Equal(db.MustExist(tname, "1"), true)
		db.MustDelete(tname, "1")
		udwTest.Equal(db.MustExist(tname, "1"), false)

		udwTest.Equal(db.MustExist("not_exist_k1", "1"), false)
	})
}

func TestInsert() {
	MustRunTestDb(func(db *udwSqlite3.Db) {
		db.MustEmptyDatabase()
		db.MustInsert("test1", "1", "1")
		udwTest.AssertPanicWithErrorMessage(func() {
			db.MustInsert("test1", "1", "1")
		}, "UNIQUE constraint")
	})

}

func MustRunTestDb(cb func(db *udwSqlite3.Db)) {
	const dbPath = "/tmp/test_sqlite3.db"
	udwFile.MustDelete(dbPath)
	db := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath: "/tmp/test_sqlite3.db",
	})
	cb(db)
	db.Close()
	udwFile.MustDelete(dbPath)
}

func TestSet() {
	MustRunTestDb(func(db *udwSqlite3.Db) {
		db.MustEmptyDatabase()
		task := udwTask.NewLimitThreadTaskManager(10)
		for i := 0; i < 100; i++ {
			i := i
			task.AddFunc(func() {
				db.MustSet("test", udwStrconv.FormatInt(i), "hello"+udwStrconv.FormatInt(i))
			})
		}
		task.Wait()
		list := db.MustGetRange(udwSqlite3.MustGetRangeRequest{
			K1: "test",
		})
		udwTest.Equal(len(list), 100, list)
		udwTest.Equal(len(db.MustGetAllK1()), 1)
	})
	db := udwSqlite3.MustNewMemoryDb()
	defer db.Close()
	db.MustSet("test", "k1", "")
	udwTest.Equal(db.MustGet("test", "k1"), "")
	udwTest.Equal(db.MustExist("test", "k1"), true)
}

func TestMustInsertAndReturnExist() {
	MustRunTestDb(func(db *udwSqlite3.Db) {
		db.MustEmptyDatabase()
		ret := db.MustInsertAndReturnHasSucc("test", "1", "1")
		udwTest.Equal(ret, true)
		ret = db.MustInsertAndReturnHasSucc("test", "1", "1")
		udwTest.Equal(ret, false)
		ret = db.MustInsertAndReturnHasSucc("test", "1", "1")
		udwTest.Equal(ret, false)
	})
}

func TestMustMulitDelete() {
	const dbPath = "/tmp/test_sqlite3.db"
	udwFile.MustDelete(dbPath)
	db := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:           "/tmp/test_sqlite3.db",
		MulitDeleteMaxSize: 10,
	})
	db.MustEmptyDatabase()
	db.MustSet("test", "1", "1")
	db.MustSet("test", "2", "1")
	db.MustMulitDelete("test", []string{"1", "2"})
	udwTest.Equal(db.MustGet("test", "1"), "")
	udwTest.Equal(db.MustGet("test", "2"), "")
	db.MustSet("test", "3", "1")
	db.MustMulitDelete("test", []string{"1", "2", "3"})
	udwTest.Equal(db.MustGet("test", "3"), "")
	db.MustMulitDelete("test", []string{})

	keyList := []string{}
	keyValueList := []udwSqlite3.KeyValuePair{}
	for i := 0; i < 100; i++ {
		k := strconv.Itoa(i)
		keyList = append(keyList, k)
		keyValueList = append(keyValueList, udwSqlite3.KeyValuePair{
			Key:   k,
			Value: k,
		})
	}
	db.MustMulitSet("test", keyValueList)

	db.MustMulitDelete("test", keyList)
	outList := db.MustGetRange(udwSqlite3.MustGetRangeRequest{
		K1: "test",
	})
	udwTest.Equal(len(outList), 0)

	db.Close()
	udwFile.MustDelete(dbPath)
}

func testMustUpdate() {
	db := udwSqlite3.MustNewMemoryDb()
	defer db.Close()
	db.MustSet("test", "k1", "v1")
	db.MustUpdate("test", "k1", "v2")
	udwTest.Equal(db.MustGet("test", "k1"), "v2")
	db.MustUpdate("test", "k1", "")
	udwTest.Equal(db.MustGet("test", "k1"), "")

	db.MustUpdate("test", "k2", "v2")
	udwTest.Equal(db.MustGet("test", "k2"), "")
	db.MustUpdate("test2", "k1", "v1")
	udwTest.Equal(db.MustGet("test2", "k1"), "")
}
