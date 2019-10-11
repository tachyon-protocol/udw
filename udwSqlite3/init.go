package udwSqlite3

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwCryptoSha3"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwSqlite3/sqlite3"
	"sync"
)

type NewDbRequest struct {
	FilePath string

	EmptyDatabaseIfDatabaseCorrupt bool

	UsingMemory bool

	UsingWALFull bool

	HighPerformanceModeWithDataLoss bool

	JournalMode string

	Synchronous string

	DatabaseCorruptCallback func()

	EncryptPskString string

	MulitDeleteMaxSize int
	SkipBcCheck        bool
}

const JournalModeDELETE = "DELETE"
const JournalModeWAL = "WAL"
const JournalModeOFF = "OFF"

const SynchronousOFF = "OFF"
const SynchronousNORMAL = "NORMAL"
const SynchronousFULL = "FULL"

func MustNewDb(req NewDbRequest) *Db {
	if req.MulitDeleteMaxSize == 0 {
		req.MulitDeleteMaxSize = defaultMulitDeleteMaxSize
	}
	db := &Db{
		req:               req,
		tableNameCacheMap: map[string]string{},
	}
	db.initLevelLocker.Lock()
	db.initLevel = 3
	db.initLevelLocker.Unlock()
	db.mustInitDbL1()
	if req.SkipBcCheck == false {
		db.MustHandleTableBc()
	}
	return db
}

func (db *Db) mustInitDbL1() {
	db.initLevelLocker.Lock()
	db.initLevel--
	if db.initLevel <= 0 {
		db.initLevelLocker.Unlock()
		panic("[mustInitDbL1] too many init db")
	}
	db.initLevelLocker.Unlock()
	var mysqlDb *sqlite3.Db
	isFinish := false
	if db.req.FilePath != "" && db.req.UsingMemory == false {
		db.req.FilePath = udwFile.MustGetFullPath(db.req.FilePath)
		udwFile.MustMkdirForFile777(db.req.FilePath)
	}
	mysqlDb, errMsg := sqlite3.NewDb(sqlite3.NewDbReq{
		UsingMemory: db.req.UsingMemory,
		FilePath:    db.req.FilePath,
	})
	if errMsg != "" {
		panic("xsea553z38 " + errMsg)
	}

	defer func() {
		if isFinish == false && mysqlDb != nil {
			mysqlDb.Close()
		}
	}()
	if db.db != nil {
		db.db.Close()
	}
	db.db = mysqlDb
	if db.req.EncryptPskString != "" {
		keyHex := udwCryptoSha3.Sha3512ToHexStringFromString(db.req.EncryptPskString)
		db.mustExec(`PRAGMA key = "x'` + keyHex[:64] + `'"`)
		db.mustExec(`PRAGMA cipher_page_size = 4096`)
	}
	if db.req.HighPerformanceModeWithDataLoss == true {

		db.req.JournalMode = JournalModeOFF
		db.req.Synchronous = SynchronousOFF
	} else if db.req.UsingWALFull {
		db.req.JournalMode = JournalModeWAL
		db.req.Synchronous = SynchronousFULL
	}
	if db.req.JournalMode == "" {
		db.req.JournalMode = JournalModeDELETE
	}
	switch db.req.JournalMode {
	case JournalModeWAL, JournalModeOFF, JournalModeDELETE:
	default:
		panic("[udwSqlite3] unknow JournalMode " + db.req.JournalMode)
	}
	if db.req.JournalMode != JournalModeDELETE {
		db.mustExec("PRAGMA journal_mode=" + db.req.JournalMode)
	}
	if db.req.Synchronous == "" {
		db.req.Synchronous = SynchronousFULL
	}
	switch db.req.Synchronous {
	case SynchronousOFF, SynchronousFULL, SynchronousNORMAL:
	default:
		panic("[udwSqlite3] unknow Synchronous " + db.req.JournalMode)
	}
	if db.req.JournalMode != SynchronousFULL {
		db.mustExec("PRAGMA synchronous=" + db.req.JournalMode)
	}

	db.initLevelLocker.Lock()
	db.initLevel = 3
	db.initLevelLocker.Unlock()
	isFinish = true
}

type Db struct {
	db                      *sqlite3.Db
	req                     NewDbRequest
	initLevel               int
	initLevelLocker         sync.Mutex
	locker                  sync.Mutex
	queryBuf                udwBytes.BufWriter
	argumentListCache       [][]byte
	tableNameCacheMapLocker sync.Mutex
	tableNameCacheMap       map[string]string
	cacheOneResult          string
	cacheOneResultDataCb    func(valueList [][]byte)
}

func (db *Db) Close() {
	db.db.Close()
}

func MustNewMemoryDb() *Db {
	return MustNewDb(NewDbRequest{
		UsingMemory: true,
	})
}
