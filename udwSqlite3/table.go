package udwSqlite3

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"sort"
	"strings"
	"time"
)

func (db *Db) MustDeleteK1(k1 string) {

	errMsg := db.querySkipResult(`DROP TABLE IF EXISTS ` + db.getEscapedTableName(k1))
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
		if !strings.HasPrefix(tableName, "udw_") {
			continue
		}
		k1, err := tableNameToK1(tableName)
		if err != nil {
			continue
		}
		output = append(output, k1)
	}
	sort.Strings(output)
	return output
}

type K1Status struct {
	Name         string
	RowNum       int64
	AvgRowLength int64
	DataLength   int64
	CreateTime   time.Time
}

func (db *Db) MustGetAllK1Status() []K1Status {
	k1List := db.MustGetAllK1()
	out := []K1Status{}
	for _, k1 := range k1List {
		out = append(out, K1Status{
			Name: k1,
		})
	}
	return out

}

func (db *Db) MustGetVersion() string {
	s, errMsg := db.queryToOneString("SELECT sqlite_version()")
	if errMsg != "" {
		panic(errMsg)
	}
	return s
}

func createTable(db *Db, k1 string) (errMsg string) {
	return db.querySkipResult(`CREATE TABLE IF NOT EXISTS ` + db.getEscapedTableName(k1) +
		` (k BLOB Primary Key,v BLOB) WITHOUT ROWID`)

}
