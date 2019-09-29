package udwFile

func MustGetFileOrDirectoryNameWithRealFold(path string) string {
	return MustGetFullPath(path)
}

func MustMoveNameFoldCorrect(fromPath string, toPath string) {
	MustMove(fromPath, toPath)
}
