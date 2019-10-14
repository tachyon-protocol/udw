package udwGoTypeMarshal

import (
	"github.com/tachyon-protocol/udw/udwMap"
)

func MapStringStringToSortedStringKeyValueList(m map[string]string) []udwMap.KeyValuePair {
	return udwMap.MapStringStringToKeyValuePairListAes(m)
}
