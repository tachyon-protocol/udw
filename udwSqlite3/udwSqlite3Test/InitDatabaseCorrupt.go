package udwSqlite3Test

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwSync"
	"github.com/tachyon-protocol/udw/udwTest"
	"strings"
)

func TestInitDatabaseCorrupt() {
	const dbFilePath = "/tmp/test_sqlite3.db"
	udwFile.MustDelete(dbFilePath)
	defer udwFile.MustDelete(dbFilePath)

	udwFile.MustWriteFile(dbFilePath, []byte("abc"))
	err := udwErr.PanicToError(func() {
		db := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
			FilePath:     dbFilePath,
			UsingWALFull: true,
		})
		defer db.Close()
		db.MustSet("test", "1", "1")
	})
	errMsg := udwErr.ErrorToMsg(err)
	udwTest.Equal(udwSqlite3.IsErrorDatabaseCorrupt(errMsg), true, errMsg)

	udwFile.MustWriteFile(dbFilePath, []byte("abc"))
	db := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:                       dbFilePath,
		EmptyDatabaseIfDatabaseCorrupt: true,
		UsingWALFull:                   true,
	})
	db.MustSet("test", "1", "1")
	udwTest.Equal(db.MustGet("test", "1"), "1")
	db.Close()

	udwFile.MustWriteFile(dbFilePath, []byte("abc"))
	db = udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:                       dbFilePath,
		EmptyDatabaseIfDatabaseCorrupt: true,
		UsingWALFull:                   true,
	})
	udwTest.Equal(db.MustGet("test", "1"), "")
	db.Close()

	udwFile.MustWriteFile(dbFilePath, []byte("abc"))
	db = udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:                       dbFilePath,
		EmptyDatabaseIfDatabaseCorrupt: true,
		UsingWALFull:                   true,
	})
	udwTest.Equal(db.MustGetRange(udwSqlite3.MustGetRangeRequest{
		K1: "test",
	}), nil)
	db.Close()
}

func TestInitDatabaseCorrupt2() {
	const dbFilePath = "/tmp/test_sqlite3.db"
	udwFile.MustDelete(dbFilePath)
	defer udwFile.MustDelete(dbFilePath)
	db := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:                        dbFilePath,
		HighPerformanceModeWithDataLoss: true,
		UsingWALFull:                    true,
	})
	db.MustSet("test", "1", "1")
	db.Close()
	buf := udwFile.MustReadFile(dbFilePath)

	tryCorruptDatabase := func(i int) {

		if i >= len(buf) {
			fmt.Println("[tryCorruptDatabase] i too large", i, len(buf))
			return
		}
		buf2 := udwBytes.Clone(buf)
		buf2[i]++
		udwFile.WriteFile(dbFilePath, buf2)
		udwFile.MustDelete(dbFilePath + "-shm")
		udwFile.MustDelete(dbFilePath + "-wal")
		db := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
			FilePath:                       dbFilePath,
			EmptyDatabaseIfDatabaseCorrupt: true,
			UsingWALFull:                   true,
		})
		db.MustSet("test", "1", "1")
		udwTest.Equal(db.MustGet("test", "1"), "1")
		db.Close()
	}
	tryCorruptDatabase(0)
	tryCorruptDatabase(47)
	tryCorruptDatabase(52)
	tryCorruptDatabase(3999)
	tryCorruptDatabase(4014)
	tryCorruptDatabase(4030)

}

func TestInitDatabaseCorrupt3() {
	const dbFilePath = "/tmp/test_sqlite3.db"
	udwSqlite3.DeleteSqliteDbFileByPath(dbFilePath)
	defer udwSqlite3.DeleteSqliteDbFileByPath(dbFilePath)
	CallbackNum := udwSync.Int{}
	testDbFn := func() {
		db := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
			FilePath:                       dbFilePath,
			EmptyDatabaseIfDatabaseCorrupt: true,
			UsingWALFull:                   true,
			DatabaseCorruptCallback: func() {
				CallbackNum.Add(1)
			},
			EncryptPskString: "123",
		})
		defer db.Close()
		b := strings.Repeat("1", 1024*1024)
		db.MustSet("test", "1", b)
		udwTest.Equal(db.MustGet("test", "1"), b)
		db.MustSet("test", "2", b)
		db.MustDelete("test", "1")
	}

	testDbFn()
	udwTest.Equal(CallbackNum.Get(), 0)
	udwFile.MustWriteFileWithMkdir(dbFilePath, []byte("abc"))

	testDbFn()
	udwTest.Equal(CallbackNum.Get(), 1)

	testDbFn()
	udwTest.Equal(CallbackNum.Get(), 1)
}
