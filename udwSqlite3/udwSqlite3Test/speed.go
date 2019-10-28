package udwSqlite3Test

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwMap"
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"github.com/tachyon-protocol/udw/udwTest"
	"strconv"
)

func TestSpeedSingleSetGet() {

	udwFile.MustDelete("/tmp/test_sqlite3.db")
	testDb := udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
		FilePath: "/tmp/test_sqlite3.db",
	})
	defer testDb.Close()

	intList := []string{}
	const num = 1000
	for i := 0; i < num; i++ {
		intList = append(intList, strconv.Itoa(i))
	}
	udwTest.BenchmarkSetName("singleSetGet")
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			testDb.MustSet("test", intList[i], intList[i])
		}
		for i := 0; i < num; i++ {
			ret := testDb.MustGet("test", intList[i])
			if ret != intList[i] {
				panic("TestSingleGet fail")
			}
		}
	})
	udwTest.BenchmarkSetName("MulitSetGet")
	var kvPairList []udwMap.KeyValuePair
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		keyPairList := make([]udwMap.KeyValuePair, num)
		for i := 0; i < num; i++ {
			keyPairList[i] = udwMap.KeyValuePair{
				Key:   intList[i],
				Value: intList[i],
			}
		}
		testDb.MustMulitSet("test", keyPairList)
		kvPairList = testDb.MustMulitGet("test", intList)
	})
	numSet := map[string]struct{}{}
	for i := 0; i < num; i++ {
		key := kvPairList[i].Key
		value := kvPairList[i].Value
		if key != value {
			panic("TestMulitGet fail [" + key + "] [" + value + "]")
		}
		numSet[key] = struct{}{}
	}
	udwTest.Equal(len(numSet), num)
}

func TestSpeedMemory() {
	intList := []string{}
	const num = 1000
	for i := 0; i < num; i++ {
		intList = append(intList, strconv.Itoa(i))
	}
	testDb := udwSqlite3.MustNewMemoryDb()
	defer testDb.Close()
	udwTest.BenchmarkSetName("memory singleSetGet")

	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			testDb.MustSet("test", intList[i], intList[i])
		}
		for i := 0; i < num; i++ {
			ret := testDb.MustGet("test", intList[i])
			if ret != intList[i] {
				panic("TestSingleGet fail")
			}
		}
	})
}
