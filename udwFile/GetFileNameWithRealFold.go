package udwFile

func MustIsFileOrDirectoryNameFoldCorrect(path string) bool {
	return MustGetFileOrDirectoryNameWithRealFold(path) == MustGetFullPath(path)
}
