package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTask"
	"github.com/tachyon-protocol/udw/udwTest"
	"strconv"
	"strings"
)

func TestThread() {
	const testDb = "/tmp/test_sqlite3.db"
	udwFile.MustDelete(testDb)
	defer udwFile.MustDelete(testDb)
	db := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath:     testDb,
		UsingWALFull: true,
	})
	const num = 1e4
	task := udwTask.New(10)
	for i := 0; i < num; i++ {
		i := i
		task.AddFunc(func() {
			iS := strconv.Itoa(i)
			msg := strings.Repeat(iS, 1000)
			db.MustSet("t1", iS, msg)
			out := db.MustGet("t1", iS)
			udwTest.Equal(out, msg)
			db.MustDelete("t1", iS)
		})
	}
	task.Close()
	db.Close()
}
