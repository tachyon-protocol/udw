// +build ios macAppStore

package udwOcFoundation

func UserDefaultsSet(key string, value string) {
	UdwOcUserDefaultsKvSet(key, value)
}

func UserDefaultsGet(key string) string {
	return UdwOcUserDefaultsKvGet(key)
}
