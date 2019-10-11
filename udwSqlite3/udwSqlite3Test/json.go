package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTest"
)

func TestJson() {
	MustRunTestDb(func(db *udwSqlite3.Db) {
		db.MustEmptyDatabase()
		a := map[string]string{}
		a["1"] = "1"
		db.MustSetJson("test", "1", a)
		var b map[string]string
		db.MustGetJson("test", "1", &b)
		udwTest.Equal(len(b), 1)
		udwTest.Equal(b["1"], "1")
	})
}
