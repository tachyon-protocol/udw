package udwGoBuildCtx

var gGetArchList []string = []string{
	`386`,
	`amd64`,
	`amd64p32`,
	`arm`,
	`arm64`,
	`mips`,
	`mips64`,
	`mips64le`,
	`mipsle`,
	`ppc64`,
	`ppc64le`,
	`riscv64`,
	`s390x`,
	`wasm`,
}

func GetArchList() []string {
	return gGetArchList
}
