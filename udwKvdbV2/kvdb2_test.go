package udwKvdbV2

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestNewKvdbV23(t *testing.T) {
	const dp = `/tmp/b466xkydj4`
	udwFile.MustDelete(dp)

	defer udwFile.MustDelete(dp)

	db := NewKvdbV2(NewKvdbV2Request{
		DataDirPath: dp,
	})
	db.loadFromList([]decryptItem{
		{
			data: []byte(`12345`),
		},
	})
	data := marshalMap(map[string]string{`1`: `233`})
	db.loadFromList([]decryptItem{
		{
			data: data,
		},
		{
			data: data,
		},
	})
	udwTest.Ok(db.Get(`1`) == `233`)
	db.Set(`3`, `2`)
	db.loadFromList([]decryptItem{
		{
			data: data,
		},
	})
	udwTest.Ok(db.Get(`3`) == `2`)

	udwErr.PanicIfError(db.FlushAndWaitDiskIO())

	udwFile.MustDelete(dp)
	udwFile.MustMkdir(dp)
	db.Set(`k2`, `v2`)
}
