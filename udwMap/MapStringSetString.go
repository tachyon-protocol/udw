package udwMap

func MapStringSetStringAdd(m map[string]map[string]struct{}, k1 string, k2 string) {
	k2M := m[k1]
	if k2M == nil {
		k2M = map[string]struct{}{}
		m[k1] = k2M
	}
	k2M[k2] = struct{}{}
}
