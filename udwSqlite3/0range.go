package udwSqlite3

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwMap"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwStrings"
	"strconv"
)

type GetRangeReq struct {
	K1                 string
	IsDescOrder        bool
	MinValue           string
	MaxValue           string
	MinValueNotInclude string
	MaxValueNotInclude string
	Prefix             string
	Limit              int
}

func addGetRangeSql(req GetRangeReq, sqlBuf *bytes.Buffer) (valueList [][]byte) {
	valueList = [][]byte{}
	kr := udwStrings.RangeMinMax{}
	kr.AddMinKeyInclude(req.MinValue)
	kr.AddMinKeyNotInclude(req.MinValueNotInclude)
	kr.AddMaxKeyInclude(req.MaxValue)
	kr.AddMaxKeyNotInclude(req.MaxValueNotInclude)
	kr.AddPrefix(req.Prefix)
	minInclude := kr.GetMinKeyInclude()
	maxNotInclude := kr.GetMaxKeyNotInclude()
	hasWhere := minInclude != "" || maxNotInclude != ""
	if hasWhere {
		sqlBuf.WriteString(" WHERE ")
		if minInclude != "" {
			sqlBuf.WriteString("k>=?")
			valueList = append(valueList, []byte(minInclude))
		}
		if maxNotInclude != "" {
			if minInclude != "" {
				sqlBuf.WriteString(" AND k<?")
			} else {
				sqlBuf.WriteString("k<?")
			}
			valueList = append(valueList, []byte(maxNotInclude))
		}
	}
	if req.IsDescOrder {
		sqlBuf.WriteString(" ORDER BY k DESC")
	} else {
		sqlBuf.WriteString(" ORDER BY k ASC")
	}
	if req.Limit > 0 {
		sqlBuf.WriteString(" LIMIT ")
		sqlBuf.WriteString(strconv.Itoa(req.Limit))
	}
	return valueList
}

func (db *Db) MustCountGetRange(req GetRangeReq) int {
	sqlBuf := bytes.NewBufferString(`SELECT count(1) FROM ` + db.getTableNameFromK1(req.K1))
	valueList := addGetRangeSql(req, sqlBuf)
	s, errMsg := db.queryToOneString(sqlBuf.String(), valueList...)
	if errMsg != "" {
		panic(errMsg)
	}
	return udwStrconv.MustParseInt(s)
}

func (db *Db) MustGetRange(req GetRangeReq) []udwMap.KeyValuePair {
	sqlBuf := bytes.NewBufferString(`SELECT k,v FROM ` + db.getTableNameFromK1(req.K1))
	valueList := addGetRangeSql(req, sqlBuf)
	output := []udwMap.KeyValuePair{}
	errMsg := db.Query(QueryReq{
		Query: sqlBuf.String(),
		Args:  valueList,
		RespDataCb: func(row [][]byte) {
			output = append(output, udwMap.KeyValuePair{
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

func (db *Db) MustGetRangeKeyList(req GetRangeReq) []string {
	sqlBuf := bytes.NewBufferString(`SELECT k FROM ` + db.getTableNameFromK1(req.K1))
	valueList := addGetRangeSql(req, sqlBuf)
	output := []string{}
	errMsg := db.Query(QueryReq{
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

func (db *Db) MustGetRangeValueList(req GetRangeReq) []string {
	sqlBuf := bytes.NewBufferString(`SELECT v FROM ` + db.getTableNameFromK1(req.K1))
	valueList := addGetRangeSql(req, sqlBuf)
	output := []string{}
	errMsg := db.Query(QueryReq{
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

func (db *Db) MustGetRangeKeyMap(req GetRangeReq) map[string]struct{} {
	out := map[string]struct{}{}
	db.MustGetRangeKeyListCallback(req, func(key string) {
		out[key] = struct{}{}
		return
	})
	return out
}

func (db *Db) MustGetRangeToMap(req GetRangeReq) map[string]string {
	ret := db.MustGetRange(req)
	m := make(map[string]string, len(ret))
	for _, pair := range ret {
		m[pair.Key] = pair.Value
	}
	return m
}
