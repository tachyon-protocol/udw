package udwGoParser

func IsTypeEqual(t1 Type, t2 Type) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if (t1 == nil && t2 != nil) || (t1 != nil && t2 == nil) {
		return false
	}
	return t1.Equal(t2)
}
