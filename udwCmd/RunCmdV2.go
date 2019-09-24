package udwCmd

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
)

type CmdV2Request struct {
	CmdSlice []string
}

type CmdV2Response struct {
	req     CmdV2Request
	execCmd *exec.Cmd
}

func (resp *CmdV2Response) MustRun() {
	resp.init()
	err := resp.execCmd.Run()
	if err != nil {
		panic(err)
	}
}

func (resp *CmdV2Response) MustRunAsync() {
	resp.init()
	err := resp.execCmd.Start()
	if err != nil {
		panic(err)
	}
}

func (resp *CmdV2Response) RunAndReturnOutputAsync() (br *bytes.Buffer, err error) {
	resp.init()
	buf := &bytes.Buffer{}
	w := io.MultiWriter(buf, os.Stdout)
	resp.execCmd.Stdout = w
	resp.execCmd.Stderr = w
	err = resp.execCmd.Start()
	return buf, err
}

func (resp *CmdV2Response) MustWait() {
	err := resp.Wait()
	if err != nil {
		panic(err)
	}
}

func (resp *CmdV2Response) Wait() error {
	execCmd := resp.execCmd
	if execCmd == nil {
		return errors.New("[CmdV2Response.MustWait] resp.execCmd==nil")
	}
	return execCmd.Wait()
}

func CmdV2(req CmdV2Request) *CmdV2Response {
	return &CmdV2Response{
		req: req,
	}
}
