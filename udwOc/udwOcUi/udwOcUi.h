#import <UIKit/UIKit.h>
#import <MobileCoreServices/MobileCoreServices.h>
#import <UserNotifications/UserNotifications.h>
#import <AVKit/AVKit.h>

#pragma GCC diagnostic ignored "-Wnullability-completeness"

#define UIColorFromRGB(rgbValue) [UIColor colorWithRed:((float)((rgbValue & 0xFF0000) >> 16))/255.0 green:((float)((rgbValue & 0xFF00) >> 8))/255.0 blue:((float)(rgbValue & 0xFF))/255.0 alpha:1.0]
#define UIColorFromRGBA(rgbValue,alphaValue) [UIColor colorWithRed:((float)((rgbValue & 0xFF0000) >> 16))/255.0 green:((float)((rgbValue & 0xFF00) >> 8))/255.0 blue:((float)(rgbValue & 0xFF))/255.0 alpha:alphaValue]
UIColor *UIColorFromRGBHexStr(NSString *HexStr);
#define SCREEN_WIDTH ([UIScreen mainScreen].bounds.size.width)
#define SCREEN_HEIGHT ([UIScreen mainScreen].bounds.size.height)

                                                        
#define SYSTEM_VERSION_EQUAL_TO(v)                  ([[[UIDevice currentDevice] systemVersion] compare:v options:NSNumericSearch] == NSOrderedSame)
#define SYSTEM_VERSION_GREATER_THAN(v)              ([[[UIDevice currentDevice] systemVersion] compare:v options:NSNumericSearch] == NSOrderedDescending)
#define SYSTEM_VERSION_GREATER_THAN_OR_EQUAL_TO(v)  ([[[UIDevice currentDevice] systemVersion] compare:v options:NSNumericSearch] != NSOrderedAscending)
#define SYSTEM_VERSION_LESS_THAN(v)                 ([[[UIDevice currentDevice] systemVersion] compare:v options:NSNumericSearch] == NSOrderedAscending)
#define SYSTEM_VERSION_LESS_THAN_OR_EQUAL_TO(v)     ([[[UIDevice currentDevice] systemVersion] compare:v options:NSNumericSearch] != NSOrderedDescending)

       
void udwAlertOk(NSString *title,NSString *content);
void udwAlertOkForC(const char* title, const char* content);
void udwAlertOkSync(NSString *title, NSString *content);
void udwAlertOkWithVcWithCb(UIViewController *vc,NSString *title, NSString *content,void (^cb)());

UIViewController *udwGetTopViewController();
void udwAlertLoginWithUsernameAndPassword(NSString *title, NSString *content,
        void(^onLogin)(NSString* username,NSString* password),
                                          void(^onCancel)());
void udwAlertWithButtonAndCbList(NSString *title,NSString *content,NSArray* buttonNameList,NSArray* buttonCbList);

                 
void udwProgressViewShowDefault();
void udwProgressViewShow(NSString* msg);
void udwProgressViewHidden();

      
void udwOcJumpToVcFromAppDelegate(id<UIApplicationDelegate> appDelegate,UIViewController* vc);
void udwOcJumpToVcFromVc(UIViewController* from,UIViewController* to);

       
UIImage *udwGetImageByBase64String(NSString *base64);

               
typedef NS_ENUM(NSInteger, UILabelCountingMethod) {
    UILabelCountingMethodEaseInOut,
    UILabelCountingMethodEaseIn,
    UILabelCountingMethodEaseOut,
    UILabelCountingMethodLinear
};

typedef NSString* (^UICountingLabelFormatBlock)(CGFloat value);
typedef NSAttributedString* (^UICountingLabelAttributedFormatBlock)(CGFloat value);

@interface UICountingLabel : UILabel

@property (nonatomic, strong) NSString *format;
@property (nonatomic, assign) UILabelCountingMethod method;
@property (nonatomic, assign) NSTimeInterval animationDuration;

@property (nonatomic, copy) UICountingLabelFormatBlock formatBlock;
@property (nonatomic, copy) UICountingLabelAttributedFormatBlock attributedFormatBlock;
@property (nonatomic, copy) void (^completionBlock)(void);

-(void)countFrom:(CGFloat)startValue to:(CGFloat)endValue;
-(void)countFrom:(CGFloat)startValue to:(CGFloat)endValue withDuration:(NSTimeInterval)duration;

-(void)countFromCurrentValueTo:(CGFloat)endValue;
-(void)countFromCurrentValueTo:(CGFloat)endValue withDuration:(NSTimeInterval)duration;

-(void)countFromZeroTo:(CGFloat)endValue;
-(void)countFromZeroTo:(CGFloat)endValue withDuration:(NSTimeInterval)duration;

- (CGFloat)currentValue;
@end

          
#define udw_DEVICE_ISIPHONE        (UI_USER_INTERFACE_IDIOM() == UIUserInterfaceIdiomPhone)
#define udw_DEVICE_SIZE_ISIPHONE_4      (udw_DEVICE_ISIPHONE && [[UIScreen mainScreen] bounds].size.height == 480.0)
#define udw_DEVICE_SIZE_ISIPHONE_5      (udw_DEVICE_ISIPHONE && [[UIScreen mainScreen] bounds].size.height == 568.0)
#define udw_DEVICE_SIZE_ISIPHONE_6      (udw_DEVICE_ISIPHONE && [[UIScreen mainScreen] bounds].size.height == 667.0)
#define udw_DEVICE_SIZE_ISIPHONE_6PLUS (udw_DEVICE_ISIPHONE && [[UIScreen mainScreen] bounds].size.height == 736.0)
#define udw_DEVICE_SIZE_ISIPHONE_X      (udw_DEVICE_ISIPHONE && [[UIScreen mainScreen] bounds].size.height == 812.0)
#define udw_DEVICE_SIZE_ISRETINA        ([[UIScreen mainScreen] scale] == 2.0)
#define udw_DEVICE_SIZE_ISIPHONE_XS     (udw_DEVICE_ISIPHONE && [UIApplication sharedApplication].statusBarFrame.size.height >= 44)
#define udw_DEVICE_ISIPAD   [(NSString*)[UIDevice currentDevice].model hasPrefix:@"iPad"]
BOOL udwDeviceIsLandscape(void);

                              
CGFloat udwFixViewSizeWithOriginalFloatAndIsFlexibleHeight(CGFloat originalFloat,BOOL isHeight);
                                            
CGFloat udwFixViewSizeWithViewWidth(CGFloat width,CGSize size);
                                            
CGFloat udwFixViewSizeWithViewHeight(CGFloat height,CGSize size);
                                                                                   
void udwLabelFixSize(UILabel *label);
CGRect udwFrameChangeX(CGRect frame,CGFloat newX);
CGRect udwFrameChangeY(CGRect frame,CGFloat newY);
CGRect udwFrameChangeWidth(CGRect frame,CGFloat width);
CGRect udwFrameChangeHeight(CGRect frame,CGFloat height);
CGRect udwFrameChangeSize(CGRect frame,CGSize newSize);
                                                                                                                                                           
void udwViewAddShadow(UIView *addShadowView,UIColor *shadowColor,CGSize shadowOffset,float shadowOpacity,CGFloat shadowRadius);
                                                                                                                                                 
void udwViewDeleteShadow(UIView *shadowView);

                                                                                                                                                        
                                                                                                                                                                           
void udwRegisterNibClass(void);

                         
#define udwViewBorderRadius(View, Radius, Width, Color)\
\
[View.layer setCornerRadius:(Radius)];\
[View.layer setMasksToBounds:YES];\
[View.layer setBorderWidth:(Width)];\
[View.layer setBorderColor:[Color CGColor]]

              
#define udwViewRadius(View, Radius)\
\
[View.layer setCornerRadius:(Radius)];\
[View.layer setMasksToBounds:YES]


                
@interface kouLoadingButton : UIButton
@property (nonatomic,strong)UIColor *loadingColor;
- (void) LoadingStatus:(BOOL) loading;
@end

                   
@interface kouGetFontSizeLabel : UILabel
@property (nonatomic,assign,readonly)CGFloat fontSize;
@end

                     
@interface kouPointInsideButton : UIButton
@end

@interface udwClickTextModel : NSObject
@property (nonatomic,copy)  NSString * _Nonnull  clickStr;
@property (nonatomic,strong)UIFont   * _Nullable clickFont;
@property (nonatomic,strong)UIColor  * _Nullable clickStrColor;
@property (nonatomic,copy) void (^ _Nullable onClick)(void);
@end
udwClickTextModel* _Nonnull udwNewClickTextModel(NSString * _Nonnull clickStr,UIFont * _Nullable clickFont,UIColor * _Nullable clickStrColor, void(^ _Nullable onClick)(void));

                 
@interface udwSizeToFitLabelModel : NSObject
@property (nonatomic, copy) NSString * _Nullable labelStr;
@property (nonatomic, assign) CGFloat labelLineH;
@property (nonatomic, copy) void (^ __nullable attributedStringCallback)(NSMutableAttributedString * _Nullable attributedString);
@end
@interface udwSizeToFitLabel : UILabel
- (instancetype _Nonnull)initWithFrame:(CGRect)frame font:(UIFont *_Nonnull)font string:(NSString *_Nonnull)context offset:(CGFloat)lineOffset;
- (instancetype _Nonnull)initWithFrame:(CGRect)frame font:(UIFont *_Nonnull)font textColor:(UIColor *_Nonnull)textColor context:(NSString *_Nonnull)context lineOffset:(CGFloat)lineOffset differentColor:(UIColor *__nullable)differentColor colorRange:(NSRange)colorRange;
- (instancetype _Nonnull)initWithFrame:(CGRect)frame font:(UIFont *_Nonnull)font textColor:(UIColor *_Nonnull)textColor context:(NSString *_Nonnull)context lineOffset:(CGFloat)lineOffset attributedStringCallback:(void (^ __nullable)(NSMutableAttributedString * _Nullable attributedString))attributedStringCallback;
- (void)changeTextWithNewContext:(NSString *_Nonnull)context;
- (void)changeTextWithNewContext:(NSString *_Nonnull)context withUnderLine:(BOOL)withUnderLine;
- (void)changeLabelWithModel:(udwSizeToFitLabelModel *)model;
@end

                                                         
@interface udwSizeToFitTextView : UITextView
                                                             
- (instancetype _Nonnull)initWithFrame:(CGRect)frame font:(UIFont *_Nonnull)font textColor:(UIColor *_Nonnull)textColor context:(NSString *_Nonnull)context lineOffset:(CGFloat)lineOffset attributedStringCallback:(void (^ __nullable)(NSMutableAttributedString * _Nullable attributedString))attributedStringCallback;
               
- (void)changeTextWithNewContext:(NSString *)context;
                                                                                                 
       
                                                                                                                                                                                                                                                                                                                                                             

                                                                                                       
                                       
                                                                                         
                                       
        
- (instancetype _Nonnull)initWithFrame:(CGRect)frame font:(UIFont *_Nonnull)font textColor:(UIColor *_Nonnull)textColor context:(NSString *_Nonnull)context lineOffset:(CGFloat)lineOffset attributedStringCallback:(void (^ __nullable)(NSMutableAttributedString * _Nullable attributedString))attributedStringCallback clickTextArr:(NSArray<udwClickTextModel *> *_Nullable)clickTextArr;
@end

@interface UIView (Frame)
           
@property CGFloat udw_x;
           
@property CGFloat udw_y;
             
@property CGFloat udw_width;
             
@property CGFloat udw_height;
              
@property CGFloat udw_centerX;
              
@property CGFloat udw_centerY;
           
@property CGSize  udw_size;
             
@property CGPoint udw_origin;
@end

               
@interface udwGifImageView : UIImageView
                                                                                          
                                                                                                                                                                   
@property (nonatomic, copy) void(^loopCompletionBlock)(NSUInteger loopCountRemaining);
@property (nonatomic, strong, readonly) UIImage *currentFrame;
@property (nonatomic, assign, readonly) NSUInteger currentFrameIndex;
                                                                                                                                     
                                                                                                                                                                                                  
@property (nonatomic, copy) NSString *runLoopMode;
- (instancetype)initWithFrame:(CGRect)frame gifName:(NSString *)gifName;
@end

              
@interface udwDrawerVCManager : NSObject
@property (nonatomic, strong, readonly)UIViewController *drawerViewController;
@property (nonatomic, strong, readonly)UIViewController *mainViewController;
@property (nonatomic, strong, readonly)UIViewController *leftViewController;
@property (nonatomic, assign) CGFloat offset;
                                                                                                                                  
+ (instancetype)sharedManager;
- (UIViewController *)InitWithMainViewController:(UIViewController *)mainViewController leftViewController:(UIViewController *)leftViewController;
- (void)openLeftDrawer;
- (void)closeLeftDrawer;
                                            
- (void)setOpenLeftCallBack:(void (^)(void))CallBack;
                                            
- (void)setCloseLeftCallBack:(void (^)(void))CallBack;
                                    
- (void)rotateReloadUi:(BOOL)isLandscape;
@end

UIImageView *udwNavigationBarFindShadowImage(UINavigationBar *bar);
void udwLogAllFontName(void);

void KeepScreenOn();

@interface UIControl (Action)
- (void)udwActionHandleEvents:(UIControlEvents)controlEvents action:(void (^)(void))action;
@end

@interface UITapGestureRecognizer (Action)
- (void)udwActionWithAction:(void (^)(void))action;
@end

void UdwLogAlert(NSString *msg);
void UdwLogAlertClean();

@interface UIViewController (PresentVcCallback)

@property (nonatomic,copy) void (^udwPresentVcCallback)(UIViewController*_vc);

@end


void UdwAddLocalNotification(NSDate *date,NSString *noticeTitle,NSString *noticeContent,NSString *identifier);
void UdwCancelOneLocalNotification(NSString *identifier);
void UdwCancelAllLocalNotification();
typedef NS_ENUM(NSInteger, udwNotificationAuthorizationStatus) {
                                             
    udwNotificationAuthorizationStatusNotDetermined = 0,
                    
    udwNotificationAuthorizationStatusDenied,
                                
    udwNotificationAuthorizationStatusAuthorized,
};
                                                                                    
void udwGetNotificationStatus(void (^statusCallback)(udwNotificationAuthorizationStatus status));
void udwRegisterNotificationPermission();
NSString* udwDeviceTokenGetString(NSData *deviceToken);

AVAuthorizationStatus udwGetAudioStatus(void);
void udwRequestRecordPermission (void (^statusCallback)(bool allow));

AVAuthorizationStatus udwGetCameraStatus();
void udwRequestCameraPermission (void (^statusCallback)(bool allow));

CGFloat udwGetStatusBarHeight();
CGFloat udwGetNavigationBarMaxY();
                                        
@interface udwCoverDrawerVC : UIViewController
@property (nonatomic, strong, readonly) UIViewController *mainViewController;
@property (nonatomic, strong, readonly) UIViewController *rightViewController;
- (void)openRightDrawer:(void (^)(void))completionCallback;
- (void)closeRightDrawer:(void (^)(void))completionCallback;
- (BOOL)isDrawerOpen;
- (CGFloat)getDrawerWidth;
- (void)setEnableOpenGesture:(BOOL)enableOpenGesture;
- (void)setOpenRightCallBack:(void (^)(void))CallBack;
- (void)setCloseRightCallBack:(void (^)(void))CallBack;
@end
@interface udwCoverDrawerVCListManager : NSObject
+ (UIViewController *)InitWithMainViewController:(UIViewController *)mainViewController rightViewController:(UIViewController *)rightViewController tag:(NSString *)tag offset:(CGFloat)offset;
+ (udwCoverDrawerVC *)GetOneDrawer:(NSString *)tag;
@end
                                          
UIView* udwGetWindow();

                                                                                                                    
@interface udwTheme : NSObject
+ (void)startTheme:(NSString *)tag;
+ (NSString *)currentThemeTag;
+ (void)setDefaultTheme:(NSString *)tag;
+ (NSArray *)allThemeTag;
@end
@class udwThemeConfigModel;
typedef void(^udwThemeConfigBlock)(UIView *item);
typedef udwThemeConfigModel *(^udwThemeConfigForColor)(NSString *tag,UIColor *color);
typedef udwThemeConfigModel *(^udwThemeConfigForBackgroundImage)(NSString *tag,UIImage *image);
typedef udwThemeConfigModel *(^udwThemeConfigForColorAndState)(NSString *tag,UIColor *color,UIControlState state);
typedef udwThemeConfigModel *(^udwThemeConfigForImageAndState)(NSString *tag,UIImage *image,UIControlState state);
typedef void(^udwThemeChangeThemeBlock)(NSString *tag,UIView *item);
typedef udwThemeConfigModel *(^udwThemeConfigForCustomConfigBlock)(NSString *tag , udwThemeConfigBlock);
typedef udwThemeConfigModel *(^udwThemeConfigForThemeChangBlock)(udwThemeChangeThemeBlock);
@interface udwThemeConfigModel : NSObject
                         
@property (nonatomic, copy, readonly)udwThemeConfigForColor udwAddBackGroundColor;
@property (nonatomic, copy, readonly)udwThemeConfigForBackgroundImage udwAddImage;
                                
@property (nonatomic, copy, readonly)udwThemeConfigForColor udwAddLabelTitleColor;
@property (nonatomic, copy, readonly)udwThemeConfigForColorAndState udwAddButtonTitleColor;
@property (nonatomic, copy, readonly)udwThemeConfigForImageAndState udwAddButtonImage;
@property (nonatomic, copy, readonly)udwThemeConfigForImageAndState udwAddButtonBackgroundImage;
                                     
@property (nonatomic, copy, readonly)udwThemeConfigForCustomConfigBlock udwAddCustomConfig;
                               
@property (nonatomic, copy, readonly)udwThemeConfigForThemeChangBlock udwThemeChangBlock;
@end
@interface UIView (udwThemeConfig)
                                                                                                     
@property (nonatomic, strong) udwThemeConfigModel *udwThemeModel;
@end

@interface UIScrollView (udwTouchEvent)
@property (nonatomic, assign) BOOL udwScrollViewTouchStrikeEnabled;
@end

                                                           
void udwSetAdjustsFontForWidth(UIView *view);

                                              
BOOL udwDeviceHaveHomeButton();

@class udwPopView;
typedef void(^WillDismissBlock)(udwPopView *popView);
typedef void(^DidDismissBlock)(udwPopView *popView);
typedef void(^WillShowBlock)(udwPopView *popView);
typedef void(^DidShowBlock)(udwPopView *popView);

typedef NS_ENUM(NSUInteger, udwPopViewAnimationStyle) {
    udwPopViewAnimationStyleTopToBottom = 0,
    udwPopViewAnimationStyleBottomToTop,
    udwPopViewAnimationStyleLeftToRight,
    udwPopViewAnimationStyleRightToLeft,
    udwPopViewAnimationStyleFade,                           
    udwPopViewAnimationStyleScale,
};

@interface udwPopView : UIView

   
        
                                                                                                      
                               
                        
   

- (instancetype)initWithSubView:(__kindof UIView *)subView;

   
        
                          
   
- (void)removeSelf;

   
        
                                          
                            
   
- (void)showPopViewOn:(UIView *)view;


                                  
@property (assign, nonatomic) BOOL animatedEnable;
                                                                         
@property (nonatomic, assign) BOOL autoRemoveEnable;
                                                    
@property (assign, nonatomic) BOOL adjustedKeyboardEnable;
                                      
@property (assign, nonatomic) NSTimeInterval duration;
                                   
@property (assign, nonatomic) udwPopViewAnimationStyle animationStyle;

                                          
@property (nonatomic, strong) WillDismissBlock willDismiss;
                                         
@property (nonatomic, strong) DidDismissBlock didDismiss;
                                       
@property (nonatomic, strong) WillShowBlock willShow;
                                      
@property (nonatomic, strong) DidShowBlock didShow;



@end
