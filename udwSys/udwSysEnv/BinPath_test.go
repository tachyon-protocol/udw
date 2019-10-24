package udwSysEnv

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"os"
	"testing"
)

func TestRecoverPath0(t *testing.T) {
	target := "/abc:/dce"
	os.Setenv("PATH", target)
	RecoverPath()
	udwTest.Ok(os.Getenv("PATH") == "/usr/local/bin:/bin:/usr/bin:/abc:/dce", os.Getenv("PATH"))
}

func TestRecoverPath1(t *testing.T) {
	target := "/abc:/dce:/abc"
	os.Setenv("PATH", target)
	RecoverPath()
	udwTest.Ok(os.Getenv("PATH") == "/usr/local/bin:/bin:/usr/bin:/abc:/dce", os.Getenv("PATH"))
}

func TestRecoverPath2(t *testing.T) {
	target := "/bin:/abc:/dce:/abc"
	os.Setenv("PATH", target)
	RecoverPath()
	udwTest.Ok(os.Getenv("PATH") == "/usr/local/bin:/bin:/usr/bin:/abc:/dce", os.Getenv("PATH"))
}

func TestRecoverPath3(t *testing.T) {
	target := "/bin:/usr/local/bin:/abc:/dce:/abc"
	os.Setenv("PATH", target)
	RecoverPath()
	udwTest.Ok(os.Getenv("PATH") == "/usr/local/bin:/bin:/usr/bin:/abc:/dce", os.Getenv("PATH"))
}
