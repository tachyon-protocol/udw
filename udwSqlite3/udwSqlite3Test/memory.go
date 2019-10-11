package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTask"
	"github.com/tachyon-protocol/udw/udwTest"
)

func testMemoryDb() {
	tasker := udwTask.New(10)
	for i := 0; i < 100; i++ {
		tasker.AddFunc(func() {
			db := udwSqlite3.MustNewMemoryDb()
			db.MustSet("test", "1", "2")
			udwTest.Equal(db.MustGet("test", "1"), "2")
			db.Close()
		})
	}
	tasker.Close()
}
