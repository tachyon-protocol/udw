package udwSqlite3

import (
	"sort"
)

func (db *Db) MustQuery(req QueryReq) {
	errMsg := db.Query(req)
	if errMsg != "" {
		panic(errMsg)
	}
}

func (db *Db) GetAllDataInTableToRowMap(tableName string) []map[string]string {
	return db.mustQueryToMapRowList("SELECT * FROM " + mustEscapeTableOrDatabaseName(tableName) + " ")
}

func (db *Db) GetTableNameList() []string {
	output := []string{}
	errMsg := db.Query(QueryReq{
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
	qReq := QueryReq{
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
	errMsg := db.Query(qReq)
	if errMsg != "" {
		panic(errMsg)
	}
	return mapRowList
}

func (db *Db) querySkipResult(query string, args ...[]byte) (errMsg string) {
	return db.Query(QueryReq{
		Query: query,
		Args:  args,
	})
}

func (db *Db) queryToOneString(query string, args ...[]byte) (s string, errMsg string) {
	qReq := QueryReq{
		Query: query,
		Args:  args,
		RespDataCb: func(valueList [][]byte) {
			if s == "" {
				s = string(valueList[0])
			}
		},
		ColumnCount: 1,
	}
	errMsg = db.Query(qReq)
	return s, errMsg
}
