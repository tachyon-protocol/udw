package udwSqlite3

import "github.com/tachyon-protocol/udw/udwStrconv"

func (db *Db) MustGetInt64(k1 string, k2 string) int64 {
	v, err := db.Get(k1, k2)
	if err != nil {
		panic(err)
	}
	if v == "" {
		return 0
	}
	return udwStrconv.MustParseInt64(v)
}

func (db *Db) MustSetInt64(k1 string, k2 string, v int64) {
	db.MustSet(k1, k2, udwStrconv.FormatInt64(v))
}
