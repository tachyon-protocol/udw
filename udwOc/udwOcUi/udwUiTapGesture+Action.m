#import "udwOcUi.h"
#import <objc/runtime.h>

static void * udwGestureRecognizerActionHandleKey = &udwGestureRecognizerActionHandleKey;
@implementation UITapGestureRecognizer (Action)
- (void)udwActionWithAction:(void (^)(void))action{
    objc_setAssociatedObject(self, udwGestureRecognizerActionHandleKey, action, OBJC_ASSOCIATION_COPY_NONATOMIC);
    [self addTarget:self action:@selector(udwTriggerActionHandleBlock)];
}

- (void)udwTriggerActionHandleBlock{
    void(^actionBlock)(void) = objc_getAssociatedObject(self, udwGestureRecognizerActionHandleKey);
    if (actionBlock) {
        actionBlock();
    }
}
@end