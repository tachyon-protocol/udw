// +build darwin,!ios

#import "udwOcMac.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"
#import "zzzig_udwRpcGoAndObjectiveC.h"
#include <errno.h>
#include <stdbool.h>
#include <stdlib.h>
#include <stdio.h>
#include <sys/sysctl.h>

void udwShowLocalNotification(NSString *context){
    NSUserNotification *notification = [[NSUserNotification alloc] init];
    [notification setInformativeText:context];
    [notification setDeliveryDate:[NSDate dateWithTimeIntervalSinceNow:0]];
    [notification setSoundName:NSUserNotificationDefaultSoundName];
    [notification setHasActionButton:true];
    [notification setActionButtonTitle:@"close"];
    [notification setValue:@YES forKey:@"_showsButtons"];
    NSUserNotificationCenter *center = [NSUserNotificationCenter defaultUserNotificationCenter];
    [center scheduleNotification:notification];
}

void udwOpenUrlWithDefaultBrowser(NSString *url){
    [[NSWorkspace sharedWorkspace] openURL:[NSURL URLWithString:url]];
}

BOOL udwCheckIsRunningApp(NSString *bundleId){
    NSArray<NSRunningApplication *> *appList = [[NSWorkspace sharedWorkspace] runningApplications];
    for (int i = 0; i<appList.count; i++) {
        NSRunningApplication * app = appList[i];
        if (udwStringIsEqual(bundleId, app.bundleIdentifier)) {
            return true;
        }
    }
    return false;
}

void udwApplicationDidBecomeActive(){
    applicationDidBecomeActiveCall();
}

typedef struct kinfo_proc kinfo_proc;
int udwGetProcessIdByName(NSString *processName){
    if (processName==nil || processName.length==0) {
        return 0;
    }
    int                 err;
    kinfo_proc *        result;
    bool                done;
    static const int    name[] = { CTL_KERN, KERN_PROC, KERN_PROC_ALL, 0 };
    size_t              length;
    result = NULL;
    done = false;
    do {
        length = 0;
        err = sysctl( (int *) name, (sizeof(name) / sizeof(*name)) - 1,
                     NULL, &length,
                     NULL, 0);
        if (err == -1) {
            err = errno;
        }
        if (err == 0) {
            result = malloc(length);
            if (result == NULL) {
                err = ENOMEM;
            }
        }
        if (err == 0) {
            err = sysctl( (int *) name, (sizeof(name) / sizeof(*name)) - 1,
                         result, &length,
                         NULL, 0);
            if (err == -1) {
                err = errno;
            }
            if (err == 0) {
                done = true;
            } else if (err == ENOMEM) {
                free(result);
                result = NULL;
                err = 0;
            }
        }
    } while (err == 0 && ! done);
    if (err != 0 && result != NULL) {
        free(result);
        result = NULL;
    }
    for (int i=0; i<length / sizeof(kinfo_proc); i++) {
        kinfo_proc onePrecess = result[i];
        if ([[NSString stringWithFormat:@"%s",onePrecess.kp_proc.p_comm] isEqualToString:processName]) {
            return onePrecess.kp_proc.p_pid;
        }
    }
    return 0;
}
