#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"

void KeepScreenOn () {
    udwRunOnMainThread(^{
        [[UIApplication sharedApplication]setIdleTimerDisabled:YES];
    });
}
