package udwRand

func GetRandOneInStringSlice(sList []string) string {
	return sList[Intn(len(sList))]
}

func ChooseOneInStringSlice(sList []string) string {
	return GetRandOneInStringSlice(sList)
}
