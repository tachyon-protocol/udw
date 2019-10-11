package udwSqlite3

import (
	"encoding/json"
	"github.com/tachyon-protocol/udw/udwJson"
)

func (db *Db) GetJson(k1 string, k2 string, obj interface{}) (err error) {
	value, err := db.Get(k1, k2)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	err = json.Unmarshal([]byte(value), obj)
	if err != nil {
		return err
	}
	return nil
}

func (db *Db) MustGetJson(k1 string, k2 string, obj interface{}) {
	err := db.GetJson(k1, k2, obj)
	if err != nil {
		panic(err)
	}
}

func (db *Db) MustSetJson(k1 string, k2 string, obj interface{}) {
	b := udwJson.MustMarshal(obj)
	db.MustSet(k1, k2, string(b))
}
