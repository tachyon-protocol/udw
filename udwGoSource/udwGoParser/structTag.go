package udwGoParser

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoReader"
	"unicode"
	"unicode/utf8"
)

func (sf StructField) GetTagL1List() (outTagL1List []TagL1, errMsg string) {
	tag := sf.Tag
	if len(tag) == 0 {
		return nil, ""
	}
	reader := udwGoReader.NewReaderWithBuf([]byte(tag))
	for {
		reader.ReadAllSpace()
		keyB := reader.ReadUntilByte(':')
		if len(keyB) == 0 {

			return outTagL1List, ""
		}
		if keyB[len(keyB)-1] != ':' {
			return nil, "[StructField.GetTagL1List] unknow format 1 " + tag
		}
		keyB = keyB[:len(keyB)-1]
		var valueB []byte
		err := udwErr.PanicToError(func() {
			valueB = mustReadGoString(reader)
		})
		if err != nil {
			return nil, err.Error()
		}
		outTagL1List = append(outTagL1List, TagL1{
			Key:   string(keyB),
			Value: string(valueB),
		})
	}
}

func (sf StructField) GetTagL1ValueByKeyIgnoreError(key string) string {
	outTagL1List, errMsg := sf.GetTagL1List()
	if errMsg != "" {
		return ""
	}
	for _, tagL1 := range outTagL1List {
		if tagL1.Key == key {
			return tagL1.Value
		}
	}
	return ""
}

type TagL1 struct {
	Key   string
	Value string
}

func IsNameGoExport(name string) bool {
	if len(name) == 0 {
		return false
	}
	firstRune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(firstRune)
}
