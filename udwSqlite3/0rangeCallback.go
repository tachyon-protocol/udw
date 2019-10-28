package udwSqlite3

import (
	"bytes"
)

func (db *Db) MustGetRangeCallback(req GetRangeReq, visitor func(key string, value string)) {
	sqlBuf := bytes.NewBufferString(`SELECT k,v FROM ` + db.getTableNameFromK1(req.K1))
	valueList := addGetRangeSql(req, sqlBuf)
	errMsg := db.Query(QueryReq{
		Query: sqlBuf.String(),
		Args:  valueList,
		RespDataCb: func(row [][]byte) {
			visitor(string(row[0]), string(row[1]))
		},
		ColumnCount: 2,
	})
	if errMsg != "" {
		if errorIsTableNotExist(errMsg) {
			return
		}
		panic(errMsg)
	}
	return
}

func (db *Db) MustGetRangeKeyListCallback(req GetRangeReq, visitor func(key string)) {
	sqlBuf := bytes.NewBufferString(`SELECT k FROM ` + db.getTableNameFromK1(req.K1))
	valueList := addGetRangeSql(req, sqlBuf)
	errMsg := db.Query(QueryReq{
		Query: sqlBuf.String(),
		Args:  valueList,
		RespDataCb: func(row [][]byte) {
			visitor(string(row[0]))
		},
		ColumnCount: 1,
	})
	if errMsg != "" {
		if errorIsTableNotExist(errMsg) {
			return
		}
		panic(errMsg)
	}
	return
}
