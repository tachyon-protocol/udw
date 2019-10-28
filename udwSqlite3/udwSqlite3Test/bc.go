package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwMap"
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTest"
	"strconv"
)

func TestBc() {
	TestBcVersion()
	TestBcEncrypt()
	TestBcPlainDb()
}

const currentSqlite3Version = "3.28.0"

func TestBcVersion() {
	version := udwSqlite3.MustGetSqlite3Version()

	udwTest.Equal(version, currentSqlite3Version)
}

func TestBcEncrypt() {
	testPath := "/tmp/test_sqlite3.db.encrypt"
	udwSqlite3.DeleteSqliteDbFileByPath(testPath)
	udwFile.MustWriteFileWithMkdir(testPath, getBcDbEncrypt())
	defer udwFile.MustDelete(testPath)
	db1 := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath: testPath,

		EncryptPskString: "y37cc83re9sp5pqmkve39s5skgvb25c4p59dk9vx9fnuckbayz5z9bpr2r3zp3zv",
	})
	checkBcDbStatus(db1)
	db1.Close()
}

func TestBcPlainDb() {
	testPath := "/tmp/test_sqlite3.db"
	udwSqlite3.DeleteSqliteDbFileByPath(testPath)
	udwFile.MustWriteFileWithMkdir(testPath, getBcDbPlain())
	defer udwFile.MustDelete(testPath)
	db1 := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath: testPath,
	})
	checkBcDbStatus(db1)
	db1.Close()
}

func checkBcDbStatus(db1 *udwSqlite3.Db) {
	shouldPairList := getShouldBcContentPairList()
	pairList := db1.MustGetRange(udwSqlite3.GetRangeReq{
		K1: "test",
	})
	udwTest.Equal(len(pairList), len(shouldPairList))
	for i := 0; i < len(shouldPairList); i++ {
		udwTest.Equal(pairList[i].Key, shouldPairList[i].Key)
		udwTest.Equal(pairList[i].Value, shouldPairList[i].Value)
	}
	for i := 0; i < len(shouldPairList); i++ {
		udwTest.Equal(db1.MustGet("test", shouldPairList[i].Key), shouldPairList[i].Value, shouldPairList[i].Value)
	}
	udwTest.Equal(db1.MustCountGetRange(udwSqlite3.GetRangeReq{
		K1: "test",
	}), len(shouldPairList))
}

func getShouldBcContentPairList() (pairList []udwMap.KeyValuePair) {
	pairList = append(pairList, udwMap.KeyValuePair{
		Key: "0", Value: "0",
	}, udwMap.KeyValuePair{
		Key: "000", Value: "0010",
	}, udwMap.KeyValuePair{
		Key: "1", Value: "1",
	}, udwMap.KeyValuePair{
		Key: "10", Value: "10",
	})
	for i := 2; i < 10; i++ {
		s := strconv.Itoa(i)
		pairList = append(pairList, udwMap.KeyValuePair{
			Key: s, Value: s,
		})
	}
	return pairList
}
