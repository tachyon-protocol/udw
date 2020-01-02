#import "udwOcUi.h"

void udwOcJumpToVcFromAppDelegate(id<UIApplicationDelegate> appDelegate,UIViewController* vc){
    appDelegate.window = [[UIWindow alloc] initWithFrame:[[UIScreen mainScreen] bounds]];
    appDelegate.window.backgroundColor = [UIColor whiteColor];                                                     
    appDelegate.window.rootViewController = vc;
    [appDelegate.window makeKeyAndVisible];
}

void udwOcJumpToVcFromVc(UIViewController* from,UIViewController* to){
    [from presentViewController:to animated:YES completion:nil];
}