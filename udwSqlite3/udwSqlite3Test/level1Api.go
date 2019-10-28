package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwSync"
	"github.com/tachyon-protocol/udw/udwTest"
)

func testSimpleQuery() {
	db := udwSqlite3.MustNewMemoryDb()

	defer db.Close()

	num := udwSync.NewInt(0)
	errMsg := db.Query(udwSqlite3.QueryReq{
		Query: "SELECT sqlite_version();",
		ColumnsCb: func(result [][]byte) {
			num.Add(1)
			udwTest.Equal(len(result), 1)
			udwTest.Equal(string(result[0]), "sqlite_version()")
		},
		RespDataCb: func(result [][]byte) {
			num.Add(1)
			udwTest.Equal(len(result), 1)
			udwTest.Equal(string(result[0]), currentSqlite3Version)
		},
	})
	udwErr.PanicIfErrorMsg(errMsg)
	udwTest.Equal(num.Get(), 2)

	num.Set(0)
	errMsg = db.Query(udwSqlite3.QueryReq{
		Query: `CREATE TABLE IF NOT EXISTS abc (
k VARBINARY(255) Primary Key,
v LONGBLOB NOT NULL
) WITHOUT ROWID`,
		RespDataCb: func(result [][]byte) {
			num.Add(1)
		},
	})
	udwErr.PanicIfErrorMsg(errMsg)
	udwTest.Equal(num.Get(), 0)
	errMsg = db.Query(udwSqlite3.QueryReq{
		Query: `REPLACE INTO abc (k, v) VALUES (?, ?)`,
		Args: [][]byte{
			[]byte("k1"),
			[]byte("v1"),
		},
	})
	udwErr.PanicIfErrorMsg(errMsg)

	num.Set(0)
	errMsg = db.Query(udwSqlite3.QueryReq{
		Query: "SELECT k,v FROM abc",
		RespDataCb: func(result [][]byte) {
			num.Add(1)
			udwTest.Equal(len(result), 2)
			udwTest.Equal(string(result[0]), "k1")
			udwTest.Equal(string(result[1]), "v1")
		},
	})
	udwErr.PanicIfErrorMsg(errMsg)
	udwTest.Equal(num.Get(), 1)

	errMsg = db.Query(udwSqlite3.QueryReq{
		Query: "SELECT k,v FROM abc",
	})
	udwErr.PanicIfErrorMsg(errMsg)
}

func testSimpleQueryAffectedRows() {
	db := udwSqlite3.MustNewMemoryDb()

	defer db.Close()
	errMsg := db.Query(udwSqlite3.QueryReq{
		Query: `CREATE TABLE IF NOT EXISTS abc (
k INTEGER PRIMARY KEY,
v INTEGER
);`,
	})
	udwErr.PanicIfErrorMsg(errMsg)
	num := udwSync.NewInt(0)
	errMsg = db.Query(udwSqlite3.QueryReq{
		Query: `INSERT INTO abc (v) VALUES (1)`,
		RespStatusCb: func(status udwSqlite3.QueryRespStatus) {
			num.Add(1)
			udwTest.Equal(status.LastInsertId, uint64(1))
			udwTest.Equal(status.AffectedRows, uint64(1))
		},
	})
	udwErr.PanicIfErrorMsg(errMsg)
	udwTest.Equal(num.Get(), 1)

	errMsg = db.Query(udwSqlite3.QueryReq{
		Query: `INSERT INTO abc (v) VALUES (1)`,
		RespStatusCb: func(status udwSqlite3.QueryRespStatus) {
			num.Add(1)
			udwTest.Equal(status.LastInsertId, uint64(2))
			udwTest.Equal(status.AffectedRows, uint64(1))
		},
	})
	udwErr.PanicIfErrorMsg(errMsg)
	udwTest.Equal(num.Get(), 2)
}

func testQueryReturnOrder() {
	db := udwSqlite3.MustNewMemoryDb()

	defer db.Close()

	errMsg := db.Query(udwSqlite3.QueryReq{
		Query: `CREATE TABLE IF NOT EXISTS abc (
k VARBINARY(255) Primary Key,
v LONGBLOB NOT NULL
) WITHOUT ROWID`,
	})
	udwErr.PanicIfErrorMsg(errMsg)
	errMsg = db.Query(udwSqlite3.QueryReq{
		Query: `REPLACE INTO abc (k, v) VALUES (?, ?),(?, ?),(?, ?)`,
		Args: [][]byte{
			[]byte("k1"),
			[]byte("v1"),
			[]byte("k2"),
			[]byte("v2"),
			[]byte("k3"),
			[]byte("v3"),
		},
	})
	udwErr.PanicIfErrorMsg(errMsg)

	keyList := []string{}
	errMsg = db.Query(udwSqlite3.QueryReq{
		Query: "SELECT k,v FROM abc ORDER BY k DESC",
		RespDataCb: func(result [][]byte) {
			keyList = append(keyList, string(result[0]))
		},
	})
	udwErr.PanicIfErrorMsg(errMsg)
	udwTest.Equal(keyList, []string{"k3", "k2", "k1"})

	keyList = []string{}
	errMsg = db.Query(udwSqlite3.QueryReq{
		Query: "SELECT k,v FROM abc WHERE k in (?,?,?)",
		Args: [][]byte{
			[]byte("k3"),
			[]byte("k1"),
			[]byte("k2"),
		},
		RespDataCb: func(result [][]byte) {
			keyList = append(keyList, string(result[0]))
		},
	})
	udwErr.PanicIfErrorMsg(errMsg)
	udwTest.Equal(keyList, []string{"k1", "k2", "k3"})
}
