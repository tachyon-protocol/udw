package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTest"
)

func TestProcessList() {
	db := udwSqlite3.MustNewMemoryDb()
	defer db.Close()
	db.MustSet("test", "test", "test")
	udwTest.Equal(db.MustGetAllK1(), []string{"test"})
	k1Status := db.MustGetAllK1Status()
	udwTest.Equal(len(k1Status), 1)
	udwTest.Equal(k1Status[0].Name, "test")
}
