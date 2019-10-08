package udwTime

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"strconv"
	"testing"
	"time"
)

func TestTtlStringTimeMap(ot *testing.T) {
	m := TtlStringTimeMapNew(time.Minute)
	has := m.Add("1", MustFromMysqlFormatDefaultTZ("2001-01-01 01:01:01"), MustFromMysqlFormatDefaultTZ("2001-01-01 01:01:01"))
	udwTest.Equal(has, false)

	has = m.Add("1", MustFromMysqlFormatDefaultTZ("2001-01-01 01:01:01"), MustFromMysqlFormatDefaultTZ("2001-01-01 01:01:01"))
	udwTest.Equal(has, true)

	has = m.Add("1", MustFromMysqlFormatDefaultTZ("2001-01-01 01:03:01"), MustFromMysqlFormatDefaultTZ("2001-01-01 01:03:01"))
	udwTest.Equal(has, false)
}

func TestTtlStringTimeMapGc(ot *testing.T) {
	m := TtlStringTimeMapNew(time.Minute)
	for i := 0; i < 1000; i++ {
		has := m.Add("a"+strconv.Itoa(i), MustFromMysqlFormatDefaultTZ("2001-01-01 01:01:01"), MustFromMysqlFormatDefaultTZ("2001-01-01 01:01:01"))
		udwTest.Equal(has, false)
	}
	udwTest.Equal(len(m.m), 1000)
	has := m.Add("1", MustFromMysqlFormatDefaultTZ("2001-01-01 01:02:02"), MustFromMysqlFormatDefaultTZ("2001-01-01 01:02:02"))
	udwTest.Equal(has, false)
	udwTest.Equal(len(m.m), 1)
}
