#import "udwOcUi.h"
#import <objc/runtime.h>

@interface UIView (udwPrivateShadow)
@property (nonatomic, weak) UIView *targetShadowView;
@end
static void * udwPrivateShadowTargetShadowViewKey = &udwPrivateShadowTargetShadowViewKey;
@implementation UIView (udwPrivateShadow)
- (UIView *)targetShadowView{
    return objc_getAssociatedObject(self,udwPrivateShadowTargetShadowViewKey);
}
- (void)setTargetShadowView:(UIView *)targetShadowView{
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        Method udw_removeFromSuperview = class_getInstanceMethod([self class], @selector(udwPrivateShadow_removeFromSuperview));
        Method removeFromSuperview = class_getInstanceMethod([self class], @selector(removeFromSuperview));
        method_exchangeImplementations(udw_removeFromSuperview,removeFromSuperview);
    });
    objc_setAssociatedObject(self, udwPrivateShadowTargetShadowViewKey, targetShadowView, OBJC_ASSOCIATION_ASSIGN);
}
- (void)udwPrivateShadow_removeFromSuperview{
    if (self.targetShadowView) {
        [self.targetShadowView udwPrivateShadow_removeFromSuperview];
        self.targetShadowView = nil;
    }
    [self udwPrivateShadow_removeFromSuperview];
}
@end
@interface udwPrivateShadowView : UIView
@end
@implementation udwPrivateShadowView
- (void)layoutSubviews {
    [super layoutSubviews];
    CAAnimation *animation = [self.layer animationForKey:@"bounds.size"];
    if (animation) {
        CABasicAnimation *shadowPathAnimation = [CABasicAnimation animationWithKeyPath:@"udwPrivateShadowPath"];
        shadowPathAnimation.timingFunction = animation.timingFunction;
        shadowPathAnimation.duration = animation.duration;
        shadowPathAnimation.toValue = [UIBezierPath bezierPathWithRoundedRect:self.layer.bounds cornerRadius:self.layer.cornerRadius];
        [self.layer addAnimation:shadowPathAnimation forKey:@"udwPrivateShadowPath"];
    }
    self.layer.shadowPath = [UIBezierPath bezierPathWithRoundedRect:self.layer.bounds cornerRadius:self.layer.cornerRadius].CGPath;
}
@end
void udwViewAddShadow(UIView *addShadowView,UIColor *shadowColor,CGSize shadowOffset,float shadowOpacity,CGFloat shadowRadius){
    void (^addshadow)(UIView *,UIColor*,CGSize,float,CGFloat) = ^void(UIView* shadowView,UIColor* _shadowColor,CGSize _shadowOffset,float _shadowOpacity,CGFloat _shadowRadius){
        UIBezierPath *shadowPath = [UIBezierPath bezierPathWithRoundedRect:addShadowView.bounds cornerRadius:addShadowView.layer.cornerRadius];
        shadowView.layer.masksToBounds = NO;
        shadowView.layer.cornerRadius = addShadowView.layer.cornerRadius;
        shadowView.layer.shadowColor = _shadowColor.CGColor;
        shadowView.layer.shadowOffset = _shadowOffset;
        shadowView.layer.shadowOpacity = _shadowOpacity;
        shadowView.layer.shadowRadius = _shadowRadius;
        shadowView.layer.shadowPath = shadowPath.CGPath;
    };
        if (addShadowView.targetShadowView!=nil||addShadowView.superview==nil) {
            return;
        }
        udwPrivateShadowView *shadowView = [[udwPrivateShadowView alloc] initWithFrame:addShadowView.frame];
        if (addShadowView.backgroundColor == nil || addShadowView.backgroundColor == [UIColor clearColor]) {                                  
            shadowView.backgroundColor = [UIColor clearColor];
            CAShapeLayer *layer = [CAShapeLayer layer];
            if (addShadowView.layer.cornerRadius !=0) {
                UIBezierPath *path = [UIBezierPath bezierPathWithRoundedRect:shadowView.bounds cornerRadius:addShadowView.layer.cornerRadius];
                layer.path = path.CGPath;
            } else {
                layer.path = [UIBezierPath bezierPathWithRect:shadowView.bounds].CGPath;
            }
            layer.fillColor = [UIColor whiteColor].CGColor;
            [shadowView.layer addSublayer:layer];
        }
        addshadow(shadowView,shadowColor,shadowOffset,shadowOpacity,shadowRadius);
        [addShadowView.superview insertSubview:shadowView belowSubview:addShadowView];                             
        [shadowView setTranslatesAutoresizingMaskIntoConstraints:false];
        [addShadowView.superview addConstraints:@[[NSLayoutConstraint constraintWithItem:shadowView attribute:NSLayoutAttributeLeading relatedBy:NSLayoutRelationEqual toItem:addShadowView attribute:NSLayoutAttributeLeading multiplier:1.0 constant:0],[NSLayoutConstraint constraintWithItem:shadowView attribute:NSLayoutAttributeTrailing relatedBy:NSLayoutRelationEqual toItem:addShadowView attribute:NSLayoutAttributeTrailing multiplier:1.0 constant:0],[NSLayoutConstraint constraintWithItem:shadowView attribute:NSLayoutAttributeTop relatedBy:NSLayoutRelationEqual toItem:addShadowView attribute:NSLayoutAttributeTop multiplier:1.0 constant:0],[NSLayoutConstraint constraintWithItem:shadowView attribute:NSLayoutAttributeBottom relatedBy:NSLayoutRelationEqual toItem:addShadowView attribute:NSLayoutAttributeBottom multiplier:1.0 constant:0]]];
        addShadowView.targetShadowView = shadowView;
}

void udwViewDeleteShadow(UIView *shadowView){
    if (shadowView.targetShadowView) {
        [shadowView.targetShadowView removeFromSuperview];
        shadowView.targetShadowView = nil;
    }
}

