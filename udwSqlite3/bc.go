package udwSqlite3

import (
	"github.com/tachyon-protocol/udw/udwMap"
	"github.com/tachyon-protocol/udw/udwSqlite3/sqlite3"
	"strings"
)

const k1BcTmpTable = "k1BcTmpTable"

func (db *Db) MustHandleTableBc() {
	tablePairList := []udwMap.KeyValuePair{}
	errMsg := db.db.Query(sqlite3.QueryReq{
		Query: "SELECT name,sql FROM SQLITE_MASTER where type='table'",
		RespDataCb: func(row [][]byte) {
			tablePairList = append(tablePairList, udwMap.KeyValuePair{
				Key:   string(row[0]),
				Value: string(row[1]),
			})
		},
		ColumnCount: 2,
	})
	if errMsg != "" {
		panic(errMsg)
	}
	for _, pair := range tablePairList {
		if strings.HasPrefix(pair.Key, "udw_") == false {
			continue
		}
		if strings.Contains(pair.Value, "VARBINARY(255)") == false {
			continue
		}
		k1, err := tableNameToK1(pair.Key)
		if err != nil {
			panic(err)
		}
		dataPairList := db.MustGetRange(MustGetRangeRequest{
			K1: k1,
		})
		db.MustDeleteK1(k1BcTmpTable)
		msb := db.NewMsb(1)
		for _, pair := range dataPairList {
			msb.Set(k1BcTmpTable, pair.Key, pair.Value)
		}
		msb.MustClose()
		newTableName := k1ToTableName(k1BcTmpTable)
		db.mustExec(`DROP TABLE ` + pair.Key)
		db.mustExec(`ALTER TABLE ` + newTableName + ` RENAME TO ` + pair.Key)
	}
}
