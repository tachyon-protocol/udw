package udwKvdbV2

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestHistory(t *testing.T) {
	var history KvdbWriteHistory
	history.switchToSnapshot()
	history.recordCmdSet(nil, `k`, `v`)
	udwTest.Ok(history.IsSnapshot && history.getLength() == 0)
}

func TestNewKvdbV25(t *testing.T) {
	data := marshalMap(nil)
	m, ok := unmarshalMap(data)
	udwTest.Ok(ok && len(m) == 0)

	data = marshalMap(map[string]string{
		`k`: `v`,
	})
	m, ok = unmarshalMap(data)
	udwTest.Ok(ok && len(m) == 1 && m[`k`] == `v`)
}
