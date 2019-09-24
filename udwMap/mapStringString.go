package udwMap

type KeyValuePair struct {
	Key   string
	Value string
}

func MapStringStringToKeyValuePairListAes(m map[string]string) []KeyValuePair {
	list := make([]KeyValuePair, len(m))
	listIndex := 0
	for k, v := range m {
		list[listIndex] = KeyValuePair{
			Key:   k,
			Value: v,
		}
		listIndex++
	}
	SortKeyValuePairList(list)

	return list
}

func GetKeyListFromPairList(pairList []KeyValuePair) []string {
	keyList := []string{}
	for _, pair := range pairList {
		keyList = append(keyList, pair.Key)
	}
	return keyList
}

func MapStringStringGetKeyList(m map[string]string) []string {
	list := make([]string, len(m))
	listIndex := 0
	for k := range m {
		list[listIndex] = k
		listIndex++
	}
	return list
}

func MapStringStringGetValueList(m map[string]string) []string {
	list := make([]string, 0, len(m))
	for _, v := range m {
		list = append(list, v)
	}
	return list
}

type KvStringBytePair struct {
	K string
	V []byte
}

func MapStringStringValueSet(m map[string]string) map[string]struct{} {
	outMap := make(map[string]struct{}, len(m))
	for _, v := range m {
		outMap[v] = struct{}{}
	}
	return outMap
}

func MapStringStringClone(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	outMap := make(map[string]string, len(m))
	for k, v := range m {
		outMap[k] = v
	}
	return outMap
}

func MapStringStringCloneAndSet(m map[string]string, kvList ...string) map[string]string {
	outMap := MapStringStringClone(m)
	if len(kvList)%2 != 0 {
		panic("[MapStringStringCloneAndSet] must passing key value pair")
	}
	for i := 0; i < len(kvList); i += 2 {
		outMap[kvList[i]] = kvList[i+1]
	}
	return outMap
}

type KeyValueSlicePair struct {
	K []byte
	V []byte
}
