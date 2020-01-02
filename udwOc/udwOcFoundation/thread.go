// +build ios macAppStore

package udwOcFoundation

func IsMainThread() bool {
	return udwRunIsMainThread()
}
