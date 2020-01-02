#import "udwOcUi.h"
#import <objc/runtime.h>

static void * udwControlActionHandleKey = &udwControlActionHandleKey;
@implementation UIControl (Action)
- (void)udwActionHandleEvents:(UIControlEvents)controlEvents action:(void (^)(void))action{
    objc_setAssociatedObject(self, udwControlActionHandleKey, action, OBJC_ASSOCIATION_COPY_NONATOMIC);
    [self addTarget:self action:@selector(udwTriggerActionHandleBlock) forControlEvents:controlEvents];
}

- (void)udwTriggerActionHandleBlock{
    void(^actionBlock)(void) = objc_getAssociatedObject(self, udwControlActionHandleKey);
    if (actionBlock) {
        actionBlock();
    }
}
@end