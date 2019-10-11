package udwSqlite3

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwLog"
	"github.com/tachyon-protocol/udw/udwSqlite3/sqlite3"
	"github.com/tachyon-protocol/udw/udwStrings"
)

type setExecReq struct {
	k1           string
	sql          string
	valueBuf     [][]byte
	respStatusCb func(status sqlite3.QueryRespStatus)
	UseStmtCache bool
}

func (db *Db) setExec(req setExecReq) (errMsg string) {
	for i := 0; i < 3; i++ {
		errMsg = db.db.Query(sqlite3.QueryReq{
			Query:        req.sql,
			Args:         req.valueBuf,
			RespStatusCb: req.respStatusCb,
			UseStmtCache: req.UseStmtCache,
		})
		if errMsg != "" {

			if errorIsTableNotExist(errMsg) {
				errMsg = createTable(db, req.k1)
				if errMsg != "" {

					if db.handleEmptyDatabaseWhenCorrupt(errMsg) {
						continue
					}
					return "err: " + errMsg + " sql: " + req.sql
				}
				continue
			}

			if db.handleEmptyDatabaseWhenCorrupt(errMsg) {
				continue
			}
			return "err: " + errMsg + " sql: " + req.sql
		}
		return ""
	}
	return "[mustSetExec]try too many times sql: " + req.sql
}

func (db *Db) mustSetExec(k1 string, sql string, valueBuf [][]byte) {
	errMsg := db.setExec(setExecReq{
		k1:       k1,
		sql:      sql,
		valueBuf: valueBuf,
	})
	if errMsg != "" {
		panic(errMsg)
	}
	return
}

func (db *Db) exec(sql string) (errMsg string) {
	for i := 0; i < 2; i++ {
		errMsg := db.querySkipResult(sql)
		if errMsg != "" {
			if db.handleEmptyDatabaseWhenCorrupt(errMsg) {
				continue
			}
			return errMsg
		}
		return ""
	}
	return "[exec]try too many times sql:[" + sql + "] err:[" + errMsg + "]"
}

func (db *Db) mustExec(sql string) {
	errMsg := db.exec(sql)
	if errMsg != "" {
		panic(errMsg)
	}
	return
}

func (db *Db) handleEmptyDatabaseWhenCorrupt(errMsg string) (ok bool) {
	if IsErrorDatabaseCorrupt(errMsg) {

		if db.req.DatabaseCorruptCallback != nil {
			db.req.DatabaseCorruptCallback()
		}
		if db.req.EmptyDatabaseIfDatabaseCorrupt {
			udwLog.Log("erorr", "[udwSqlite3.isEmptyDatabaseWhenCorrupt]", "DatabaseCorrupt and emtry database now.", db.req.FilePath, errMsg)
			db.db.Close()
			db.db = nil
			DeleteSqliteDbFileByPath(db.req.FilePath)
			db.mustInitDbL1()
			return true
		}
	}
	return false
}

func DeleteSqliteDbFileByPath(path string) {
	if udwFile.MustIsFile(path) {

		udwFile.MustDelete(path)
	}
	if udwFile.MustIsFile(path + "-shm") {

		udwFile.MustDelete(path + "-shm")
	}
	if udwFile.MustIsFile(path + "-wal") {

		udwFile.MustDelete(path + "-wal")
	}
}

func stoB(s string) (b []byte) {
	return udwStrings.GetByteArrayFromStringNoAlloc(s)
}

func btoS(b []byte) (s string) {
	return udwStrings.GetStringFromByteArrayNoAlloc(b)
}

func (c *Db) getArgumentListCache(needSize int) [][]byte {
	if cap(c.argumentListCache) < needSize {
		c.argumentListCache = make([][]byte, needSize)
	}
	return c.argumentListCache[:needSize]
}
