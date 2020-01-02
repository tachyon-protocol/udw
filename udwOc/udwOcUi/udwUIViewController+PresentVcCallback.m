#import "udwOcUi.h"
#import <objc/message.h>
#import <objc/runtime.h>

static void * udwPresentVcCallbackHandleKey = &udwPresentVcCallbackHandleKey;
@implementation UIViewController (PresentVcCallback)
@dynamic udwPresentVcCallback;
- (void)udw_presentViewController:(UIViewController *)viewControllerToPresent animated:(BOOL)flag completion:(void (^)(void))completion{
    if (self && viewControllerToPresent) {
        [self udw_presentViewController:viewControllerToPresent animated:flag completion:^{
            if (completion) {
                completion();
            }
            if (self.udwPresentVcCallback) {
                self.udwPresentVcCallback(viewControllerToPresent);
            }
        }];
    }
}

- (void (^)(UIViewController *))udwPresentVcCallback{
    return objc_getAssociatedObject(self, udwPresentVcCallbackHandleKey);
}

- (void)setUdwPresentVcCallback:(void (^)(UIViewController *))udwPresentVcCallback{
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        Method udw_presentViewController = class_getInstanceMethod([self class], @selector(udw_presentViewController:animated:completion:));
        Method presentViewController = class_getInstanceMethod([self class], @selector(presentViewController:animated:completion:));
        method_exchangeImplementations(udw_presentViewController,presentViewController);
    });
    objc_setAssociatedObject(self, udwPresentVcCallbackHandleKey, udwPresentVcCallback, OBJC_ASSOCIATION_COPY_NONATOMIC);
}

@end
