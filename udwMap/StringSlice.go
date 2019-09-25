package udwMap

func StringSliceDeleteString(ss []string, sWillDelete string) []string {
	after := []string{}
	for _, str := range ss {
		if str != sWillDelete {
			after = append(after, str)
		}
	}
	return after
}
