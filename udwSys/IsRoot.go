package udwSys

import (
	"fmt"
	"os"
)

func MustIsRoot() bool {
	return mustIsRoot()
}

func MustIsRootOnCmd() {
	if !MustIsRoot() {
		fmt.Println("need root to run this command.")
		os.Exit(1)
	}
}

func MustRoot() {
	if !mustIsRoot() {
		panic("need root to run")
	}
}
