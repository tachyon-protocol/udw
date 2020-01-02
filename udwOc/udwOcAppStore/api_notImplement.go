// +build !ios,!macAppStore

package udwOcAppStore

func SetPayLoadingCallback(f func(loadType string)) {
}

func SetPayAlertCallback(f func(alertType string)) {
}

func SetPayPurchasedSuccessCallback(f func(productId, transactionId string)) {
}

func SetPayRestoredCallback(f func(success bool)) {
}

func PaySubscribe(productId string) {
}

func PayUpgrade(productId string) {
}

func RestoreAction() {
}
func GetAppStoreReceiptData() string {
	return ""
}
