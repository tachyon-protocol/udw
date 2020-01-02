// +build ios macAppStore

#import "udwOcFoundation.h"

                             
void udwAsyncRun(void (^handler)()){
    dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), handler);
}

void udwRunOnMainThread(void (^handler)()){
    if([NSThread isMainThread]){
        handler();
    } else {
        dispatch_async(dispatch_get_main_queue(), handler);
    }
}

bool udwRunIsMainThread(){
    return [NSThread isMainThread];
}

@interface udwThreadCallback: NSObject
@property (strong, nonatomic) void(^handler)();
@end

@implementation udwThreadCallback
- (void)call{
    self.handler();
}
@end

                                                                            
                            
NSTimer *udwRunWithTimeIntervalOnMainThread(NSTimeInterval seconds,void(^handler)()){
    udwThreadCallback *cb = [udwThreadCallback new];
    cb.handler = handler;
    NSTimer *timer = [NSTimer timerWithTimeInterval:seconds target:cb selector:@selector(call) userInfo:nil repeats:YES];
    [[NSRunLoop mainRunLoop] addTimer:timer forMode:NSDefaultRunLoopMode];
    return timer;
}

@interface udwSignalChan()
@property (strong, nonatomic) dispatch_semaphore_t sysChan;
@end

@implementation udwSignalChan
- (void) send{
    dispatch_semaphore_signal(self.sysChan);
}
- (void) recv{
    dispatch_semaphore_wait(self.sysChan, DISPATCH_TIME_FOREVER);
}
@end
udwSignalChan* udwNewSignalChan(){
    udwSignalChan* chan = [udwSignalChan new];
    chan.sysChan =dispatch_semaphore_create(0);
    return chan;
}
