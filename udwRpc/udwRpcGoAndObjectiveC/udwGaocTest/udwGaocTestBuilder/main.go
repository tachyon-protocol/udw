package udwGaocTestBuilder

import (
	"github.com/tachyon-protocol/udw/udwRpc/udwRpcGoAndObjectiveC/udwRpcGoAndObjectiveCBuilder"
	"strings"
)

func MustBuild() {
	udwRpcGoAndObjectiveCBuilder.MustBuildNoCache(udwRpcGoAndObjectiveCBuilder.MustBuildRequest{
		PackagePath: "github.com/tachyon-protocol/udw/udwRpc/udwRpcGoAndObjectiveC/udwGaocTest",
		Filter: func(name string) bool {
			if strings.HasPrefix(name, "test") {
				return false
			}
			if strings.HasPrefix(name, "Test") {
				return false
			}
			if strings.HasPrefix(name, "goToOc") {
				return false
			}
			return true
		},
		GoToOcFunctionList: []string{
			"func goToOcVoid()",
			"func goToOcHasGoToOcVoidRun()bool",
			"func goToOcInAndOutStringSlice(inSlice []string)(outSlice []string)",
			"func goToOcInNoOut(msg string,msg2 string)",
			"func goToOcInt64Test()bool",
			"func goToOcTestType()bool",
			"func goToOcTestGuSEqual()int",
		},
	})
}
