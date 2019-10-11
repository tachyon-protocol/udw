package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTest"
	"strconv"
)

func TestEncrypt() {
	testPath := "/tmp/test_sqlite3.db.encrypt"
	udwFile.MustDelete(testPath)
	defer udwFile.MustDelete(testPath)
	db1 := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:         testPath,
		UsingWALFull:     true,
		EncryptPskString: "123",
	})
	updateDbFn := func(db *udwSqlite3.Db) {
		for i := 0; i < 10; i++ {
			s := strconv.Itoa(i)
			db.MustSet("test", s, s)
		}
	}
	updateDbFn(db1)
	db1.Close()

	db1 = udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:         testPath,
		UsingWALFull:     true,
		EncryptPskString: "123",
	})
	for i := 0; i < 10; i++ {
		s := strconv.Itoa(i)
		udwTest.Equal(db1.MustGet("test", s), s)
	}
	db1.Close()

	udwTest.AssertPanicWithErrorMessage(func() {
		udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
			FilePath:     testPath,
			UsingWALFull: true,
		})
	}, "file is not a database")

	db2 := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:                       testPath,
		UsingWALFull:                   true,
		EmptyDatabaseIfDatabaseCorrupt: true,
	})
	updateDbFn(db2)
	db2.Close()

	db1 = udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:     testPath,
		UsingWALFull: true,
	})
	for i := 0; i < 10; i++ {
		s := strconv.Itoa(i)
		udwTest.Equal(db1.MustGet("test", s), s)
	}
	db1.Close()

	db1 = udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:                       testPath,
		UsingWALFull:                   true,
		EncryptPskString:               "123",
		EmptyDatabaseIfDatabaseCorrupt: true,
	})
	updateDbFn(db1)
	for i := 0; i < 10; i++ {
		s := strconv.Itoa(i)
		udwTest.Equal(db1.MustGet("test", s), s)
	}
	db1.Close()

	db1 = udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:         testPath,
		UsingWALFull:     true,
		EncryptPskString: "123",
	})
	for i := 0; i < 10; i++ {
		s := strconv.Itoa(i)
		udwTest.Equal(db1.MustGet("test", s), s)
	}
	db1.Close()
}
