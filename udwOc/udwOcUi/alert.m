#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"

UIViewController *udwGetTopViewController() {
    UIWindow *window = [[UIApplication sharedApplication] keyWindow];
    if (window.windowLevel != UIWindowLevelNormal) {
        NSArray *windows = [[UIApplication sharedApplication] windows];
        for (window in windows) {
            if (window.windowLevel == UIWindowLevelNormal) {
                break;
            }
        }
    }

    for (UIView *subView in [window subviews]) {
        UIResponder *responder = [subView nextResponder];
                                                                                                                                 
         if ([responder isEqual:window]) {
                                     
             if ([[subView subviews] count]) {
                 UIView *subSubView = [subView subviews][0];                                           
                 if ([NSStringFromClass(subSubView.class) isEqualToString:@"UIDropShadowView"]) {              
                     if ([[subSubView subviews] count]) {
                         subSubView = subSubView.subviews[0];
                      }
                  }             
                 responder = [subSubView nextResponder];
             }
         }
        if ([responder isKindOfClass:[UIViewController class]]) {
            UIViewController *controller = (UIViewController *)responder;
            BOOL isPresenting;
            do {
                                                                                                
                UIViewController *presented = [controller presentedViewController];
                isPresenting = presented != nil;
                if (presented != nil) {
                    controller = presented;
                }

            } while (isPresenting);

            return controller;
        }
    }
    return nil;
}

                                               
                                
void udwAlertOk(NSString *title, NSString *content) {
    udwAlertOkWithVcWithCb(udwGetTopViewController(),title,content,nil);
}

                                       
void UdwAlertOkV2(NSString *title, NSString *content) {
    udwAlertOkWithVcWithCb(udwGetTopViewController(),title,content,nil);
}

void udwAlertOkForC(const char* title, const char* content) {
    NSString* _title = [[NSString alloc] initWithCString:title encoding:NSUTF8StringEncoding];
    NSString* _content = [[NSString alloc] initWithCString:content encoding:NSUTF8StringEncoding];
    udwAlertOk(_title,_content);
}

                                                
                                     
void udwAlertOkSync(NSString *title, NSString *content) {
    udwSignalChan* chan = udwNewSignalChan();
    udwAlertOkWithVcWithCb(udwGetTopViewController(),title,content,^{
        [chan send];
    });
    [chan recv];
}

                                                        
                            
void udwAlertOkWithVcWithCb(UIViewController *vc,NSString *title, NSString *content,void (^cb)()) {
    udwRunOnMainThread(^{
        UIAlertController *alert = [UIAlertController alertControllerWithTitle:title
                                                                       message:content
                                                                preferredStyle:UIAlertControllerStyleAlert];

        UIAlertAction *defaultAction = [UIAlertAction actionWithTitle:@"OK" style:UIAlertActionStyleDefault
                                                              handler:^(UIAlertAction *action) {
                                                                  if (cb!=nil){
                                                                      cb();
                                                                  }
                                                              }];
        [alert addAction:defaultAction];
        [vc presentViewController:alert animated:YES completion:nil];
    });
}

void udwAlertWithButtonAndCbList(NSString *title,NSString *content,NSArray* buttonNameList,NSArray* buttonCbList){
    udwRunOnMainThread(^{
        UIAlertController *alert = [UIAlertController alertControllerWithTitle:title
                                                                       message:content
                                                                preferredStyle:UIAlertControllerStyleAlert];

        for (int i=0;i<buttonNameList.count;i++) {
            UIAlertAction *action = [UIAlertAction actionWithTitle:buttonNameList[i] style:UIAlertActionStyleDefault
                                                                  handler:^(UIAlertAction *action) {
                                                                      ((void (^)())buttonCbList[i])();
                                                                  }];
            [alert addAction:action];
        }
        [udwGetTopViewController() presentViewController:alert animated:YES completion:nil];
    });
}

                    
                                                                                                              
void udwAlertLoginWithUsernameAndPassword(NSString *title, NSString *content,
        void(^onLogin)(NSString* username,NSString* password),
        void(^onCancel)()) {
    udwRunOnMainThread(^{
        UIAlertController *alert = [UIAlertController alertControllerWithTitle:title
                                                                       message:content
                                                                preferredStyle:UIAlertControllerStyleAlert];

        [alert addTextFieldWithConfigurationHandler:^(UITextField *textField){
            textField.placeholder = @"UserName";
        }];
        [alert addTextFieldWithConfigurationHandler:^(UITextField *textField){
            textField.placeholder = @"Password";
            textField.secureTextEntry = true;
        }];
        [alert addAction:[UIAlertAction actionWithTitle:@"Login" style:UIAlertActionStyleDefault
                                                handler:^(UIAlertAction *action) {
                                                    if (onLogin !=nil){
                                                        onLogin(((UITextField *)alert.textFields[0]).text,
                                                                ((UITextField *)alert.textFields[1]).text);
                                                    }
                                                }]];
        [alert addAction:[UIAlertAction actionWithTitle:@"Cancel" style:UIAlertActionStyleDefault
                                                handler:^(UIAlertAction *action) {
                                                    if (onCancel !=nil){
                                                        onCancel();
                                                    }
                                                }]];
        [udwGetTopViewController() presentViewController:alert animated:YES completion:nil];
    });
}

UIView* udwGetWindow(){
    UIWindow* window = [UIApplication sharedApplication].keyWindow;
    if (!window) {
        window = [[UIApplication sharedApplication].windows objectAtIndex:0];
    }
    return window;
}