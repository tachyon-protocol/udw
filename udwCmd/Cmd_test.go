package udwCmd

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestExist(t *testing.T) {
	udwTest.Ok(!Exist("absadfklsjdfa"))
	udwTest.Ok(Exist("ls"))
	udwTest.Ok(Exist("top"))
}

func TestWriter(ot *testing.T) {
	buf := &udwBytes.BufWriter{}
	CmdBash(`echo 1`).MustRunWithWriter(buf)
	udwTest.Equal(buf.GetBytes(), []byte(`> bash -c echo 1
1
`))
	buf = &udwBytes.BufWriter{}
	CmdBash(`echo 1`).MustOutputToWriter(buf)
	udwTest.Equal(buf.GetBytes(), []byte(`1
`))
}
