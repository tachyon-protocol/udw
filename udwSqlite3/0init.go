package udwSqlite3

import (
	"github.com/tachyon-protocol/udw/udwCryptoSha3"
	"github.com/tachyon-protocol/udw/udwFile"
)

type NewDbRequest struct {
	FilePath string

	EmptyDatabaseIfDatabaseCorrupt bool

	UsingMemory bool

	JournalMode string

	Synchronous string

	DatabaseCorruptCallback func()

	EncryptPskString string

	MulitDeleteMaxSize int
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

	isFinish := false
	if db.req.FilePath != "" && db.req.UsingMemory == false {
		db.req.FilePath = udwFile.MustGetFullPath(db.req.FilePath)
		udwFile.MustMkdirForFile777(db.req.FilePath)
	}
	defer func() {
		if isFinish == false {
			db.cClose()
		}
	}()
	db.cClose()
	errMsg := db.newInnerDb(newDbReq{
		UsingMemory: db.req.UsingMemory,
		FilePath:    db.req.FilePath,
	})
	if errMsg != "" {
		panic("xsea553z38 " + errMsg)
	}
	if db.req.EncryptPskString != "" {
		keyHex := udwCryptoSha3.Sha3512ToHexStringFromString(db.req.EncryptPskString)
		db.mustExec(`PRAGMA key = "x'` + keyHex[:64] + `'"`)
		db.mustExec(`PRAGMA cipher_page_size = 4096`)
	}
	if db.req.JournalMode == "" {
		db.req.JournalMode = JournalModeWAL
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

func (db *Db) Close() {
	db.cClose()
}

func MustNewMemoryDb() *Db {
	return MustNewDb(NewDbRequest{
		UsingMemory: true,
	})
}
