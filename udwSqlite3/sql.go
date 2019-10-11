package udwSqlite3

import (
	"github.com/tachyon-protocol/udw/udwSqlite3/sqlite3"
	"sort"
	"strings"
)

func (db *Db) MustQuery(req sqlite3.QueryReq) {
	errMsg := db.db.Query(req)
	if errMsg != "" {
		panic(errMsg)
	}
}

func (db *Db) GetAllDataInTableToRowMap(tableName string) []map[string]string {
	if strings.Contains(tableName, "`") {
		panic("awhufezbgx")
	}
	return db.mustQueryToMapRowList("SELECT * FROM `" + tableName + "` ")
}

func (db *Db) GetTableNameList() []string {
	output := []string{}
	errMsg := db.db.Query(sqlite3.QueryReq{
		Query: "SELECT name FROM SQLITE_MASTER where type='table'",
		RespDataCb: func(row [][]byte) {
			output = append(output, string(row[0]))
		},
		ColumnCount: 1,
	})
	if errMsg != "" {
		panic(errMsg)
	}
	sort.Strings(output)
	return output

}

func (db *Db) mustQueryToMapRowList(query string, args ...[]byte) (mapRowList []map[string]string) {
	mapRowList = []map[string]string{}
	columnNameList := []string{}
	qReq := sqlite3.QueryReq{
		Query: query,
		Args:  args,
		ColumnsCb: func(cl [][]byte) {
			for _, b := range cl {
				columnNameList = append(columnNameList, string(b))
			}
		},
		RespDataCb: func(valueList [][]byte) {
			thisMap := map[string]string{}
			for i, c := range columnNameList {
				thisMap[c] = string(valueList[i])
			}
			mapRowList = append(mapRowList, thisMap)
		},
	}
	errMsg := db.db.Query(qReq)
	if errMsg != "" {
		panic(errMsg)
	}
	return mapRowList
}

func (db *Db) querySkipResult(query string, args ...[]byte) (errMsg string) {
	return db.db.Query(sqlite3.QueryReq{
		Query: query,
		Args:  args,
	})
}

func (db *Db) queryToOneString(query string, args ...[]byte) (s string, errMsg string) {
	qReq := sqlite3.QueryReq{
		Query: query,
		Args:  args,
		RespDataCb: func(valueList [][]byte) {
			if s == "" {
				s = string(valueList[0])
			}
		},
		ColumnCount: 1,
	}
	errMsg = db.db.Query(qReq)
	return s, errMsg
}

func (db *Db) queryToOneString2__noLock(query string, args ...[]byte) (s string, errMsg string) {
	db.cacheOneResult = ""
	if db.cacheOneResultDataCb == nil {
		db.cacheOneResultDataCb = func(valueList [][]byte) {

			if db.cacheOneResult == "" {
				db.cacheOneResult = string(valueList[0])
			}
		}
	}
	qReq := sqlite3.QueryReq{
		Query:        query,
		Args:         args,
		RespDataCb:   db.cacheOneResultDataCb,
		ColumnCount:  1,
		UseStmtCache: true,
	}

	errMsg = db.db.Query(qReq)
	return db.cacheOneResult, errMsg
}
