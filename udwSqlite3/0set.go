package udwSqlite3

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwMap"
)

func (db *Db) MustSet(k1 string, k2 string, value string) {
	errMsg := db.Set(k1, k2, value)
	if errMsg != "" {
		panic(errMsg)
	}
}

func (db *Db) Set(k1 string, k2 string, value string) (errMsg string) {
	queryBuf := &udwBytes.BufWriter{}
	queryBuf.WriteString(`REPLACE INTO `)
	queryBuf.WriteString(db.getTableNameFromK1(k1))
	queryBuf.WriteString(` (k, v) VALUES (?, ?)`)
	queryBuf.WriteByte(0)
	argumentList := db.getArgumentList(2)
	argumentList[0] = stoB(k2)
	argumentList[1] = stoB(value)
	errMsg = db.setExec(setExecReq{
		k1:  k1,
		sql: btoS(queryBuf.GetBytes()),

		valueBuf:     argumentList,
		UseStmtCache: true,
	})
	return errMsg
}

func (db *Db) MustInsertAndReturnHasSucc(k1 string, k2 string, value string) bool {
	hasSucc := false
	errMsg := db.setExec(setExecReq{
		k1:       k1,
		sql:      `INSERT OR IGNORE INTO ` + db.getTableNameFromK1(k1) + ` (k, v) VALUES (?, ?)`,
		valueBuf: [][]byte{[]byte(k2), []byte(value)},
		respStatusCb: func(resp QueryRespStatus) {
			hasSucc = resp.AffectedRows > 0
		},
	})
	if errMsg != "" {
		panic(errMsg)
	}
	return hasSucc
}

func (db *Db) MustInsert(k1 string, k2 string, value string) {
	db.mustSetExec(k1, `INSERT INTO `+db.getTableNameFromK1(k1)+` (k, v) VALUES (?, ?)`, [][]byte{[]byte(k2), []byte(value)})
	return
}

func (db *Db) MustUpdate(k1 string, k2 string, value string) {
	db.mustSetExec(k1, `UPDATE `+db.getTableNameFromK1(k1)+` SET v=? WHERE k=?`, [][]byte{[]byte(value), []byte(k2)})
	return
}

func (db *Db) MustMulitSetByKeySet(k1 string, keySet map[string]struct{}) {
	if len(keySet) == 0 {
		return
	}
	keyValuePairList := make([]udwMap.KeyValuePair, 0, len(keySet))
	for k := range keySet {
		keyValuePairList = append(keyValuePairList, udwMap.KeyValuePair{
			Key:   k,
			Value: "1",
		})
	}
	db.MustMulitSet(k1, keyValuePairList)
}

func (db *Db) MustMulitSet(k1 string, pairList []udwMap.KeyValuePair) {
	udwMap.SortKeyValuePairList(pairList)

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

func (db *Db) mustMulitSetL1(k1 string, pairList []udwMap.KeyValuePair) {
	queryBuf := &udwBytes.BufWriter{}
	queryBuf.WriteString(`REPLACE INTO `)
	queryBuf.WriteString(db.getTableNameFromK1(k1))
	queryBuf.WriteString(` (k, v) VALUES `)
	valueBuf := db.getArgumentList(len(pairList) * 2)
	l := len(pairList)
	for i, pair := range pairList {
		if i < l-1 {
			queryBuf.WriteString("(?,?),")
		} else {
			queryBuf.WriteString("(?,?)")
		}
		valueBuf[2*i] = stoB(pair.Key)
		valueBuf[2*i+1] = stoB(pair.Value)
	}
	queryBuf.WriteByte(0)
	errMsg := db.setExec(setExecReq{
		k1:       k1,
		sql:      btoS(queryBuf.GetBytes()),
		valueBuf: valueBuf,
	})
	if errMsg != "" {
		panic(errMsg)
	}
	return
}
