#import "udwOcUi.h"
@interface udwDrawerVC : UIViewController
@property (nonatomic, strong) UIViewController                 *mainViewController;
@property (nonatomic, strong) UIViewController                 *leftViewController;
@property (nonatomic, strong) UIView                           *maskView;
@property (nonatomic, strong) UIScreenEdgePanGestureRecognizer *openGesture;
@property (nonatomic, strong) UIScreenEdgePanGestureRecognizer *rightCloseGesture;
@property (nonatomic, strong) UIPanGestureRecognizer           *closeGesture;
@property (nonatomic, assign) CGFloat offset;
@property (nonatomic, copy) void (^beforeOpenLeftCallback)(void);
@property (nonatomic, copy) void (^afterOpenLeftCallback)(void);
@property (nonatomic, copy) void (^beforeCloseLeftCallback)(void);
@property (nonatomic, copy) void (^afterCloseLeftCallback)(void);
@property (nonatomic,copy) NSMutableArray<void (^)(void)> *openLeftCallbackArray;
@property (nonatomic,copy) NSMutableArray<void (^)(void)> *closeLeftCallbackArray;
@property (nonatomic, assign) CGSize preRotateSize;
@end
@implementation udwDrawerVC
- (instancetype)initWithMainViewController:(UIViewController *)mainViewController leftViewController:(UIViewController *)leftViewController{
    if (self = [super init]) {
        self.mainViewController = mainViewController;
        self.leftViewController = leftViewController;
    }
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];
    self.view.backgroundColor = [UIColor whiteColor];
    [self addChildViewController:self.leftViewController];
    self.leftViewController.view.frame = CGRectMake(-CGRectGetWidth(self.view.bounds)/2, 0, CGRectGetWidth(self.view.bounds), CGRectGetHeight(self.view.bounds));
    [self.leftViewController didMoveToParentViewController:self];
    [self addChildViewController:self.mainViewController];
    [self.view addSubview:self.mainViewController.view];
    [self.mainViewController didMoveToParentViewController:self];
    UIScreenEdgePanGestureRecognizer *openGesture = [[UIScreenEdgePanGestureRecognizer alloc] initWithTarget:self action:@selector(onPanGesture:)];
    openGesture.edges = UIRectEdgeLeft;
    [self.view addGestureRecognizer:openGesture];
    self.openGesture = openGesture;
    UIScreenEdgePanGestureRecognizer *rightCloseGesture = [[UIScreenEdgePanGestureRecognizer alloc] initWithTarget:self action:@selector(closePanGestre:)];
    rightCloseGesture.edges = UIRectEdgeRight;
    self.rightCloseGesture = rightCloseGesture;
    UIPanGestureRecognizer *closeGesture = [[UIPanGestureRecognizer alloc] initWithTarget:self action:@selector(closePanGestre:)];
    self.closeGesture = closeGesture;
    __weak typeof(self) weakSelf = self;
    if ([self.mainViewController isKindOfClass:[UINavigationController class]]) {
        UIViewController *navigationRootVC = [((UINavigationController *)self.mainViewController).viewControllers firstObject];
        if (navigationRootVC) {
            navigationRootVC.udwPresentVcCallback = ^(UIViewController *_vc) {
                [weakSelf closeLeftDrawerWithPresentVc:_vc];
            };
        }
    }
    [[NSNotificationCenter defaultCenter] addObserver:self
                                             selector:@selector(observerViewControllerChange:)
                                                 name:@"UINavigationControllerWillShowViewControllerNotification" object:nil];
    self.mainViewController.udwPresentVcCallback = ^(UIViewController *_vc) {
        [weakSelf closeLeftDrawerWithPresentVc:_vc];
    };
    self.leftViewController.udwPresentVcCallback = ^(UIViewController *_vc) {
        [weakSelf closeLeftDrawerWithPresentVc:_vc];
    };
}

- (void)closeLeftDrawerWithPresentVc:(UIViewController *)vc {
    if ([[vc class] isSubclassOfClass:[UIActivityViewController class]] || [[vc class] isSubclassOfClass:[UIAlertController class]]) {
        return;
    }
    [self closeLeftDrawerIfNeed];
}

- (void)presentViewController:(UIViewController *)viewControllerToPresent animated:(BOOL)flag completion:(void (^)(void))completion{
    [super presentViewController:viewControllerToPresent animated:flag completion:^{
        if (completion) {
            completion();
        }
        [self closeLeftDrawerWithPresentVc:viewControllerToPresent];
    }];
}

- (CGFloat)interpolateFrom:(CGFloat)from to:(CGFloat)to percent:(CGFloat)percent{
    return from + (to - from) * percent;
}

- (void)openLeftDrawer:(void (^)(void))completionCallback{
    [self handleDraweWithFlag:YES complete:completionCallback];
}

- (void)closeLeftDrawer:(void (^)(void))completionCallback{
    [self handleDraweWithFlag:NO complete:completionCallback];
}

- (void)handleDraweWithFlag:(BOOL)isOpen complete:(void (^)(void))completionCallback{
    CGFloat width = CGRectGetWidth(self.view.frame);
    CGFloat height = CGRectGetHeight(self.view.frame);
    CGRect mainToFrame;
    CGRect leftToFrame;
    CGFloat maskToAlpha;
    if (isOpen) {
        [self.view removeGestureRecognizer:self.openGesture];
        [self.maskView addGestureRecognizer:self.closeGesture];
        [self.view addSubview:self.maskView];
        [self.view insertSubview:self.leftViewController.view atIndex:0];
        leftToFrame = self.view.bounds;
        mainToFrame = CGRectMake(width-self.offset, 0, width, height);
        maskToAlpha = 0.6;
        if (self.offset <= 0) {
            [self.view addGestureRecognizer:self.rightCloseGesture];
        }
    }else{
        [self.view addGestureRecognizer:self.openGesture];
        [self.maskView removeGestureRecognizer:self.closeGesture];
        leftToFrame = CGRectMake(-width*0.5, 0, width, height);
        mainToFrame = self.view.bounds;
        maskToAlpha = 0;
        if (self.offset <= 0) {
            [self.view removeGestureRecognizer:self.rightCloseGesture];
        }
    }
    if (isOpen) {
        if (self.beforeOpenLeftCallback!=nil) {
            self.beforeOpenLeftCallback();
        }
    } else {
        if (self.beforeCloseLeftCallback!=nil) {
            self.beforeCloseLeftCallback();
        }
    }
    if (isOpen) {
        for (void (^callback)(void) in self.openLeftCallbackArray) {
            if (callback) {
                callback();
            }
        }
    } else {
        for (void (^callback)(void) in self.closeLeftCallbackArray) {
            if (callback) {
                callback();
            }
        }
    }
    [UIView animateWithDuration:0.2 animations:^{
        self.leftViewController.view.frame = leftToFrame;
        self.mainViewController.view.frame = mainToFrame;
        self.maskView.frame = mainToFrame;
        self.maskView.alpha = maskToAlpha;
        if (!isOpen) {
            [self.leftViewController.view removeFromSuperview];
            [self.maskView removeFromSuperview];
        }
    } completion:^(BOOL finished) {
        if (completionCallback != nil) {
            completionCallback();
        }
        if (isOpen) {
            if (self.afterOpenLeftCallback!=nil) {
                self.afterOpenLeftCallback();
            }
        } else {
            if (self.afterCloseLeftCallback!=nil) {
                self.afterCloseLeftCallback();
            }
        }
    }];
}

                       
- (void)tapOnMask:(UIPanGestureRecognizer *)sender{
    [self closeLeftDrawer:nil];
}
- (void)closePanGestre:(UIPanGestureRecognizer *)sender{
    if (sender.state == UIGestureRecognizerStateBegan) {
        
    }else if (sender.state == UIGestureRecognizerStateChanged){
        CGFloat maxTranslation = 300.f;
        CGPoint translation = [sender translationInView:sender.view];
        if (translation.x > 0) {
            translation.x = 0;
        }
        CGFloat percent = -translation.x / maxTranslation;
        if (percent > 1) {
            percent = 1;
        }
        CGFloat fromAlpha = 0.6;
        CGFloat toAlpha = 0;
        CGFloat alpha = [self interpolateFrom:fromAlpha to:toAlpha percent:percent];
        self.maskView.alpha = alpha;
        CGFloat width = CGRectGetWidth(self.mainViewController.view.bounds);
        CGFloat height = CGRectGetHeight(self.view.bounds);
        CGFloat mainFromX = width - self.offset;
        CGFloat mainToX = 0;
        CGFloat x = [self interpolateFrom:mainFromX to:mainToX percent:percent];
        self.mainViewController.view.frame = CGRectMake(x, 0, width, height);
        self.maskView.frame = self.mainViewController.view.frame;
        CGFloat leftFromX = 0;
        CGFloat leftTox = -width/2;
        CGFloat leftX = [self interpolateFrom:leftFromX to:leftTox percent:percent];
        self.leftViewController.view.frame = CGRectMake(leftX, 0, width, height);
    }else if (sender.state == UIGestureRecognizerStateCancelled || sender.state == UIGestureRecognizerStateEnded){
        CGPoint velocity = [sender velocityInView:sender.view];
        if (velocity.x > 0) {
            [self openLeftDrawer:nil];
        }else {
            [self closeLeftDrawer:nil];
        }
    }
}
- (void)onPanGesture:(UIPanGestureRecognizer *)sender{
    if (sender.state == UIGestureRecognizerStateBegan) {
        [self.view insertSubview:self.leftViewController.view atIndex:0];
        [self.view addSubview:self.maskView];
    }else if (sender.state == UIGestureRecognizerStateChanged){
        CGFloat maxTranslation = 300.f;
        CGPoint translation = [sender translationInView:sender.view];
        if (translation.x < 0) {
            translation.x = 0;
        }
        CGFloat percent = translation.x / maxTranslation;
        if (percent > 1) {
            percent = 1;
        }
        CGFloat fromAlpha = 0;
        CGFloat toAlpha = 0.6;
        CGFloat alpha = [self interpolateFrom:fromAlpha to:toAlpha percent:percent];
        self.maskView.alpha = alpha;
        CGFloat width = CGRectGetWidth(self.mainViewController.view.bounds);
        CGFloat height = CGRectGetHeight(self.view.bounds);
        CGFloat mainFromX = 0;
        CGFloat mainToX = width - self.offset;
        CGFloat x = [self interpolateFrom:mainFromX to:mainToX percent:percent];
        self.mainViewController.view.frame = CGRectMake(x, 0, width, height);
        self.maskView.frame = self.mainViewController.view.frame;
        CGFloat leftFromX = -width/2;
        CGFloat leftTox = 0;
        CGFloat leftX = [self interpolateFrom:leftFromX to:leftTox percent:percent];
        self.leftViewController.view.frame = CGRectMake(leftX, 0, width, height);
    }else if (sender.state == UIGestureRecognizerStateCancelled || sender.state == UIGestureRecognizerStateEnded){
        CGPoint velocity = [sender velocityInView:sender.view];
        if (velocity.x > 0) {
            [self openLeftDrawer:nil];
        }else {
            [self closeLeftDrawer:nil];
        }
    }
}

- (BOOL)isDrawerOpen{
    return self.maskView.alpha!=0;
}

- (void)closeLeftDrawerIfNeed{
    if ([self isDrawerOpen]) {
        [self closeLeftDrawer:nil];
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
- (NSMutableArray<void (^)(void)> *)openLeftCallbackArray{
    if (!_openLeftCallbackArray) {
        _openLeftCallbackArray = [NSMutableArray array];
    }
    return _openLeftCallbackArray;
}
- (NSMutableArray<void (^)(void)> *)closeLeftCallbackArray{
    if (!_closeLeftCallbackArray) {
        _closeLeftCallbackArray = [NSMutableArray array];
    }
    return _closeLeftCallbackArray;
}
                           
- (void)observerViewControllerChange:(NSNotification *)notification {
    NSDictionary *userInfo = [notification userInfo];
    UIViewController * fromVc = [userInfo objectForKey:@"UINavigationControllerLastVisibleViewController"];
    UIViewController * toVc = [userInfo objectForKey:@"UINavigationControllerNextVisibleViewController"];
                                                                   
    BOOL (^isActiveNav)(UINavigationController *) = ^BOOL (UINavigationController *nav) {
        if (toVc.navigationController == nav||fromVc.navigationController == nav) {
            return true;
        }
        return false;
    };
                                                
    if ([self.mainViewController isKindOfClass:[UINavigationController class]]) {
        if (isActiveNav((UINavigationController*)self.mainViewController)) {
            UIViewController * rootVc = ((UINavigationController*)self.mainViewController).viewControllers.firstObject;
            if ([fromVc isKindOfClass:[rootVc class]] && ![toVc isKindOfClass:[rootVc class]]) {
                                             
                [self closeLeftDrawerIfNeed];
            }
                                          
            self.openGesture.enabled = toVc==rootVc;
        }
    }
    if ([self.leftViewController isKindOfClass:[UINavigationController class]]) {
        if (isActiveNav((UINavigationController*)self.leftViewController)) {
            UIViewController * leftRootVc = ((UINavigationController*)self.leftViewController).viewControllers.firstObject;
            self.rightCloseGesture.enabled = toVc==leftRootVc;
        }
    }
}
- (void)rotateReloadUi:(BOOL)isLandscape {
    CGFloat width = CGRectGetWidth(self.view.frame);
    CGFloat height = CGRectGetHeight(self.view.frame);
    self.leftViewController.view.frame = CGRectMake(-CGRectGetWidth(self.view.bounds)/2, 0, CGRectGetWidth(self.view.bounds), CGRectGetHeight(self.view.bounds));
    self.maskView.frame = self.mainViewController.view.frame;
    if ([self isDrawerOpen]) {
        CGRect leftToFrame = self.view.bounds;
        CGRect  mainToFrame = CGRectMake(width-self.offset, 0, width, height);
        self.leftViewController.view.frame = leftToFrame;
        self.mainViewController.view.frame = mainToFrame;
    }
}

- (void)dealloc{
    [[NSNotificationCenter defaultCenter] removeObserver:self];
    NSLog(@"----------dealloc %@----------",NSStringFromClass([self class]));
}
@end

@interface udwDrawerVCManager()
@property (nonatomic, strong)udwDrawerVC *drawerVC;
@property (nonatomic, strong)UIViewController *mainViewController;
@property (nonatomic, strong)UIViewController *leftViewController;
@end
@implementation udwDrawerVCManager

+ (instancetype)sharedManager{
    static udwDrawerVCManager *manageer = nil;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        manageer = [[udwDrawerVCManager alloc] init];
    });
    return manageer;
}

- (UIViewController *)InitWithMainViewController:(UIViewController *)mainViewController leftViewController:(UIViewController *)leftViewController{
                                 
                                                                                                                     
    _drawerVC = [[udwDrawerVC alloc] initWithMainViewController:mainViewController leftViewController:leftViewController];
    self.mainViewController = mainViewController;
    self.leftViewController = leftViewController;
    if (_drawerVC.offset == 0) {
        CGFloat defaultOffset = 50;
        _drawerVC.offset = defaultOffset;
        self.offset = defaultOffset;
    }
           
    return self.drawerVC;
}

- (void)rotateReloadUi:(BOOL)isLandscape {
    [self.drawerVC rotateReloadUi:isLandscape];
}

- (UIViewController *)drawerViewController{
    return self.drawerVC;
}

- (void)openLeftDrawer{
    [self.drawerVC openLeftDrawer:nil];
}

- (void)closeLeftDrawer{
    [self.drawerVC closeLeftDrawer:nil];
}

- (void)setOffset:(CGFloat)offset{
    self.drawerVC.offset = offset;
    _offset = offset;
}

                                                                                                          
- (void)setOpenLeftCallBackWithBeforeAnimation:(BOOL)BeforeAnimation callback:(void (^)(void))CallBack{
    if (BeforeAnimation) {
        self.drawerVC.beforeOpenLeftCallback = nil;
        self.drawerVC.beforeOpenLeftCallback = CallBack;
    } else {
        self.drawerVC.afterOpenLeftCallback = nil;
        self.drawerVC.afterOpenLeftCallback = CallBack;
    }
}

                                                                                                          
- (void)setCloseLeftCallBackWithBeforeAnimation:(BOOL)BeforeAnimation callback:(void (^)(void))CallBack{
    if (BeforeAnimation) {
        self.drawerVC.beforeCloseLeftCallback = nil;
        self.drawerVC.beforeCloseLeftCallback = CallBack;
    } else {
        self.drawerVC.afterCloseLeftCallback = nil;
        self.drawerVC.afterCloseLeftCallback = CallBack;
    }
}

- (void)setOpenLeftCallBack:(void (^)(void))CallBack{
    if (CallBack == nil) {
        return;
    }
    [self.drawerVC.openLeftCallbackArray addObject:CallBack];
}

- (void)setCloseLeftCallBack:(void (^)(void))CallBack{
    if (CallBack == nil) {
        return;
    }
    [self.drawerVC.closeLeftCallbackArray addObject:CallBack];
}

@end




