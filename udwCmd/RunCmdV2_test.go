package udwCmd

import (
	"fmt"
	"testing"
	"time"
)

func TestCmdV2Request(t *testing.T) {
	resp := CmdV2(CmdV2Request{
		CmdSlice: []string{"bash", "-c", "top"},
	})
	buf, err := resp.RunAndReturnOutputAsync()
	if err != nil {
		panic("0" + err.Error())
	}
	go func() {
		time.Sleep(time.Second * 3)
		resp.MustKill()
	}()
	err = resp.Wait()
	fmt.Println("err:", err, "log:", buf.String())
}
