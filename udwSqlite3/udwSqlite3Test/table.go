package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTest"
	"strings"
)

func TestMustGetAllK1() {
	db := udwSqlite3.MustNewMemoryDb()
	for i := 0; i < 256; i++ {

		db.MustEmptyDatabase()

		k1 := string([]byte{byte(i)})
		db.MustSet(k1, "1", "1")
		v := db.MustGet(k1, "1")
		udwTest.Equal(v, "1")
		k1List := db.MustGetAllK1()
		udwTest.Equal(len(k1List), 1)
		udwTest.Equal(k1List[0], k1)

	}
	MustRunTestDb(func(db *udwSqlite3.Db) {
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

func testGetAllDataInTableToRowMap() {
	db := udwSqlite3.MustNewMemoryDb()
	db.MustQuery(udwSqlite3.QueryReq{
		Query: `CREATE TABLE IF NOT EXISTS abc (
k VARBINARY(255) Primary Key,
v LONGBLOB NOT NULL
) WITHOUT ROWID`,
	})
	db.MustQuery(udwSqlite3.QueryReq{
		Query: `REPLACE INTO abc (k, v) VALUES (?, ?)`,
		Args: [][]byte{
			[]byte("k1"),
			[]byte("v1"),
		},
	})
	m := db.MustGetAllDataInTableToRowMap("abc")
	udwTest.Equal(len(m), 1)
	udwTest.Equal(m[0]["k"], "k1")
	udwTest.Equal(m[0]["v"], "v1")
}
