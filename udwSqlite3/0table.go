package udwSqlite3

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"sort"
)

func (db *Db) MustDeleteK1(k1 string) {

	errMsg := db.querySkipResult(`DROP TABLE IF EXISTS ` + db.getTableNameFromK1(k1))
	udwErr.PanicIfErrorMsg(errMsg)
}

func (db *Db) MustEmptyDatabase() {
	k1List := db.MustGetAllK1()
	for _, k1 := range k1List {
		db.MustDeleteK1(k1)
	}
}

func (db *Db) MustGetAllK1() []string {
	tableNameList := db.GetTableNameList()
	output := []string{}
	for _, tableName := range tableNameList {

		k1, err := tableNameToK1(tableName)
		if err != nil {
			continue
		}
		output = append(output, k1)
	}
	sort.Strings(output)
	return output
}

func createTable(db *Db, k1 string) (errMsg string) {
	return db.querySkipResult(`CREATE TABLE IF NOT EXISTS ` + db.getTableNameFromK1(k1) +
		` (k BLOB Primary Key,v BLOB) WITHOUT ROWID`)
}
