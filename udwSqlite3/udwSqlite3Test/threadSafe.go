package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwTask"
	"github.com/tachyon-protocol/udw/udwTest"
)

func TestThreadSafe() {
	const testDb = "/tmp/test_sqlite3.db"
	udwSqlite3.DeleteSqliteDbFileByPath(testDb)
	defer udwSqlite3.DeleteSqliteDbFileByPath(testDb)
	db1 := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:                       testDb,
		EmptyDatabaseIfDatabaseCorrupt: true,
	})
	defer db1.Close()
	db1.MustEmptyDatabase()
	db2 := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:                       testDb,
		EmptyDatabaseIfDatabaseCorrupt: true,
	})
	defer db2.Close()
	task := udwTask.New(10)
	for i := 0; i < 100; i++ {
		i := i
		task.AddFunc(func() {
			db1.MustSet("test", udwStrconv.FormatInt(i), "hello"+udwStrconv.FormatInt(i))
		})
		task.AddFunc(func() {
			db2.MustSet("test", udwStrconv.FormatInt(i), "hello"+udwStrconv.FormatInt(i))
		})
	}
	task.Close()
}

func testThreadSafe2() {
	task := udwTask.New(10)
	udwFile.MustDelete("/tmp/test_sqlite3.db")
	for i := 0; i < 100; i++ {
		i := 0
		task.AddFunc(func() {
			db1 := udwSqlite3.MustNewMemoryDb()
			defer db1.Close()
			k := udwStrconv.FormatInt(i)
			v := "hello" + udwStrconv.FormatInt(i)
			db1.MustSet("test", k, v)
			udwTest.Equal(db1.MustGet("test", k), v)
		})
	}
	task.Close()
}
