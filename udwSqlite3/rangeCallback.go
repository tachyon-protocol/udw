package udwSqlite3

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwDbI"
	"github.com/tachyon-protocol/udw/udwSqlite3/sqlite3"
)

func (db *Db) MustGetRangeValueListCallback(req MustGetRangeRequest, visitor func(value string)) {
	sqlBuf := bytes.NewBufferString(`SELECT v FROM ` + db.getEscapedTableName(req.K1))
	valueList := addGetRangeSql(req, sqlBuf)
	errMsg := db.db.Query(sqlite3.QueryReq{
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

func (db *Db) MustGetRangeCallback(req MustGetRangeRequest, visitor func(key string, value string) bool) {
	sqlBuf := bytes.NewBufferString(`SELECT k,v FROM ` + db.getEscapedTableName(req.K1))
	valueList := addGetRangeSql(req, sqlBuf)
	errMsg := db.db.Query(sqlite3.QueryReq{
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

func (db *Db) IMustGetRangeCallback(req udwDbI.GetRangeReq, visitor func(key string, value string)) {
	db.MustGetRangeCallback(MustGetRangeRequest{
		K1:                 req.K1,
		IsDescOrder:        req.IsDescOrder,
		MinValue:           req.MinValue,
		MaxValue:           req.MaxValue,
		MinValueNotInclude: req.MinValueNotInclude,
		MaxValueNotInclude: req.MaxValueNotInclude,
		Prefix:             req.Prefix,
		Limit:              req.Limit,
	}, func(key string, value string) bool {
		visitor(key, value)
		return true
	})
}

func (db *Db) MustGetRangeKeyListCallback(req MustGetRangeRequest, visitor func(key string) bool) {
	sqlBuf := bytes.NewBufferString(`SELECT k FROM ` + db.getEscapedTableName(req.K1))
	valueList := addGetRangeSql(req, sqlBuf)
	errMsg := db.db.Query(sqlite3.QueryReq{
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
