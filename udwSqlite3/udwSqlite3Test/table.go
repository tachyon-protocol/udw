package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwSqlite3/sqlite3"
	"github.com/tachyon-protocol/udw/udwTest"
	"strings"
)

func TestMustGetAllK1() {
	MustRunTestDb(func(db *udwSqlite3.Db) {
		db.MustEmptyDatabase()
		tableList := db.MustGetAllK1()
		udwTest.Equal(len(tableList), 0, tableList)
		db.MustSet("test1", "1", "1")
		db.MustSet("test2", "1", "1")
		db.MustSet("\x00", "1", "1")
		db.MustSet("a/a", "1", "1")
		db.MustSet("z"+strings.Repeat("\x00", 200), "1", "1")
		tableList = db.MustGetAllK1()
		udwTest.Equal(len(tableList), 5, tableList)
		udwTest.Equal(tableList[0], "\x00")
		udwTest.Equal(tableList[1], "a/a")
		udwTest.Equal(tableList[2], "test1")
		udwTest.Equal(tableList[3], "test2")
		udwTest.Equal(tableList[4], "z"+strings.Repeat("\x00", 200))
	})
}

func TestVersion() {
	MustRunTestDb(func(db *udwSqlite3.Db) {
		version := db.MustGetVersion()

		udwTest.Equal(version, "3.20.1")
	})
}

func testGetAllDataInTableToRowMap() {
	db := udwSqlite3.MustNewMemoryDb()
	db.MustQuery(sqlite3.QueryReq{
		Query: `CREATE TABLE IF NOT EXISTS abc (
k VARBINARY(255) Primary Key,
v LONGBLOB NOT NULL
) WITHOUT ROWID`,
	})
	db.MustQuery(sqlite3.QueryReq{
		Query: `REPLACE INTO abc (k, v) VALUES (?, ?)`,
		Args: [][]byte{
			[]byte("k1"),
			[]byte("v1"),
		},
	})
	m := db.GetAllDataInTableToRowMap("abc")
	udwTest.Equal(len(m), 1)
	udwTest.Equal(m[0]["k"], "k1")
	udwTest.Equal(m[0]["v"], "v1")
}
