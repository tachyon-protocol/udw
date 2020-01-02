package udwOc_Build

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwRpc/udwRpcGoAndObjectiveC/udwRpcGoAndObjectiveCBuilder"
	"github.com/tachyon-protocol/udw/udwStrings"
)

func DsnBuild() {
	udwFile.MustDeleteFile("src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwRpcGoAndObjectiveC__Gen.go")
	udwFile.MustDeleteFile("src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwRpcGoAndObjectiveC__Gen.h")
	udwFile.MustDeleteFile("src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwRpcGoAndObjectiveC__Gen.m")
	udwRpcGoAndObjectiveCBuilder.MustBuild(udwRpcGoAndObjectiveCBuilder.MustBuildRequest{
		PackagePath: "github.com/tachyon-protocol/udw/udwOc/udwOcFoundation",
		Filter: func(name string) bool {
			return udwStrings.IsInSlice([]string{}, name)
		},
		GoToOcFunctionList: []string{
			"func UdwOcGetAppName() string",
			"func udwRunIsMainThread() bool",
			"func UdwOcUserDefaultsKvGet(key string) string",
			"func UdwOcUserDefaultsKvSet(key string,value string)",
			"func UdwTimeGetTimeZoneName() string",
			"func UdwTimeGetTimeZoneOffset() int",
		},
	})
	udwRpcGoAndObjectiveCBuilder.MustBuild(udwRpcGoAndObjectiveCBuilder.MustBuildRequest{
		PackagePath: "github.com/tachyon-protocol/udw/udwOc/udwOcUi",
		Filter: func(name string) bool {
			return udwStrings.IsInSlice([]string{}, name)
		},
		GoToOcFunctionList: []string{
			"func UdwAlertOkV2(title string,message string)",
			"func KeepScreenOn()",
			"func UdwLogAlert(msg string)",
			"func UdwLogAlertClean()",
			"func UdwAddLocalNotification(date time.Time,title string,content string,id string)",
			"func UdwCancelAllLocalNotification()",
			"func UdwCancelOneLocalNotification(identifier string)",
		},
	})
	udwRpcGoAndObjectiveCBuilder.MustBuild(udwRpcGoAndObjectiveCBuilder.MustBuildRequest{
		PackagePath:      "github.com/tachyon-protocol/udw/udwOc/udwOcAppStore",
		BuildFlagContent: "ios macAppStore",
		Filter: func(name string) bool {
			return udwStrings.IsInSlice([]string{
				`setAppStoreLoading`,
				`setAppStoreAlert`,
				`setAppStorePurchasedSuccess`,
				`setAppStoreRestoredSuccess`,
			}, name)
		},
		GoToOcFunctionList: []string{
			"func cForGoAppStorePaySubscribe(productId string)",
			"func cForGoAppStorePayUpgrade(productId string)",
			"func cForGoAppStorePayRestore()",
			"func udwGetAppStoreReceiptData() string",
		},
	})
	udwRpcGoAndObjectiveCBuilder.MustBuild(udwRpcGoAndObjectiveCBuilder.MustBuildRequest{
		PackagePath:      "github.com/tachyon-protocol/udw/udwOc/udwOcMac",
		BuildFlagContent: "darwin,!ios",
		Filter: func(name string) bool {
			return udwStrings.IsInSlice([]string{
				"applicationDidBecomeActiveCall",
			}, name)
		},
		GoToOcFunctionList: []string{
			"func udwShowLocalNotification(context string)",
			"func udwOpenUrlWithDefaultBrowser(url string)",
			"func udwCheckIsRunningApp(bundleId string) bool",
			"func udwGetProcessIdByName(processName string) int",
		},
	})
}
