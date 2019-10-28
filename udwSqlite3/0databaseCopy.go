package udwSqlite3

import (
	"github.com/tachyon-protocol/udw/udwTask"
	"sync"
)

type dbCopyContext struct {
	dbFrom       *Db
	dbTo         *Db
	copyedRowNum int64
	locker       sync.Mutex
}

func mustDatabaseCopyL1OneTable(ctx *dbCopyContext, Fromk1Name string, Tok1Name string) {
	tasker := udwTask.New(10)
	lastMaxKey := ""
	for {
		keyPairList := ctx.dbFrom.MustGetRange(GetRangeReq{
			K1:       Fromk1Name,
			MinValue: lastMaxKey,
			Limit:    databaseCopyPerQuery,
		})
		tasker.AddFunc(func() {
			ctx.dbTo.MustMulitSet(Tok1Name, keyPairList)
			ctx.locker.Lock()
			ctx.copyedRowNum += int64(len(keyPairList))
			ctx.locker.Unlock()
		})

		if len(keyPairList) < databaseCopyPerQuery {
			break
		}
		lastMaxKey = keyPairList[databaseCopyPerQuery-1].Key
	}
	tasker.Close()
}

func (db *Db) MustTableCopy(fromK1Name string, toK1Name string) {
	ctx := &dbCopyContext{
		dbFrom: db,
		dbTo:   db,
	}
	mustDatabaseCopyL1OneTable(ctx, fromK1Name, toK1Name)
}
