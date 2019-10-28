package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTest"
)

func testMulitGet() {
	db := udwSqlite3.MustNewMemoryDb()
	defer db.Close()
	output := db.MustMulitGet("test", []string{"1"})
	udwTest.Equal(len(output), 0)

	m := db.MustMulitExist("test", []string{"1"})
	udwTest.Equal(len(m), 0)

	db.MustMulitSetByKeySet("test", map[string]struct{}{
		"1": struct{}{},
		"2": struct{}{},
	})
	output = db.MustMulitGet("test", []string{"1"})
	udwTest.Equal(len(output), 1)
	udwTest.Equal(output[0].Key, "1")

	m = db.MustMulitExist("test", []string{"1"})
	udwTest.Equal(len(m), 1)
	_, ok := m["1"]
	udwTest.Equal(ok, true)

	db.MustMulitSetByKeySet("test", nil)
}
