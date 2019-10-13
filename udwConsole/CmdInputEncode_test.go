package udwConsole

import (
	"github.com/tachyon-protocol/udw/udwHex"
	"github.com/tachyon-protocol/udw/udwTest"
	"strings"
	"testing"
	"time"
)

type TestRequest struct {
	L              string
	ListenAddrList []string
	Profile        bool
	L1             string `CmdFlag:"l"`
	InPort         int
	InDur          time.Duration
	Fn             func()
}

type TestFailRequest struct {
	L2 string `CmdFlag:" l"`
}

func mustRunCommandLineFromFuncV3L1(argList []string, fn interface{}) {
	errMsg := RunCommandLineFromFuncV3(RunCommandLineFromFuncV3Request{
		OsArgList: argList,
		F:         fn,
	})
	if errMsg != "" {
		panic(errMsg)
	}
}

func TestMustRunCommandLineFromFunc(t *testing.T) {
	udwTest.AssertPanicWithErrorMessage(func() {
		mustRunCommandLineFromFuncV3L1(nil, func(req TestRequest) (err error) {
			return nil
		})
	}, "t.NumOut()")

	hasRun := false
	mustRunCommandLineFromFuncV3L1([]string{"udw", "-L", ":40001"}, func(req TestRequest) {
		udwTest.Equal(req.L, ":40001")
		udwTest.Equal(len(req.ListenAddrList), 0)
		hasRun = true
	})
	udwTest.Equal(hasRun, true)

	hasRun = false
	mustRunCommandLineFromFuncV3L1([]string{"udw", "-L", ":40002", "-ListenAddrList", `[":40003",":40004"]`},
		func(req TestRequest) {
			udwTest.Equal(req.L, ":40002")
			udwTest.Equal(req.ListenAddrList, []string{":40003", ":40004"})
			hasRun = true
		})
	udwTest.Equal(hasRun, true)

	hasRun = false
	mustRunCommandLineFromFuncV3L1([]string{"udw", "-jsonHex", udwHex.EncodeStringToString(`{"L":":40003","ListenAddrList":[":40005",":40006"]}`)},
		func(req TestRequest) {
			udwTest.Equal(req.L, ":40003")
			udwTest.Equal(req.ListenAddrList, []string{":40005", ":40006"})
			hasRun = true
		})
	udwTest.Equal(hasRun, true)

	hasRun = false
	mustRunCommandLineFromFuncV3L1([]string{"udw", "-Profile"}, func(req TestRequest) {
		udwTest.Equal(req.Profile, true)
		hasRun = true
	})
	udwTest.Equal(hasRun, true)

	hasRun = false
	mustRunCommandLineFromFuncV3L1([]string{"udw", "-l", ":40001"}, func(req TestRequest) {
		udwTest.Equal(req.L1, ":40001")
		hasRun = true
	})
	udwTest.Equal(hasRun, true)

	hasRun = false
	mustRunCommandLineFromFuncV3L1([]string{"udw", "-l", ":40001"}, func(req *TestRequest) {
		udwTest.Equal(req.L1, ":40001")
		hasRun = true
	})
	udwTest.Equal(hasRun, true)

	hasRun = false
	mustRunCommandLineFromFuncV3L1(nil, func() {
		hasRun = true
	})
	udwTest.Equal(hasRun, true)

	hasRun = false
	udwTest.AssertPanicWithErrorMessage(func() {
		mustRunCommandLineFromFuncV3L1([]string{"udw", "- l", ":40001"}, func(req *TestFailRequest) {
			hasRun = true
		})
	}, "has not support char")
	udwTest.Equal(hasRun, false)

	hasRun = false
	mustRunCommandLineFromFuncV3L1([]string{"udw", ":40001"}, func(req *TestRequest) {
		udwTest.Equal(req.L, ":40001")
		hasRun = true
	})
	udwTest.Equal(hasRun, true)

	hasRun = false
	mustRunCommandLineFromFuncV3L1([]string{"udw", "-notExistFlag"}, func(req *TestRequest) {
		udwTest.Equal(req.L, "")
		hasRun = true
	})
	udwTest.Equal(hasRun, true)

	hasRun = false
	udwTest.AssertPanicWithErrorMessage(func() {
		mustRunCommandLineFromFuncV3L1([]string{"udw", "-h"}, func(req TestRequest) {
			hasRun = true
		})
	}, "Usage of udw:")
	udwTest.Equal(hasRun, false)

	hasRun = false
	mustRunCommandLineFromFuncV3L1([]string{"udw", "-InPort", "40001"}, func(req TestRequest) {
		udwTest.Equal(req.InPort, 40001)
		udwTest.Equal(len(req.ListenAddrList), 0)
		hasRun = true
	})
	udwTest.Equal(hasRun, true)

	udwTest.Ok(strings.Contains(GetUsageString(), "Usage of udw:"))

	hasRun = false
	mustRunCommandLineFromFuncV3L1([]string{"udw", "--L=:40001"}, func(req TestRequest) {
		udwTest.Equal(req.L, ":40001")
		udwTest.Equal(len(req.ListenAddrList), 0)
		hasRun = true
	})
	udwTest.Equal(hasRun, true)

	hasRun = false
	mustRunCommandLineFromFuncV3L1([]string{"udw", "-ListenAddrList", `[":40003"]`},
		func(req TestRequest) {
			udwTest.Equal(req.ListenAddrList, []string{":40003"})
			hasRun = true
		})
	udwTest.Equal(hasRun, true)

	type TestStepRequest struct {
		IpList []string
	}
	hasRun = false
	mustRunCommandLineFromFuncV3L1([]string{"udw", `["127.0.0.1"]`},
		func(req TestStepRequest) {
			udwTest.Equal(req.IpList, []string{"127.0.0.1"})
			hasRun = true
		})
	udwTest.Equal(hasRun, true)

	type TestStepRequest2 struct {
		UseFreeId bool
		Ip        string
	}
	hasRun = false
	mustRunCommandLineFromFuncV3L1([]string{"udw", "-UseFreeId", "-Ip", "35.189.172.135"},
		func(req TestStepRequest2) {
			udwTest.Equal(req.UseFreeId, true)
			udwTest.Equal(req.Ip, "35.189.172.135")
			hasRun = true
		})
	udwTest.Equal(hasRun, true)
}
