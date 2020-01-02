// +build ios macAppStore

package udwOcAppStore

func SetPayLoadingCallback(f func(loadType string)) {
	initContextIfNeed()
	gContext.loadingCallback = f
}

func SetPayAlertCallback(f func(alertType string)) {
	initContextIfNeed()
	gContext.alertCallback = f
}

func SetPayPurchasedSuccessCallback(f func(productId, transactionId string)) {
	initContextIfNeed()
	gContext.purchasedSuccessCallback = f
}

func SetPayRestoredCallback(f func(success bool)) {
	initContextIfNeed()
	gContext.restoredCallback = f
}

func PaySubscribe(productId string) {
	initContextIfNeed()
	cForGoAppStorePaySubscribe(productId)
}

func PayUpgrade(productId string) {
	initContextIfNeed()
	cForGoAppStorePayUpgrade(productId)
}

func RestoreAction() {
	initContextIfNeed()
	cForGoAppStorePayRestore()
}
func GetAppStoreReceiptData() string {
	initContextIfNeed()
	return udwGetAppStoreReceiptData()
}
