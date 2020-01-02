
#import <objc/runtime.h>
#import <Cocoa/Cocoa.h>

@implementation NSView (Add)

- (void)setTy_backgroundColor:(NSColor *)backgroundColor{
    self.wantsLayer = true;
    self.layer.backgroundColor = backgroundColor.CGColor;
    objc_setAssociatedObject(self, @selector(ty_backgroundColor), backgroundColor, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
}

- (NSColor *)ty_backgroundColor{
    return (NSColor *)objc_getAssociatedObject(self, @selector(ty_backgroundColor));
}

-(void)setCornerRadius:(CGFloat)cornerRadius{
    self.wantsLayer = true;
    self.layer.cornerRadius = cornerRadius;
      objc_setAssociatedObject(self, @selector(cornerRadius), @(cornerRadius), OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    
}

-(CGFloat)cornerRadius{
     return [(NSNumber *)objc_getAssociatedObject(self, @selector(cornerRadius)) floatValue];
}

- (void)setAnchorPoint:(CGPoint)anchor{
    self.wantsLayer = true;
    CGPoint newPoint = CGPointMake(self.bounds.size.width * anchor.x, self.bounds.size.height * anchor.y);
    CGPoint oldPoint = CGPointMake(self.bounds.size.width * self.layer.anchorPoint.x, self.bounds.size.height * self.layer.anchorPoint.y);
    newPoint = CGPointApplyAffineTransform(newPoint, self.layer.affineTransform);
    oldPoint = CGPointApplyAffineTransform(oldPoint, self.layer.affineTransform);
    
    CGPoint position = self.layer.position;
    position.x -= oldPoint.x;
    position.x += newPoint.x;
    
    position.y -= oldPoint.y;
    position.y += newPoint.y;
    
    self.layer.position = position;
    self.layer.anchorPoint = anchor;
}

- (void)setOrigin:(CGPoint)origin{
    self.frame = CGRectMake(origin.x, origin.y, self.width, self.height);
}
- (void)setSize:(CGSize)size{
    self.frame = CGRectMake(self.originX, self.originY, size.width, size.height);
}

- (void)setOriginX:(CGFloat)originX{
    self.frame = CGRectMake(originX, self.originY, self.width, self.height);
}

- (void)setOriginY:(CGFloat)originY{
    self.frame = CGRectMake(self.originX, originY, self.width, self.height);
}

- (void)setWidth:(CGFloat)width{
    self.frame = CGRectMake(self.originX, self.originY, width, self.height);
}

- (void)setHeight:(CGFloat)height{
    self.frame = CGRectMake(self.originX, self.originY, self.width, height);
}

- (CGSize)size{
    return self.frame.size;
}

- (CGPoint)origin{
    return self.frame.origin;
}

- (CGFloat)width{
    return self.size.width;
}

- (CGFloat)height{
    return self.size.height;
}

- (CGFloat)originX{
    return self.origin.x;
}

- (CGFloat)originY{
    return self.origin.y;
}

@end
