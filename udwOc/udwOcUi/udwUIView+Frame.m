#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"

@implementation UIView (Frame)

- (CGFloat)udw_x {
    return self.frame.origin.x;
}

- (CGFloat)udw_y {
    return self.frame.origin.y;
}

- (CGFloat)udw_height {
    return self.frame.size.height;
}

- (CGFloat)udw_width {
    return self.frame.size.width;
}

- (CGFloat)udw_centerX {
    return self.center.x;
}

- (CGFloat)udw_centerY {
    return self.center.y;
}

- (CGSize)udw_size {
    return self.frame.size;
}

-(CGPoint)udw_origin {
    return self.frame.origin;
}

- (void)setUdw_x:(CGFloat)udw_x {
    CGRect rect = self.frame;
    rect.origin.x = udw_x;
    self.frame = rect;
}

- (void)setUdw_y:(CGFloat)udw_y {
    CGRect rect = self.frame;
    rect.origin.y = udw_y;
    self.frame = rect;
}

- (void)setUdw_width:(CGFloat)udw_width {
    CGRect rect = self.frame;
    rect.size.width = udw_width;
    self.frame = rect;
}

- (void)setUdw_height:(CGFloat)udw_height {
    CGRect rect = self.frame;
    rect.size.height = udw_height;
    self.frame = rect;
}

- (void)setUdw_centerX:(CGFloat)udw_centerX {
    CGPoint center = self.center;
    center.x = udw_centerX;
    self.center = center;
}

- (void)setUdw_centerY:(CGFloat)udw_centerY {
    CGPoint center = self.center;
    center.y = udw_centerY;
    self.center = center;
}

- (void)setUdw_size:(CGSize)udw_size {
    CGRect rect = self.frame;
    rect.size = udw_size;
    self.frame = rect;
}

-(void)setUdw_origin:(CGPoint)udw_origin {
    CGRect rect = self.frame;
    rect.origin = udw_origin;
    self.frame = rect;
}

@end