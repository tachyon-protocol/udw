package udwGoBuildCtx

var gGetOsList []string = []string{
	`android`,
	`darwin`,
	`dragonfly`,
	`freebsd`,
	`js`,
	`linux`,
	`nacl`,
	`netbsd`,
	`openbsd`,
	`plan9`,
	`solaris`,
	`windows`,
}

func GetOsList() []string {
	return gGetOsList
}
