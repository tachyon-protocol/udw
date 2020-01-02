// +build ios macAppStore

#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"
#import "udwAppStorePayHandler.h"
#import <objc/runtime.h>
#import "zzzig_udwRpcGoAndObjectiveC.h"

NSString* udwGetAppStoreReceiptData(){
    NSURL   *receiptURL  = [[NSBundle mainBundle] appStoreReceiptURL];
    NSData  *receiptData = [NSData dataWithContentsOfURL:receiptURL];
    return [receiptData base64EncodedStringWithOptions:0];
}

typedef void(^udwPayLoadingBlock)(NSString *loadType);
typedef void(^udwPayAlertBlock)(NSString *alertType);
typedef void(^udwPayHandlerBlock)(BOOL isSuccess);
typedef void(^udwPayHandlerProductIdBlock)(NSString *productId, NSString *transactionId);

@interface udwAppStorePayHandler : NSObject
+ (instancetype)shareHandler;
- (void)setPayLoadingCallback:(udwPayLoadingBlock)payLoadingCallback;
- (void)setPayAlertCallback:(udwPayAlertBlock)payAlertCallback;
- (void)SubscribeByProductId:(NSString *)productId callback:(udwPayHandlerProductIdBlock)callback;
- (void)UpgradeByProductId:(NSString *)productId callback:(udwPayHandlerProductIdBlock)callbac;
- (void)RestoreActionWithCallback:(udwPayHandlerBlock)callback;
- (void)setIapPromoteCallback:(BOOL (^)(NSString *))iapPromoteCallback;
@end

@interface udwAppStorePayHandler()<SKPaymentTransactionObserver, SKProductsRequestDelegate>
@property (nonatomic, strong) SKProductsRequest *payRequest;
@property (nonatomic, copy) NSString *selectProductId;
@property (nonatomic, assign) BOOL isUpgrade;
@property (nonatomic, assign) BOOL IsInitPay;                                            
                                         
@property (nonatomic,copy) udwPayLoadingBlock payLoadingCallback;
                  
@property (nonatomic,copy) udwPayAlertBlock payAlertCallback;
                                                                        
@property (nonatomic,copy) udwPayHandlerProductIdBlock PurchasedSuccessCallback;
                                                           
@property (nonatomic,copy) udwPayHandlerBlock RestoredSuccessCallback;
                   
@property (nonatomic,copy) BOOL(^IapPromoteCallback)(NSString *productId);
@end

static udwAppStorePayHandler *_handler = nil;

@implementation udwAppStorePayHandler
+ (instancetype)shareHandler {
    udwAppStorePayHandler *handler = [self shareInstance];
    return handler;
}

+ (instancetype)shareInstance {
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        _handler = [[super allocWithZone:NULL] init];
        [[SKPaymentQueue defaultQueue] addTransactionObserver:_handler];
    });
    return _handler;
}

+ (instancetype)allocWithZone:(struct _NSZone *)zone {
    return [udwAppStorePayHandler shareInstance];
}

                       
static void * udwPayLogingKey = &udwPayLogingKey;
- (void)setPayLoadingCallback:(udwPayLoadingBlock)payLoadingCallback{
    objc_setAssociatedObject(self, udwPayLogingKey, payLoadingCallback, OBJC_ASSOCIATION_COPY_NONATOMIC);
}
- (udwPayLoadingBlock)payLoadingCallback{
    return ^(NSString *loadingType){
        setAppStoreLoading(loadingType);
        udwPayLoadingBlock block=objc_getAssociatedObject(self, udwPayLogingKey);
        if (block) {
            block(loadingType);
        }
    };
}
static void * udwPayAlertKey = &udwPayAlertKey;
- (void)setPayAlertCallback:(udwPayAlertBlock)payAlertCallback{
    objc_setAssociatedObject(self, udwPayAlertKey, payAlertCallback, OBJC_ASSOCIATION_COPY_NONATOMIC);
}
- (udwPayAlertBlock)payAlertCallback{
    return ^(NSString *alertType){
        setAppStoreAlert(alertType);
        udwPayAlertBlock block=objc_getAssociatedObject(self, udwPayAlertKey);
        if (block) {
            block(alertType);
        }
    };
}
static void * udwPayPurchasedSuccessKey = &udwPayPurchasedSuccessKey;
- (void)setPurchasedSuccessCallback:(udwPayHandlerProductIdBlock)PurchasedSuccessCallback{
    objc_setAssociatedObject(self, udwPayPurchasedSuccessKey, PurchasedSuccessCallback, OBJC_ASSOCIATION_COPY_NONATOMIC);
}
- (udwPayHandlerProductIdBlock)PurchasedSuccessCallback{
    return ^(NSString *productId, NSString *transactionId){
        setAppStorePurchasedSuccess(productId,transactionId);
        udwPayHandlerProductIdBlock block=objc_getAssociatedObject(self, udwPayPurchasedSuccessKey);
        if (block) {
            block(productId,transactionId);
        }
    };
}
static void * udwPayRestoredSuccessKey = &udwPayRestoredSuccessKey;
- (void)setRestoredSuccessCallback:(udwPayHandlerBlock)RestoredSuccessCallback{
    objc_setAssociatedObject(self, udwPayRestoredSuccessKey, RestoredSuccessCallback, OBJC_ASSOCIATION_COPY_NONATOMIC);
}
- (udwPayHandlerBlock)RestoredSuccessCallback{
    return ^(BOOL isSuccess){
        setAppStoreRestoredSuccess(isSuccess);
        udwPayHandlerBlock block=objc_getAssociatedObject(self, udwPayRestoredSuccessKey);
        if (block) {
            block(isSuccess);
        }
    };
}

- (void)setIapPromoteCallback:(BOOL (^)(NSString *))iapPromoteCallback{
    if (iapPromoteCallback) {
        _IapPromoteCallback = iapPromoteCallback;
    }
}

                     
- (void)SubscribeByProductId:(NSString *)productId{
    [self SubscribeByProductId:productId callback:nil];
}
- (void)SubscribeByProductId:(NSString *)productId callback:(udwPayHandlerProductIdBlock)callback {
    self.isUpgrade = false;
    self.PurchasedSuccessCallback = callback;
    self.RestoredSuccessCallback = nil;
    [self inner_SubscribeByProductId:productId];
}
- (void)UpgradeByProductId:(NSString *)productId{
    [self UpgradeByProductId:productId callback:nil];
}
- (void)UpgradeByProductId:(NSString *)productId callback:(udwPayHandlerProductIdBlock)callback {
    self.isUpgrade = true;
    self.PurchasedSuccessCallback = callback;
    self.RestoredSuccessCallback = nil;
    [self inner_SubscribeByProductId:productId];
}
- (void)RestoreAction{
    [self RestoreActionWithCallback:nil];
}
               
- (void)RestoreActionWithCallback:(udwPayHandlerBlock)callback {
    self.IsInitPay = true;
    self.RestoredSuccessCallback = callback;
    self.PurchasedSuccessCallback = nil;
    [[SKPaymentQueue defaultQueue] restoreCompletedTransactions];
                             
    udwRunOnMainThread(^{
        if (self.payLoadingCallback) {
            self.payLoadingCallback(udwPayLoadingTypeRestore);
        }
    });
}

                      
- (void)inner_SubscribeByProductId:(NSString *)productId {
    self.IsInitPay = true;
    if ([SKPaymentQueue canMakePayments]) {
                                                           
        udwRunOnMainThread(^{
            if (self.payLoadingCallback) {
                self.payLoadingCallback(udwPayLoadingTypeSubscribing);
            }
        });
        [self getProductInfoWithProductId:productId];
        self.selectProductId = productId;
    } else {
                                                     
        udwRunOnMainThread(^{
            if (self.payAlertCallback) {
                self.payAlertCallback(udwPayAlertTypeCantMakePayment);
            }
        });
    }
}

- (void)getProductInfoWithProductId:(NSString *)productID {
    if (udwStringIsNilOrEmpty(productID)) {
        udwRunOnMainThread(^{
                                                                  
            if (self.payLoadingCallback) {
                self.payLoadingCallback(udwPayLoadingTypeHide);
            }
            if (self.payAlertCallback) {
                self.payAlertCallback(udwPayAlertTypeProductIdNil);
            }
        });
        return;
    }
    NSSet * set = [NSSet setWithArray:@[productID]];
    self.payRequest = [[SKProductsRequest alloc] initWithProductIdentifiers:set];
    self.payRequest.delegate = self;
    [self.payRequest start];
}

                              
- (void)productsRequest:(SKProductsRequest *)request didReceiveResponse:(SKProductsResponse *)response {
    self.payRequest.delegate = nil;
    NSArray *myProduct = response.products;
    if (myProduct.count == 0) {
        udwRunOnMainThread(^{
            if (self.payLoadingCallback) {
                self.payLoadingCallback(udwPayLoadingTypeHide);
            }
            if (self.payAlertCallback) {
                self.payAlertCallback(udwPayAlertTypeNotExistProduct);
            }
        });
        if (udwAppStoreOcDebug) {NSLog(@"Unable to get product information, purchase failed.");}
        return;
    }
    for (SKProduct *product in myProduct) {
        if (udwAppStoreOcDebug) {NSLog(@"Get product information and initiate a purchase.%@", product.productIdentifier);}
    }
    SKPayment * payment = [SKPayment paymentWithProduct:myProduct[0]];
    [[SKPaymentQueue defaultQueue] addPayment:payment];
}

- (void)paymentQueue:(SKPaymentQueue *)queue restoreCompletedTransactionsFailedWithError:(NSError *)error{
                   
    udwRunOnMainThread(^{
        if (self.payLoadingCallback) {
            self.payLoadingCallback(udwPayLoadingTypeHide);
        }
        if (self.RestoredSuccessCallback) {
            self.RestoredSuccessCallback(false);
        }
        if (self.payAlertCallback) {
            self.payAlertCallback(udwPayAlertTypeRestoreFailed);
        }
    });
}

                       
- (void)paymentQueue:(SKPaymentQueue *)queue updatedTransactions:(NSArray *)transactions {
    if (udwAppStoreOcDebug){NSLog(@"[Payment callback] len:%lu",(unsigned long)transactions.count);}
    for (SKPaymentTransaction *transaction in transactions)
    {
        switch (transaction.transactionState)
        {
            case SKPaymentTransactionStatePurchasing:
                if (udwAppStoreOcDebug){NSLog(@"[Payment callback] transactionId[%@] Purchasing",transaction.transactionIdentifier);}
                break;
            case SKPaymentTransactionStateDeferred:
                if (udwAppStoreOcDebug){NSLog(@"[Payment callback] transactionId[%@] Deferred",transaction.transactionIdentifier);}
                break;
            case SKPaymentTransactionStateFailed: {
                if (udwAppStoreOcDebug){NSLog(@"[Payment callback] transactionId[%@] error[%@] Failed",transaction.transactionIdentifier,transaction.error.localizedDescription );}
                [self failedTransaction:transaction];
                break;
            }
            case SKPaymentTransactionStateRestored: {
                if (udwAppStoreOcDebug){NSLog(@"[Payment callback] transactionId[%@] Restored",transaction.transactionIdentifier);}
                [self restoreTransaction:transaction];
                break;
            }
            case SKPaymentTransactionStatePurchased: {
                                                                 
                [[SKPaymentQueue defaultQueue] finishTransaction:transaction];
                if (udwAppStoreOcDebug){NSLog(@"[Payment callback] transactionId[%@] Purchased",transaction.transactionIdentifier);}
                if (!self.IsInitPay) {
                    return;
                }
                if (self.PurchasedSuccessCallback) {
                    self.PurchasedSuccessCallback(self.selectProductId, transaction.transactionIdentifier);
                }
                break;
            }
            default:
                break;
        }
    }
}

- (void)failedTransaction:(SKPaymentTransaction *)transaction {
                                      
    udwRunOnMainThread(^{
        if (self.payLoadingCallback) {
            self.payLoadingCallback(udwPayLoadingTypeHide);
        }
        if (self.payAlertCallback) {
            self.payAlertCallback(udwPayAlertTypeFailedTransaction);
        }
    });
    [[SKPaymentQueue defaultQueue] finishTransaction: transaction];
}

- (void)restoreTransaction:(SKPaymentTransaction *)transaction {
    if (udwAppStoreOcDebug){NSLog(@"Resume purchase");}
    [[SKPaymentQueue defaultQueue] finishTransaction: transaction];
}

- (void)paymentQueueRestoreCompletedTransactionsFinished:(SKPaymentQueue *)queue {
    if (self.RestoredSuccessCallback) {
        self.RestoredSuccessCallback(true);
    }
}

- (void)dealloc{
    [[SKPaymentQueue defaultQueue] removeTransactionObserver:self];
    if (self.payRequest) {
        self.payRequest.delegate = nil;
        [self.payRequest cancel];
        self.payRequest = nil;
    }
    if (udwAppStoreOcDebug){NSLog(@"---------------udwAppStorePayHandler>>>>>>>dealloc");}
}

- (void)cancelPayResuest{
    self.payRequest.delegate = nil;
    [self.payRequest cancel];
}

#if TARGET_OS_IPHONE
                         
- (BOOL)paymentQueue:(SKPaymentQueue *)queue shouldAddStorePayment:(SKPayment *)payment forProduct:(SKProduct *)product{
    if (self.IapPromoteCallback) {
        return self.IapPromoteCallback(product.productIdentifier);
    }
    return false;
}
#else
#endif

@end

void udwAppStorePaySetPayLoadingCallback(void (^callback)(NSString *loadType)){
    [[udwAppStorePayHandler shareHandler] setPayLoadingCallback:callback];
}
void udwAppStorePaySetPayAlertCallback(void (^callback)(NSString *alertType)){
    [[udwAppStorePayHandler shareHandler] setPayAlertCallback:callback];
}
void udwAppStorePaySubscribeCallback(NSString *productId,void (^successCallback)(NSString *productId, NSString *transactionId)){
    [[udwAppStorePayHandler shareHandler] SubscribeByProductId:productId callback:successCallback];
}
void udwAppStorePayUpgradeCallback(NSString *productId,void (^successCallback)(NSString *productId, NSString *transactionId)){
    [[udwAppStorePayHandler shareHandler] UpgradeByProductId:productId callback:successCallback];
}
void udwAppStorePayRestoreCallback(void (^callback)(BOOL success)){
    [[udwAppStorePayHandler shareHandler] RestoreActionWithCallback:callback];
}
           
void udwAppStoreIapPromoteCallback(BOOL (^iapPromoteCallback)(NSString *productId)){
    [[udwAppStorePayHandler shareHandler] setIapPromoteCallback:iapPromoteCallback];
}

void cForGoAppStorePaySubscribe(NSString *productId){
    [[udwAppStorePayHandler shareHandler] SubscribeByProductId:productId];
}
void cForGoAppStorePayUpgrade(NSString *productId){
    [[udwAppStorePayHandler shareHandler] UpgradeByProductId:productId];
}
void cForGoAppStorePayRestore(){
    [[udwAppStorePayHandler shareHandler] RestoreAction];
}

