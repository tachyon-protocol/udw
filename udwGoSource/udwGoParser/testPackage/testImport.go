package testPackage

func ImportTester() string {
	return `import (
		"not_exist_package"
	)` + "import ()"
}
