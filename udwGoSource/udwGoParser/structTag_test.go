package udwGoParser

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"strings"
	"testing"
)

func TestStructTag(ot *testing.T) {
	sf := StructField{
		Tag: "",
	}
	tagl1, errMsg := sf.GetTagL1List()
	udwTest.Equal(errMsg, "")
	udwTest.Equal(len(tagl1), 0)

	sf.Tag = "bad stuff"
	tagl1, errMsg = sf.GetTagL1List()
	udwTest.Ok(strings.Contains(errMsg, "unknow format 1"))
	udwTest.Equal(len(tagl1), 0)

	sf.Tag = `json:"name" xml:"name"`
	tagl1, errMsg = sf.GetTagL1List()
	udwTest.Equal(errMsg, "")
	udwTest.Equal(len(tagl1), 2)
	udwTest.Equal(tagl1[0].Key, "json")
	udwTest.Equal(tagl1[0].Value, "name")
	udwTest.Equal(tagl1[1].Key, "xml")
	udwTest.Equal(tagl1[1].Value, "name")

	sf.Tag = `json:"name,abc"`
	tagl1, errMsg = sf.GetTagL1List()
	udwTest.Equal(errMsg, "")
	udwTest.Equal(len(tagl1), 1)
	udwTest.Equal(tagl1[0].Key, "json")
	udwTest.Equal(tagl1[0].Value, "name,abc")
}
