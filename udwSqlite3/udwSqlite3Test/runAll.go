package udwSqlite3Test

func RunAllTest() {
	TestBc()

	TestMustGetAllK1()
	testZeroBug()
	testInt()
	testMustUpdate()

	testSimpleQuery()
	testSimpleQueryAffectedRows()
	testQueryReturnOrder()
	TestMustTableCopy()
	TestEncrypt()
	TestThread()
	TestInitDatabaseCorrupt()
	TestInitDatabaseCorrupt2()
	TestInitDatabaseCorrupt3()
	TestMustGetRange()
	TestThreadSafe()

	TestRangeCallback2()
	TestNoSql()
	TestInsert()
	TestSet()
	TestMustInsertAndReturnExist()
	TestMustMulitDelete()

	testMulitGet()
	testGetAllDataInTableToRowMap()
	testThreadSafe2()
	testMemoryDb()

	TestSpeedMemory()
	TestSpeedSingleSetGet()
}
