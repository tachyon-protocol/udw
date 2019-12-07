package udwSqlite3

import (
	"bytes"
)

func (db *Db) MustDelete(k1 string, k2 string) {
	errMsg := db.querySkipResult(`DELETE FROM `+db.getTableNameFromK1(k1)+` WHERE k=?`, []byte(k2))
	if errMsg != "" {
		if errorIsTableNotExist(errMsg) {
			return
		}
		panic(errMsg)
	}
	return
}

func (db *Db) MustMulitDelete(k1 string, k2List []string) {

	for {
		if len(k2List) == 0 {
			return
		}
		toIndex := db.req.MulitDeleteMaxSize
		if toIndex > len(k2List) {
			toIndex = len(k2List)
		}
		db.mustMulitDeleteL1(k1, k2List[:toIndex])
		k2List = k2List[toIndex:]
	}
}

const defaultMulitDeleteMaxSize = 499

func (db *Db) mustMulitDeleteL1(k1 string, k2List []string) {
	sqlBuf := bytes.NewBufferString(`DELETE FROM ` + db.getTableNameFromK1(k1) + ` WHERE k IN (`)
	valueBuf := [][]byte{}
	l := len(k2List)
	for i, v := range k2List {
		if i < l-1 {
			sqlBuf.WriteString("?,")
		} else {
			sqlBuf.WriteString("?")
		}
		valueBuf = append(valueBuf, []byte(v))
	}
	sqlBuf.WriteString(")")
	db.mustSetExec(k1, sqlBuf.String(), valueBuf)
	return
}

func (db *Db) MustEmptyK1(k1 string) {
	db.mustSetExec(k1, `DELETE FROM `+db.getTableNameFromK1(k1), nil)
}

func (db *Db) IMustDelete(k1 string, k2 string) {
	db.MustDelete(k1, k2)
}

func (db *Db) MustDeleteWithKv(k1 string, k2 string, v string) {
	errMsg := db.querySkipResult(`DELETE FROM `+db.getTableNameFromK1(k1)+` WHERE k=? AND v=?`, []byte(k2), []byte(v))
	if errMsg != "" {
		if errorIsTableNotExist(errMsg) {
			return
		}
		panic(errMsg)
	}
	return
}
