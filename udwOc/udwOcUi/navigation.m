#import "udwOcUi.h"

UIImageView *findShadowImage(UIView *view){
    if ([view isKindOfClass:[UIImageView class]] && view.bounds.size.height <= 1) {
        return (UIImageView *)view;
    }
    for (UIView *subview in view.subviews) {
        UIImageView *imageView = findShadowImage(subview);
        if (imageView) {
            return imageView;
        }
    }
    return nil;
}

UIImageView *udwNavigationBarFindShadowImage(UINavigationBar *bar){
    return findShadowImage(bar);
}

CGFloat udwGetStatusBarHeight(){
    return [UIApplication sharedApplication].statusBarFrame.size.height;
}

CGFloat udwGetNavigationBarMaxY(){
    return udwGetStatusBarHeight()+44;
}