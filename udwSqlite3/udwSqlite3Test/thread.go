package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTask"
	"github.com/tachyon-protocol/udw/udwTest"
	"strconv"
	"strings"
)

func TestThread() {

	const testDb = "/tmp/test_sqlite3.db"
	udwSqlite3.DeleteSqliteDbFileByPath(testDb)
	defer udwSqlite3.DeleteSqliteDbFileByPath(testDb)
	db := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath: testDb,
	})
	const num = 1e2
	task := udwTask.New(10)
	for i := 0; i < num; i++ {
		i := i
		task.AddFunc(func() {
			iS := strconv.Itoa(i)
			msg := strings.Repeat(iS, 1000)
			db.MustSet("t1", iS, msg)
			out := db.MustGet("t1", iS)
			udwTest.Equal(out, msg)
			db.MustDelete("t1", iS)
		})
	}
	task.Close()
	task = udwTask.New(10)
	for i := 0; i < num; i++ {
		i := i
		task.AddFunc(func() {
			iS := strconv.Itoa(i)
			msg := "hello" + iS
			db.MustSet("t1"+iS, iS, msg)
			out := db.MustGet("t1"+iS, iS)
			udwTest.Equal(out, msg)
			db.MustDelete("t1"+iS, msg)
		})
	}
	task.Close()
	db.Close()
	TestThread2()
}

func TestThread2() {

	const testDb = "/tmp/test_sqlite3.db"
	udwSqlite3.DeleteSqliteDbFileByPath(testDb)
	defer udwSqlite3.DeleteSqliteDbFileByPath(testDb)
	db := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath: testDb,
	})
	msb := db.NewMsb(3)
	const num = 1e2
	for i := 0; i < num; i++ {
		iS := strconv.Itoa(i)
		msg := strings.Repeat(iS, 1000)
		msb.Set("t1", "a"+iS, msg)
	}
	msb.MustClose()

	db.MustGetRangeCallback(udwSqlite3.GetRangeReq{
		K1: "t1",
	}, func(k string, v string) {
		db.MustSet("t2", k, v)
		return
	})

	db.MustGetRangeCallback(udwSqlite3.GetRangeReq{
		K1:     "t1",
		Prefix: "a",
	}, func(k string, v string) {
		db.MustSet("t1", "b"+k, v)
		return
	})

}
