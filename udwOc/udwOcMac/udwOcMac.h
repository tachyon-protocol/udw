// +build darwin,!ios

#import <Cocoa/Cocoa.h>
#import "NSView+Add.h"

void udwShowLocalNotification(NSString *context);
void udwOpenUrlWithDefaultBrowser(NSString *url);
BOOL udwCheckIsRunningApp(NSString *bundleId);
void udwApplicationDidBecomeActive();
int  udwGetProcessIdByName(NSString *processName);

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

@interface udwPopView : NSView

   
        
                                                                                                      
                               
                        
   

- (instancetype)initWithSubView:(__kindof NSView *)subView;

   
        
                          
   
- (void)removeSelf;

   
        
                                          
                            
   
- (void)showPopViewOn:(NSView *)view;

                                  
@property (assign, nonatomic) BOOL animatedEnable;
                                                                         
@property (nonatomic, assign) BOOL autoRemoveEnable;
                                      
@property (assign, nonatomic) NSTimeInterval duration;
                                   
@property (assign, nonatomic) udwPopViewAnimationStyle animationStyle;

                                          
@property (nonatomic, strong) WillDismissBlock willDismiss;
                                         
@property (nonatomic, strong) DidDismissBlock didDismiss;
                                       
@property (nonatomic, strong) WillShowBlock willShow;
                                      
@property (nonatomic, strong) DidShowBlock didShow;



@end
