package udwSqlite3

import (
	"github.com/tachyon-protocol/udw/udwSqlite3/sqlite3"
	"sort"
)

func (db *Db) MustSet(k1 string, k2 string, value string) {
	errMsg := db.Set(k1, k2, value)
	if errMsg != "" {
		panic(errMsg)
	}
}

func (db *Db) Set(k1 string, k2 string, value string) (errMsg string) {
	db.locker.Lock()
	db.queryBuf.Reset()
	db.queryBuf.WriteString(`REPLACE INTO `)
	db.queryBuf.WriteString(db.getEscapedTableName(k1))
	db.queryBuf.WriteString(` (k, v) VALUES (?, ?)`)
	db.queryBuf.WriteByte(0)
	argumentList := db.getArgumentListCache(2)
	argumentList[0] = stoB(k2)
	argumentList[1] = stoB(value)
	errMsg = db.setExec(setExecReq{
		k1:  k1,
		sql: btoS(db.queryBuf.GetBytes()),

		valueBuf:     argumentList,
		UseStmtCache: true,
	})
	db.locker.Unlock()
	return errMsg

}

func (db *Db) MustInsertAndReturnHasSucc(k1 string, k2 string, value string) bool {
	hasSucc := false
	errMsg := db.setExec(setExecReq{
		k1:       k1,
		sql:      `INSERT OR IGNORE INTO ` + db.getEscapedTableName(k1) + ` (k, v) VALUES (?, ?)`,
		valueBuf: [][]byte{[]byte(k2), []byte(value)},
		respStatusCb: func(resp sqlite3.QueryRespStatus) {
			hasSucc = resp.AffectedRows > 0
		},
	})
	if errMsg != "" {
		panic(errMsg)
	}
	return hasSucc

}

func (db *Db) MustInsert(k1 string, k2 string, value string) {
	db.mustSetExec(k1, `INSERT INTO `+db.getEscapedTableName(k1)+` (k, v) VALUES (?, ?)`, [][]byte{[]byte(k2), []byte(value)})
	return
}

func (db *Db) MustUpdate(k1 string, k2 string, value string) {
	db.mustSetExec(k1, `UPDATE `+db.getEscapedTableName(k1)+` SET v=? WHERE k=?`, [][]byte{[]byte(value), []byte(k2)})
	return
}

type KeyValuePair struct {
	Key   string
	Value string
}

func sortKeyValuePairList(pairList []KeyValuePair) {
	sort.Slice(pairList, func(i int, j int) bool {
		return pairList[i].Key < pairList[j].Key
	})
}

func (db *Db) MustMulitSetByKeySet(k1 string, keySet map[string]struct{}) {
	if len(keySet) == 0 {
		return
	}
	keyValuePairList := make([]KeyValuePair, 0, len(keySet))
	for k := range keySet {
		keyValuePairList = append(keyValuePairList, KeyValuePair{
			Key:   k,
			Value: "1",
		})
	}
	db.MustMulitSet(k1, keyValuePairList)
}

func (db *Db) MustMulitSet(k1 string, pairList []KeyValuePair) {
	sortKeyValuePairList(pairList)

	for {
		if len(pairList) == 0 {
			return
		}
		toIndex := mustMulitSetMaxSize
		if toIndex > len(pairList) {
			toIndex = len(pairList)
		}
		db.mustMulitSetL1(k1, pairList[:toIndex])
		pairList = pairList[toIndex:]
	}
}

func (db *Db) mustMulitSetL1(k1 string, pairList []KeyValuePair) {
	db.locker.Lock()
	db.queryBuf.Reset()
	db.queryBuf.WriteString(`REPLACE INTO `)
	db.queryBuf.WriteString(db.getEscapedTableName(k1))
	db.queryBuf.WriteString(` (k, v) VALUES `)

	valueBuf := db.getArgumentListCache(len(pairList) * 2)
	l := len(pairList)
	for i, pair := range pairList {
		if i < l-1 {
			db.queryBuf.WriteString("(?,?),")
		} else {
			db.queryBuf.WriteString("(?,?)")
		}
		valueBuf[2*i] = stoB(pair.Key)
		valueBuf[2*i+1] = stoB(pair.Value)
	}
	db.queryBuf.WriteByte(0)
	errMsg := db.setExec(setExecReq{
		k1:       k1,
		sql:      btoS(db.queryBuf.GetBytes()),
		valueBuf: valueBuf,
	})
	db.locker.Unlock()
	if errMsg != "" {
		panic(errMsg)
	}
	return
}
