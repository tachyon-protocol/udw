#import "udwOcMac.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"
#import <QuartzCore/QuartzCore.h>

@interface udwPopView ()<CAAnimationDelegate>
{
    CGRect superViewFrame;
    CGRect subViewFrame;
    CGRect startFrame;
    CGRect endFrame;
    
    int removeLock;
}
@property (strong, nonatomic) NSView *animatedView;

@end

@implementation udwPopView

- (instancetype)init{
    if (self = [super init]) {
        [self loadPopView];
    }
    return self;
}

                                          
- (void)mouseDown:(NSEvent *)event{
    if (!self.autoRemoveEnable) {
        return;
    }
    CGPoint point = [self convertPoint:event.locationInWindow toView:self.animatedView];
    if (point.x < 0 ||
        point.x > subViewFrame.size.width ||
        point.y < 0 ||
        point.y > subViewFrame.size.height) {
        [self removeSelf];
    }
}

- (instancetype)initWithFrame:(CGRect)frame{
    if (self = [super initWithFrame:frame]) {
        superViewFrame = frame;
        self.frame = frame;
        [self loadPopView];
    }
    return self;
}

- (instancetype)initWithSubView:(__kindof NSView *)subView{
    if (self = [super init]) {
        [self addSubview:subView];
        [self loadPopView];
    }
    return self;
}

- (void)loadPopView{
    self.ty_backgroundColor = [NSColor colorWithRed:0 green:0 blue:0 alpha:0.7];
    removeLock = 0;
    self.layer.masksToBounds = true;
    self.duration = 0.3;
    self.autoRemoveEnable = true;
    self.animatedEnable = true;
    self.animationStyle = udwPopViewAnimationStyleFade;
}

                                          

- (void)addSubview:(NSView *)view{
    self.animatedView = view;
    self.animatedView.wantsLayer = true;
    self.wantsLayer = true;
    [super addSubview:view];
}

                       

- (void)willShowCallBack{
    if (self.willShow) {
        self.willShow(self);
    }
}
- (void)willDismissCallBack{
    if (self.willDismiss) {
        self.willDismiss(self);
    }
}
- (void)didDismissCallBack{
    if (self.didDismiss) {
        self.didDismiss(self);
    }
}
- (void)didShowCallBack{
    if (self.didShow) {
        self.didShow(self);
    }
}

                      

- (void)showPopViewOn:(NSView *)view{
    removeLock = 1;
    [self willShowCallBack];
    if (!self.superview) {
        [view addSubview:self];
    }
    if (CGRectEqualToRect(CGRectZero, superViewFrame)) {
        superViewFrame = view.bounds;
    }
    self.frame = superViewFrame;
    [self settingSubViewsWithAnimated:self.animatedEnable Duration:self.duration];
}

                     

- (void)removeSelf{
    removeLock = 0;
    [self willDismissCallBack];
    [[NSNotificationCenter defaultCenter] removeObserver:self];
    [self removeSelfWithAnimated:self.animatedEnable];
}

- (void)removeSelfWithAnimated:(BOOL)animated{
    if (animated) {
        if (self.animationStyle == udwPopViewAnimationStyleScale) {
            [self executeTransformAnimationIsShowing:NO];
            return;
        }else if (self.animationStyle == udwPopViewAnimationStyleFade){
            [self executeFadeAnimationIsShowing:NO];
            return;
        }
        [self executeAnimationIsShowing:NO];
    }else{
        [self didRemove];
    }
}

- (void)didRemove{
    if (removeLock) {
        return;
    }
    [self didDismissCallBack];
    [self removeFromSuperview];
}

                                

- (void)settingSubViewsWithAnimated:(BOOL)animated Duration:(NSTimeInterval)duration{
    if (animated) {
        [self layoutSubtreeIfNeeded];
        if (CGRectEqualToRect(CGRectZero, subViewFrame)) {
            subViewFrame = self.animatedView.frame;
        }
        if (self.animationStyle == udwPopViewAnimationStyleBottomToTop) {
            [self prepareAnimationFromBottomToTop];
        }else if (self.animationStyle == udwPopViewAnimationStyleTopToBottom){
            [self prepareAnimationFromTopToBottom];
        }else if (self.animationStyle == udwPopViewAnimationStyleRightToLeft){
            [self prepareAnimationFromRightToLeft];
        }else if (self.animationStyle == udwPopViewAnimationStyleLeftToRight){
            [self prepareAnimationFromLeftToRight];
        }else if (self.animationStyle == udwPopViewAnimationStyleFade){
            [self executeFadeAnimationIsShowing:YES];return;
        }else if (self.animationStyle == udwPopViewAnimationStyleScale){
            [self executeTransformAnimationIsShowing:YES];return;
        }
        [self executeAnimationIsShowing:YES];
    }else{
        [self didShowCallBack];
    }
}

- (void)prepareAnimationFromTopToBottom{
    startFrame = CGRectMake(subViewFrame.origin.x, superViewFrame.size.height, subViewFrame.size.width, subViewFrame.size.height);
    endFrame = subViewFrame;
}
- (void)prepareAnimationFromLeftToRight{
    startFrame = CGRectMake(- subViewFrame.size.width, subViewFrame.origin.y, subViewFrame.size.width, subViewFrame.size.height);
    endFrame = subViewFrame;
}
- (void)prepareAnimationFromRightToLeft{
    startFrame = CGRectMake(superViewFrame.size.width, subViewFrame.origin.y, subViewFrame.size.width, subViewFrame.size.height);
    endFrame = subViewFrame;
}
- (void)prepareAnimationFromBottomToTop{
    startFrame = CGRectMake(subViewFrame.origin.x, -subViewFrame.size.height, subViewFrame.size.width, subViewFrame.size.height);
    endFrame = subViewFrame;
}

- (void)executeAnimationIsShowing:(BOOL)isShowing{
    if (isShowing) {
        [self.animatedView setFrame:self->startFrame];
        [NSAnimationContext runAnimationGroup:^(NSAnimationContext * _Nonnull context) {
            context.allowsImplicitAnimation = true;
            [context setDuration:self.duration];
            [self.animatedView.animator setFrame:self->endFrame];
        } completionHandler:^{
            [self didShowCallBack];
        }];
    }else{
        [self.animatedView setFrame:endFrame];
        [NSAnimationContext runAnimationGroup:^(NSAnimationContext * _Nonnull context) {
            context.allowsImplicitAnimation = true;
            [context setDuration:self.duration];
            [self.animatedView.animator setFrame:self->startFrame];
        } completionHandler:^{
            [self didRemove];
        }];
    }
}

- (void)executeFadeAnimationIsShowing:(BOOL)isShowing{
    if (isShowing) {
        self.animatedView.alphaValue = 0;
        [NSAnimationContext runAnimationGroup:^(NSAnimationContext * _Nonnull context) {
            [context setDuration:self.duration];
            [self.animatedView.animator setAlphaValue:1];
        } completionHandler:^{
            [self didShowCallBack];
        }];
    }else{
        [self.animatedView setAlphaValue:1];
        [NSAnimationContext runAnimationGroup:^(NSAnimationContext * _Nonnull context) {
            [context setDuration:self.duration];
            [self.animatedView.animator setAlphaValue:0];
        } completionHandler:^{
            [self didRemove];
        }];
    }
}


- (void)executeTransformAnimationIsShowing:(BOOL)isShowing{
    if (isShowing) {
                                                                
        CABasicAnimation *animation = [CABasicAnimation animationWithKeyPath:@"transform"];
        CATransform3D tr = CATransform3DIdentity;
        tr = CATransform3DTranslate(tr, self.animatedView.bounds.size.width / 2, self.animatedView.bounds.size.height / 2, 0);
        tr = CATransform3DScale(tr, 0.001, 0.001, 1);
        tr = CATransform3DTranslate(tr, -self.animatedView.bounds.size.width / 2, -self.animatedView.bounds.size.height / 2, 0);
        animation.delegate = [udwWeakProxy proxyWithTarget:self];
        animation.duration = self.duration;
        animation.autoreverses = false;
        animation.repeatCount = 0;
        animation.removedOnCompletion = false;
        animation.fillMode = kCAFillModeForwards;
        animation.fromValue = [NSValue valueWithCATransform3D:tr];
        animation.toValue = [NSValue valueWithCATransform3D:CATransform3DIdentity];
        [self.animatedView.layer addAnimation:animation forKey:@"popview_show"];
    }else{
        CABasicAnimation *animation = [CABasicAnimation animationWithKeyPath:@"transform"];
        CATransform3D tr = CATransform3DIdentity;
        tr = CATransform3DTranslate(tr, self.animatedView.bounds.size.width / 2, self.animatedView.bounds.size.height / 2, 0);
        tr = CATransform3DScale(tr, 0.001, 0.001, 1);
        tr = CATransform3DTranslate(tr, -self.animatedView.bounds.size.width / 2, -self.animatedView.bounds.size.height / 2, 0);
        animation.delegate = [udwWeakProxy proxyWithTarget:self];
        animation.duration = self.duration;
        animation.autoreverses = false;
        animation.repeatCount = 0;
        animation.removedOnCompletion = false;
        animation.fillMode = kCAFillModeForwards;
        animation.fromValue = [NSValue valueWithCATransform3D:CATransform3DIdentity];
        animation.toValue = [NSValue valueWithCATransform3D:tr];
        [self.animatedView.layer addAnimation:animation forKey:@"popview_remove"];
    }
}

- (void)animationDidStop:(CAAnimation *)anim finished:(BOOL)flag{
    if ([anim isEqual:[self.animatedView.layer animationForKey:@"popview_show"]]) {
        [self didShowCallBack];
    }else if([anim isEqual:[self.animatedView.layer animationForKey:@"popview_remove"]]){
        [self didRemove];
    }
}

                              

- (void)setAnimationStyle:(udwPopViewAnimationStyle)animationStyle{
    _animationStyle = animationStyle;
}

@end
