#import "udwOcUi.h"

BOOL udwDeviceIsLandscape(){
  return [UIScreen mainScreen].bounds.size.width>[UIScreen mainScreen].bounds.size.height;
}

BOOL udwDeviceHaveHomeButton() {
    if (@available(iOS 11.0, *))  {
        if ([UIApplication sharedApplication].keyWindow.safeAreaInsets.bottom > 0) {
            return true;
        }
    }
    return false;
}
