package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTest"
)

func TestMustTableCopy() {
	const dbPath = "/tmp/test_sqlite3.db"
	udwFile.MustDelete(dbPath)
	defer udwFile.MustDelete(dbPath)
	db := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath: "/tmp/test_sqlite3.db",
	})
	db.MustEmptyDatabase()
	db.MustSet("test", "k1", "v1")
	db.MustSet("test", "k2", "v2")
	db.MustTableCopy("test", "test2")
	udwTest.Equal(db.MustGet("test", "k1"), "v1")
	udwTest.Equal(db.MustGet("test", "k2"), "v2")
	udwTest.Equal(db.MustGet("test2", "k1"), "v1")
	udwTest.Equal(db.MustGet("test2", "k2"), "v2")
	db.Close()
}
