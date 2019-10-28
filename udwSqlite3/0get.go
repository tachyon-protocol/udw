package udwSqlite3

import (
	"bytes"
	"errors"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwMap"
	"sort"
)

func (db *Db) MustGet(k1 string, k2 string) string {
	v, err := db.Get(k1, k2)
	if err != nil {
		panic(err)
	}
	return v
}

func (db *Db) Get(k1 string, k2 string) (value string, err error) {

	queryBuf := &udwBytes.BufWriter{}
	queryBuf.Reset()
	queryBuf.WriteString(`SELECT v FROM `)
	queryBuf.WriteString(db.getTableNameFromK1(k1))
	queryBuf.WriteString(` WHERE k = ? LIMIT 1;`)
	queryBuf.WriteByte(0)
	argumentList := db.getArgumentList(1)
	argumentList[0] = stoB(k2)
	s, errMsg := db.queryToOneString(btoS(queryBuf.GetBytes()), argumentList...)

	if errMsg != "" {
		if errorIsTableNotExist(errMsg) {
			return "", nil
		}
		return "", errors.New(errMsg)
	}
	return s, nil
}

func (db *Db) MustExist(k1 string, k2 string) bool {
	queryBuf := &udwBytes.BufWriter{}
	queryBuf.WriteString(`SELECT '1' FROM `)
	queryBuf.WriteString(db.getTableNameFromK1(k1))
	queryBuf.WriteString(` WHERE k = ? LIMIT 1;`)
	queryBuf.WriteByte(0)
	argumentList := db.getArgumentList(1)
	argumentList[0] = stoB(k2)
	s, errMsg := db.queryToOneString(btoS(queryBuf.GetBytes()), argumentList...)
	if errMsg != "" {
		if errorIsTableNotExist(errMsg) {
			return false
		}
		panic(errMsg)
	}
	return s == "1"
}

func (db *Db) MustMulitGet(k1 string, k2List []string) (output []udwMap.KeyValuePair) {
	sort.Strings(k2List)

	for {
		if len(k2List) == 0 {
			return output
		}
		toIndex := mustMulitGetMaxSize
		if toIndex > len(k2List) {
			toIndex = len(k2List)
		}
		output = append(output, db.mustMulitGetL1(k1, k2List[:toIndex])...)
		k2List = k2List[toIndex:]
	}
}

func (db *Db) mustMulitGetL1(k1 string, k2List []string) []udwMap.KeyValuePair {
	queryBuf := &udwBytes.BufWriter{}
	queryBuf.WriteString(`SELECT k,v FROM `)
	queryBuf.WriteString(db.getTableNameFromK1(k1))
	queryBuf.WriteString(` WHERE k in (`)
	valueBuf := db.getArgumentList(len(k2List))
	l := len(k2List)
	for i, k := range k2List {
		if i < l-1 {
			queryBuf.WriteString("?,")
		} else {
			queryBuf.WriteString("?)")
		}
		valueBuf[i] = stoB(k)
	}
	output := []udwMap.KeyValuePair{}
	errMsg := db.Query(QueryReq{
		Query: btoS(queryBuf.GetBytes()),
		Args:  valueBuf,
		RespDataCb: func(row [][]byte) {
			output = append(output, udwMap.KeyValuePair{
				Key:   string(row[0]),
				Value: string(row[1]),
			})
		},
		ColumnCount: 2,
	})
	if errMsg != "" {
		if errorIsTableNotExist(errMsg) {
			return nil
		}
		panic(errMsg)
	}
	return output
}

func (db *Db) MustMulitExist(k1 string, k2List []string) (output map[string]struct{}) {

	output = map[string]struct{}{}
	for {
		if len(k2List) == 0 {
			return output
		}
		toIndex := mustMulitGetMaxSize
		if toIndex > len(k2List) {
			toIndex = len(k2List)
		}
		db.mustMulitExistL1(k1, k2List[:toIndex], output)
		k2List = k2List[toIndex:]
	}
}

func (db *Db) mustMulitExistL1(k1 string, k2List []string, output map[string]struct{}) {
	sqlBuf := bytes.NewBufferString(`SELECT k FROM ` + db.getTableNameFromK1(k1) +
		` WHERE k in (`)
	valueBuf := [][]byte{}
	l := len(k2List)
	for i, k := range k2List {
		if i < l-1 {
			sqlBuf.WriteString("?,")
		} else {
			sqlBuf.WriteString("?)")
		}
		valueBuf = append(valueBuf, []byte(k))
	}
	errMsg := db.Query(QueryReq{
		Query: sqlBuf.String(),
		Args:  valueBuf,
		RespDataCb: func(row [][]byte) {
			output[string(row[0])] = struct{}{}
		},
		ColumnCount: 1,
	})
	if errMsg != "" {
		if errorIsTableNotExist(errMsg) {
			return
		}
		panic(errMsg)
	}
	return
}
