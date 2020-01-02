#import "udwOcUi.h"
#import <objc/runtime.h>

static void * udwScrollViewTouchStrikeEnabledKey = &udwScrollViewTouchStrikeEnabledKey;
@implementation UIScrollView (udwTouchEvent)
@dynamic udwScrollViewTouchStrikeEnabled;
- (void)touchesBegan:(NSSet *)touches withEvent:(UIEvent *)event {
    if (self.udwScrollViewTouchStrikeEnabled) {
        [[self nextResponder] touchesBegan:touches withEvent:event];
        [super touchesBegan:touches withEvent:event];
    }
}

-(void)touchesMoved:(NSSet *)touches withEvent:(UIEvent *)event {
    if (self.udwScrollViewTouchStrikeEnabled) {
        [[self nextResponder] touchesMoved:touches withEvent:event];
        [super touchesMoved:touches withEvent:event];
    }
}

- (void)touchesEnded:(NSSet *)touches withEvent:(UIEvent *)event {
    if (self.udwScrollViewTouchStrikeEnabled) {
        [[self nextResponder] touchesEnded:touches withEvent:event];
        [super touchesEnded:touches withEvent:event];
    }
}

- (void)setUdwScrollViewTouchStrikeEnabled:(BOOL)udwScrollViewTouchStrikeEnabled{
    objc_setAssociatedObject(self, udwScrollViewTouchStrikeEnabledKey, @(udwScrollViewTouchStrikeEnabled), OBJC_ASSOCIATION_ASSIGN);
}

- (BOOL)udwScrollViewTouchStrikeEnabled{
    return [objc_getAssociatedObject(self, udwScrollViewTouchStrikeEnabledKey) boolValue];
}

@end