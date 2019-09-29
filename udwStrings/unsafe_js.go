// +build js

package udwStrings

func GetStringFromByteArrayNoAlloc(b []byte) string {
	return string(b)
}

func GetByteArrayFromStringNoAlloc(s string) []byte {
	return []byte(s)
}
