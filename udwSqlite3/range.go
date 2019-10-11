package udwSqlite3

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwSqlite3/sqlite3"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwStrings"
	"strconv"
	"strings"
)

type MustGetRangeRequest struct {
	K1                 string
	IsDescOrder        bool
	MinValue           string
	MaxValue           string
	MinValueNotInclude string
	MaxValueNotInclude string
	Prefix             string
	Limit              int
}

func addGetRangeSql(req MustGetRangeRequest, sqlBuf *bytes.Buffer) (valueList [][]byte) {
	whereSqlList := []string{}
	if req.MinValue != "" {
		whereSqlList = append(whereSqlList, "k>=?")
		valueList = append(valueList, []byte(req.MinValue))
	}
	if req.MinValueNotInclude != "" {
		whereSqlList = append(whereSqlList, "k>?")
		valueList = append(valueList, []byte(req.MinValueNotInclude))
	}
	if req.MaxValue != "" {
		whereSqlList = append(whereSqlList, "k<=?")
		valueList = append(valueList, []byte(req.MaxValue))
	}
	if req.MaxValueNotInclude != "" {
		whereSqlList = append(whereSqlList, "k<?")
		valueList = append(valueList, []byte(req.MaxValueNotInclude))
	}
	if req.Prefix != "" {
		whereSqlList = append(whereSqlList, "k>=?", "k<?")
		valueList = append(valueList, []byte(req.Prefix), []byte(udwStrings.GetSmallestBiggerStringPrefix(req.Prefix)))
	}
	if len(whereSqlList) > 0 {
		sqlBuf.WriteString(" WHERE ")
		sqlBuf.WriteString(strings.Join(whereSqlList, " AND "))
	}
	if req.IsDescOrder {
		sqlBuf.WriteString(" ORDER BY k DESC")
	} else {
		sqlBuf.WriteString(" ORDER BY k ASC")
	}
	if req.Limit > 0 {
		sqlBuf.WriteString(" LIMIT " + strconv.Itoa(req.Limit))
	}
	return valueList
}

func (db *Db) MustCountGetRange(req MustGetRangeRequest) int {
	sqlBuf := bytes.NewBufferString(`SELECT count(1) FROM ` + db.getEscapedTableName(req.K1))
	valueList := addGetRangeSql(req, sqlBuf)
	s, errMsg := db.queryToOneString(sqlBuf.String(), valueList...)
	if errMsg != "" {
		panic(errMsg)
	}
	return udwStrconv.MustParseInt(s)

}

func (db *Db) MustGetRange(req MustGetRangeRequest) []KeyValuePair {
	sqlBuf := bytes.NewBufferString(`SELECT k,v FROM ` + db.getEscapedTableName(req.K1))
	valueList := addGetRangeSql(req, sqlBuf)
	output := []KeyValuePair{}
	errMsg := db.db.Query(sqlite3.QueryReq{
		Query: sqlBuf.String(),
		Args:  valueList,
		RespDataCb: func(row [][]byte) {
			output = append(output, KeyValuePair{
				Key:   string(row[0]),
				Value: string(row[1]),
			})
		},
		ColumnCount: 2,
	})
	if errMsg != "" {
		if errorIsTableNotExist(errMsg) {
			return nil
		}
		panic(errMsg)
	}
	return output

}

func (db *Db) MustGetRangeKeyList(req MustGetRangeRequest) []string {
	sqlBuf := bytes.NewBufferString(`SELECT k FROM ` + db.getEscapedTableName(req.K1))
	valueList := addGetRangeSql(req, sqlBuf)
	output := []string{}
	errMsg := db.db.Query(sqlite3.QueryReq{
		Query: sqlBuf.String(),
		Args:  valueList,
		RespDataCb: func(row [][]byte) {
			output = append(output, string(row[0]))
		},
		ColumnCount: 1,
	})
	if errMsg != "" {
		if errorIsTableNotExist(errMsg) {
			return nil
		}
		panic(errMsg)
	}
	return output

}

func (db *Db) MustGetRangeValueList(req MustGetRangeRequest) []string {
	sqlBuf := bytes.NewBufferString(`SELECT v FROM ` + db.getEscapedTableName(req.K1))
	valueList := addGetRangeSql(req, sqlBuf)
	output := []string{}
	errMsg := db.db.Query(sqlite3.QueryReq{
		Query: sqlBuf.String(),
		Args:  valueList,
		RespDataCb: func(row [][]byte) {
			output = append(output, string(row[0]))
		},
		ColumnCount: 1,
	})
	if errMsg != "" {
		if errorIsTableNotExist(errMsg) {
			return nil
		}
		panic(errMsg)
	}
	return output

}

func (db *Db) MustGetRangeKeyMap(req MustGetRangeRequest) map[string]struct{} {
	out := map[string]struct{}{}
	db.MustGetRangeKeyListCallback(req, func(key string) bool {
		out[key] = struct{}{}
		return true
	})
	return out
}

func (db *Db) MustGetRangeToMap(req MustGetRangeRequest) map[string]string {
	ret := db.MustGetRange(req)
	m := make(map[string]string, len(ret))
	for _, pair := range ret {
		m[pair.Key] = pair.Value
	}
	return m
}
