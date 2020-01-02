#import "udwOcUi.h"
@interface udwCoverDrawerVC()
@property (nonatomic, strong) UIViewController                 *mainViewController;
@property (nonatomic, strong) UIViewController                 *rightViewController;
@property (nonatomic, strong) UIView                           *maskView;
@property (nonatomic, strong) UIScreenEdgePanGestureRecognizer *openGesture;
@property (nonatomic, strong) UIPanGestureRecognizer           *closeGesture;
@property (nonatomic, assign) BOOL enableOpenGesture;
@property (nonatomic, assign) CGFloat offset;
@property (nonatomic,copy) NSMutableArray<void (^)(void)> *openRightCallbackArray;
@property (nonatomic,copy) NSMutableArray<void (^)(void)> *closeRightCallbackArray;
@end
@implementation udwCoverDrawerVC

                     
- (instancetype)initWithMainViewController:(UIViewController *)mainViewController rightViewController:(UIViewController *)rightViewController{
    if (self = [super init]) {
        self.mainViewController = mainViewController;
        self.rightViewController = rightViewController;
        self.enableOpenGesture = true;
    }
    return self;
}

- (void)openRightDrawer:(void (^)(void))completionCallback{
    [self handleDraweWithFlag:YES complete:completionCallback];
}

- (void)closeRightDrawer:(void (^)(void))completionCallback{
    [self handleDraweWithFlag:NO complete:completionCallback];
}

- (CGFloat)getDrawerWidth{
    return self.view.frame.size.width-self.offset;
}

- (void)setEnableOpenGesture:(BOOL)enableOpenGesture{
    _enableOpenGesture = enableOpenGesture;
}

- (BOOL)isDrawerOpen{
    return self.maskView.alpha!=0;
}

- (void)setOpenRightCallBack:(void (^)(void))CallBack{
    if (CallBack == nil) {
        return;
    }
    [self.openRightCallbackArray addObject:CallBack];
}

- (void)setCloseRightCallBack:(void (^)(void))CallBack{
    if (CallBack == nil) {
        return;
    }
    [self.closeRightCallbackArray addObject:CallBack];
}

                      
- (void)viewDidLoad {
    [super viewDidLoad];
    self.view.backgroundColor = [UIColor whiteColor];
    [self addChildViewController:self.mainViewController];
    [self.view addSubview:self.mainViewController.view];
    self.mainViewController.view.frame = self.view.bounds;
    [self.mainViewController didMoveToParentViewController:self];
    [self addChildViewController:self.rightViewController];
    self.rightViewController.view.frame = CGRectMake(CGRectGetWidth(self.view.bounds), 0, CGRectGetWidth(self.view.bounds), CGRectGetHeight(self.view.bounds));
    [self.rightViewController didMoveToParentViewController:self];
    UIScreenEdgePanGestureRecognizer *openGesture = [[UIScreenEdgePanGestureRecognizer alloc] initWithTarget:self action:@selector(onPanGesture:)];
    openGesture.edges = UIRectEdgeRight;
    [self.view addGestureRecognizer:openGesture];
    self.openGesture = openGesture;
    UIPanGestureRecognizer *closeGesture = [[UIPanGestureRecognizer alloc] initWithTarget:self action:@selector(closePanGestre:)];
    self.closeGesture = closeGesture;
    __weak typeof(self) weakSelf = self;
    if ([self.mainViewController isKindOfClass:[UINavigationController class]]) {
        [[NSNotificationCenter defaultCenter] addObserver:self
                                                 selector:@selector(observerViewControllerChange:)
                                                     name:@"UINavigationControllerWillShowViewControllerNotification" object:nil];
        UIViewController *navigationRootVC = [((UINavigationController *)self.mainViewController).viewControllers firstObject];
        if (navigationRootVC) {
            navigationRootVC.udwPresentVcCallback = ^(UIViewController *_vc) {
                [weakSelf closeRightDrawerWithPresentVc:_vc];
            };
        }
    }
    self.mainViewController.udwPresentVcCallback = ^(UIViewController *_vc) {
        [weakSelf closeRightDrawerWithPresentVc:_vc];
    };
    self.rightViewController.udwPresentVcCallback = ^(UIViewController *_vc) {
        [weakSelf closeRightDrawerWithPresentVc:_vc];
    };
}

- (void)closeRightDrawerWithPresentVc:(UIViewController *)vc {
    if ([[vc class] isSubclassOfClass:[UIActivityViewController class]]) {
        return;
    }
    [self closeLeftDrawerIfNeed];
}

- (void)presentViewController:(UIViewController *)viewControllerToPresent animated:(BOOL)flag completion:(void (^)(void))completion{
    [super presentViewController:viewControllerToPresent animated:flag completion:^{
        if (completion) {
            completion();
        }
        [self closeRightDrawerWithPresentVc:viewControllerToPresent];
    }];
}

- (CGFloat)interpolateFrom:(CGFloat)from to:(CGFloat)to percent:(CGFloat)percent{
    return from + (to - from) * percent;
}

- (void)handleDraweWithFlag:(BOOL)isOpen complete:(void (^)(void))completionCallback{
    CGFloat width = CGRectGetWidth(self.view.frame);
    CGFloat height = CGRectGetHeight(self.view.frame);
    CGRect rightToFrame;
    CGFloat maskToAlpha;
    if (isOpen) {
        [self.view removeGestureRecognizer:self.openGesture];
        [self.maskView addGestureRecognizer:self.closeGesture];
        [self.view addSubview:self.maskView];
        [self.view addSubview:self.rightViewController.view];
        [self.view bringSubviewToFront:self.rightViewController.view];
        rightToFrame = CGRectMake(self.offset, 0, width, height);
        maskToAlpha = 0.6;
    }else{
        [self.view addGestureRecognizer:self.openGesture];
        [self.maskView removeGestureRecognizer:self.closeGesture];
        rightToFrame = CGRectMake(width, 0, width, height);
        maskToAlpha = 0;
    }
    if (isOpen) {
        for (void (^callback)(void) in self.openRightCallbackArray) {
            if (callback) {
                callback();
            }
        }
    } else {
        for (void (^callback)(void) in self.closeRightCallbackArray) {
            if (callback) {
                callback();
            }
        }
    }
    [UIView animateWithDuration:0.2 animations:^{
        self.rightViewController.view.frame = rightToFrame;
        self.maskView.alpha = maskToAlpha;
    } completion:^(BOOL finished) {
        if (completionCallback != nil) {
            completionCallback();
        }
    }];
}
                       
- (void)tapOnMask:(UIPanGestureRecognizer *)sender{
    [self closeRightDrawer:nil];
}
- (void)closePanGestre:(UIPanGestureRecognizer *)sender{
    if (sender.state == UIGestureRecognizerStateBegan) {

    }else if (sender.state == UIGestureRecognizerStateChanged){
        CGFloat width = CGRectGetWidth(self.view.bounds);
        CGFloat maxTranslation = width-20.f;
        CGPoint translation = [sender translationInView:sender.view];
        if (translation.x < 0) {
            translation.x = 0;
        }
        CGFloat percent = translation.x / maxTranslation;
        if (percent > 1) {
            percent = 1;
        }
        CGFloat fromAlpha = 0.6;
        CGFloat toAlpha = 0;
        CGFloat alpha = [self interpolateFrom:fromAlpha to:toAlpha percent:percent];
        self.maskView.alpha = alpha;
        CGFloat height = CGRectGetHeight(self.view.bounds);
        CGFloat rightFromX = self.offset;
        CGFloat rightTox = width;
        CGFloat rightX = [self interpolateFrom:rightFromX to:rightTox percent:percent];
        self.rightViewController.view.frame = CGRectMake(rightX, 0, width, height);
    }else if (sender.state == UIGestureRecognizerStateCancelled || sender.state == UIGestureRecognizerStateEnded){
        CGPoint velocity = [sender velocityInView:sender.view];
        if (velocity.x < 0) {
            [self openRightDrawer:nil];
        }else {
            [self closeRightDrawer:^{
                [self.rightViewController.view removeFromSuperview];
                [self.maskView removeFromSuperview];
            }];
        }
    }
}
- (void)onPanGesture:(UIPanGestureRecognizer *)sender{
    if (!self.enableOpenGesture) {
        return;
    }
    if (sender.state == UIGestureRecognizerStateBegan) {
        [self.view addSubview:self.maskView];
        [self.view addSubview:self.rightViewController.view];
        [self.view bringSubviewToFront:self.rightViewController.view];
    }else if (sender.state == UIGestureRecognizerStateChanged){
        CGFloat width = CGRectGetWidth(self.view.bounds);
        CGFloat maxTranslation = width-20.f;
        CGPoint translation = [sender translationInView:sender.view];
        if (translation.x > 0) {
            translation.x = 0;
        }
        CGFloat percent = -translation.x / maxTranslation;
        if (percent > 1) {
            percent = 1;
        }
        CGFloat fromAlpha = 0;
        CGFloat toAlpha = 0.6;
        CGFloat alpha = [self interpolateFrom:fromAlpha to:toAlpha percent:percent];
        self.maskView.alpha = alpha;
        CGFloat height = CGRectGetHeight(self.view.bounds);
        CGFloat rightFromX = width;
        CGFloat rightTox = self.offset;
        CGFloat rightX = [self interpolateFrom:rightFromX to:rightTox percent:percent];
        self.rightViewController.view.frame = CGRectMake(rightX, 0, width, height);
    }else if (sender.state == UIGestureRecognizerStateCancelled || sender.state == UIGestureRecognizerStateEnded){
        CGPoint velocity = [sender velocityInView:sender.view];
        if (velocity.x < 0) {
            [self openRightDrawer:nil];
        }else {
            [self closeRightDrawer:nil];
        }
    }
}

- (void)closeLeftDrawerIfNeed{
    if ([self isDrawerOpen]) {
        [self closeRightDrawer:nil];
    }
}

                     
- (UIView *)maskView{
    if (!_maskView) {
        _maskView = [[UIView alloc] initWithFrame:self.mainViewController.view.bounds];
        _maskView.layer.shadowColor = [UIColor blackColor].CGColor;
        [_maskView addGestureRecognizer:[[UITapGestureRecognizer alloc] initWithTarget:self action:@selector(tapOnMask:)]];
        _maskView.layer.shadowOpacity = 0.1;
                                       
        _maskView.layer.shadowPath = [UIBezierPath bezierPathWithRect:_maskView.bounds].CGPath;
        _maskView.backgroundColor = [UIColor blackColor];
        _maskView.alpha = 0;
    }
    return _maskView;
}
- (NSMutableArray<void (^)(void)> *)openRightCallbackArray{
    if (!_openRightCallbackArray) {
        _openRightCallbackArray = [NSMutableArray array];
    }
    return _openRightCallbackArray;
}
- (NSMutableArray<void (^)(void)> *)closeRightCallbackArray{
    if (!_closeRightCallbackArray) {
        _closeRightCallbackArray = [NSMutableArray array];
    }
    return _closeRightCallbackArray;
}
                           
- (void)observerViewControllerChange:(NSNotification *)notification {
    NSDictionary *userInfo = [notification userInfo];
    UIViewController * fromVc = [userInfo objectForKey:@"UINavigationControllerLastVisibleViewController"];
    UIViewController * toVc = [userInfo objectForKey:@"UINavigationControllerNextVisibleViewController"];
    UIViewController * rootVc = ((UINavigationController*)self.mainViewController).viewControllers.firstObject;
                                                                                                        
                                                
    if (toVc.navigationController == ((UINavigationController*)self.mainViewController)) {
                                                                                        
        if ([fromVc isKindOfClass:[rootVc class]] && ![toVc isKindOfClass:[rootVc class]]) {
                                         
            [self closeLeftDrawerIfNeed];
        }
                                      
        self.openGesture.enabled = [toVc isKindOfClass:[rootVc class]];
    }
}
- (void)dealloc{
    [[NSNotificationCenter defaultCenter] removeObserver:self];
    NSLog(@"----------dealloc %@----------",NSStringFromClass([self class]));
}
@end

@interface udwCoverDrawerVCListManager()
@property (nonatomic,strong)NSMutableDictionary *dict;
@end
@implementation udwCoverDrawerVCListManager
+ (instancetype)sharedManager{
    static udwCoverDrawerVCListManager *manageer = nil;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        manageer = [[udwCoverDrawerVCListManager alloc] init];
    });
    return manageer;
}
+ (UIViewController *)InitWithMainViewController:(UIViewController *)mainViewController rightViewController:(UIViewController *)rightViewController tag:(NSString *)tag offset:(CGFloat)offset{
    return [[self sharedManager] InitWithMainViewController:mainViewController rightViewController:rightViewController tag:tag offset:offset];
}
+ (udwCoverDrawerVC *)GetOneDrawer:(NSString *)tag{
    return [[self sharedManager] GetOneDrawer:tag];
}
- (UIViewController *)InitWithMainViewController:(UIViewController *)mainViewController rightViewController:(UIViewController *)rightViewController tag:(NSString *)tag offset:(CGFloat)offset{
    udwCoverDrawerVC *drawerVc =  [[udwCoverDrawerVC alloc] initWithMainViewController:mainViewController rightViewController:rightViewController];
    if (offset == 0) {
        offset = 20;
    }
    drawerVc.offset = offset;
    [self.dict setObject:drawerVc forKey:tag];
    return drawerVc;
}
- (udwCoverDrawerVC *)GetOneDrawer:(NSString *)tag{
    udwCoverDrawerVC *drawerVc = [self.dict objectForKey:tag];
    return drawerVc;
}
- (NSMutableDictionary *)dict{
    if (!_dict) {
        _dict = [NSMutableDictionary dictionary];
    }
    return _dict;
}
@end
