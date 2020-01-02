#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"

@implementation kouPointInsideButton

- (BOOL)pointInside:(CGPoint)point withEvent:(UIEvent *)event {
    CGRect bounds = self.bounds;
    CGFloat widthDelta = 60.0 - bounds.size.width;
    CGFloat heightDelta = 60.0 - bounds.size.height;
    if (((bounds.size.width > bounds.size.height) && widthDelta > 0) || ((bounds.size.width < bounds.size.height) && heightDelta > 0) || ((bounds.size.width == bounds.size.height) && widthDelta > 0)) {
         bounds = CGRectInset(bounds, -0.5 * widthDelta, -0.5 * heightDelta);
    }
    return CGRectContainsPoint(bounds, point);
}

@end
