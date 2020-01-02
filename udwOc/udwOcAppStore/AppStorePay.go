// +build ios macAppStore

package udwOcAppStore

func setAppStoreLoading(loadType string) {
	if gContext == nil {
		return
	}
	if gContext.loadingCallback != nil {
		gContext.loadingCallback(loadType)
	}
}

func setAppStoreAlert(alertType string) {
	if gContext == nil {
		return
	}
	if gContext.alertCallback != nil {
		gContext.alertCallback(alertType)
	}
}

func setAppStorePurchasedSuccess(productId, transactionId string) {
	if gContext == nil {
		return
	}
	if gContext.purchasedSuccessCallback != nil {
		gContext.purchasedSuccessCallback(productId, transactionId)
	}
}

func setAppStoreRestoredSuccess(success bool) {
	if gContext == nil {
		return
	}
	if gContext.restoredCallback != nil {
		gContext.restoredCallback(success)
	}
}

type callbackContext struct {
	loadingCallback          func(loadType string)
	alertCallback            func(alertType string)
	purchasedSuccessCallback func(productId, transactionId string)
	restoredCallback         func(success bool)
}

func initContextIfNeed() {
	if gContext == nil {
		gContext = &callbackContext{}
	}
}

var gContext *callbackContext
