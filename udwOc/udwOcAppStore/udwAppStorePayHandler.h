// +build ios macAppStore

#import <Foundation/Foundation.h>
#import <StoreKit/StoreKit.h>

NSString* udwGetAppStoreReceiptData(void);

#define udwPayLoadingTypeRestore            @"Restore"
#define udwPayLoadingTypeSubscribing        @"Subscribing"
#define udwPayLoadingTypeHide               @"Hide"

#define udwPayAlertTypeCantMakePayment      @"CantMakePayment"
#define udwPayAlertTypeProductIdNil         @"ProductIdNil"
#define udwPayAlertTypeNotExistProduct      @"NotExistProduct"
#define udwPayAlertTypeRestoreFailed        @"RestoreFailed"
#define udwPayAlertTypeFailedTransaction    @"FailedTransaction"
                 
                                                           
void udwAppStorePaySetPayLoadingCallback(void (^callback)(NSString *loadType));
                                                               
void udwAppStorePaySetPayAlertCallback(void (^callback)(NSString *alertType));
                                                                     
void udwAppStorePaySubscribeCallback(NSString *productId,void (^successCallback)(NSString *productId, NSString *transactionId));
                                                                     
void udwAppStorePayUpgradeCallback(NSString *productId,void (^successCallback)(NSString *productId, NSString *transactionId));
                                                                      
void udwAppStorePayRestoreCallback(void (^callback)(BOOL success));
                            
void udwAppStoreIapPromoteCallback(BOOL (^iapPromoteCallback)(NSString *productId));
                     
void cForGoAppStorePaySubscribe(NSString *productId);
void cForGoAppStorePayUpgrade(NSString *productId);
void cForGoAppStorePayRestore(void);

#define udwAppStoreOcDebug true
