package udwSqlite3

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwHex"
	"strings"
)

func mustEscapeTableOrDatabaseName(s string) string {
	if strings.Contains(s, "\x00") {
		panic("d2daft7kd5")
	}
	return "'" + strings.Replace(s, "'", "''", -1) + "'"
}

func (db *Db) getTableNameFromK1(k1 string) string {
	if k1 == "" {
		panic(errors.New(`[getEscapedTableName] k1==""`))
	}
	db.tableNameCacheMapLocker.Lock()
	tableName := db.tableNameCacheMap[k1]
	if tableName != "" {
		db.tableNameCacheMapLocker.Unlock()
		return tableName
	}
	tableName = k1ToTableName(k1)
	if len(db.tableNameCacheMap) > 10000 {
		for k := range db.tableNameCacheMap {
			delete(db.tableNameCacheMap, k)
		}
	}
	db.tableNameCacheMap[k1] = tableName
	db.tableNameCacheMapLocker.Unlock()
	return tableName
}

func k1ToTableName(k1 string) string {

	return "'" + udwHex.EncodeStringToString(k1) + "'"
}

func tableNameToK1(tableName string) (string, error) {

	return udwHex.DecodeStringToString(tableName)
}
