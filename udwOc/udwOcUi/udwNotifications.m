#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"

void UdwAddLocalNotification(NSDate *date,NSString *noticeTitle,NSString *noticeContent,NSString *identifier) {
    NSTimeInterval time = [date timeIntervalSinceDate:[NSDate dateWithTimeIntervalSinceNow:0]];
    if (time < 0 ) {
        return;
    }
    if (SYSTEM_VERSION_GREATER_THAN_OR_EQUAL_TO(@"10.0")) {
        UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];

                                                                                                                                                          
        UNMutableNotificationContent* content = [[UNMutableNotificationContent alloc] init];
        content.title = [NSString localizedUserNotificationStringForKey:noticeTitle arguments:nil];
        content.body = [NSString localizedUserNotificationStringForKey:noticeContent arguments:nil];
        content.sound = [UNNotificationSound defaultSound];

        UNTimeIntervalNotificationTrigger* trigger = [UNTimeIntervalNotificationTrigger
                                                      triggerWithTimeInterval:time repeats:NO];

        UNNotificationRequest* request = [UNNotificationRequest requestWithIdentifier:identifier
                                                                              content:content trigger:trigger];
        [center addNotificationRequest:request withCompletionHandler:^(NSError * _Nullable error) {

        }];
    } else {
        UILocalNotification *notification = [[UILocalNotification alloc] init];
        notification.fireDate = date;
        notification.timeZone = [NSTimeZone defaultTimeZone];
        notification.repeatInterval = 0;                  
        notification.alertTitle = noticeTitle;
        notification.alertBody =  noticeContent;
        notification.applicationIconBadgeNumber = 1;
        notification.soundName = UILocalNotificationDefaultSoundName;
        notification.userInfo = @{identifier:identifier};
        [[UIApplication sharedApplication] scheduleLocalNotification:notification];
    }
}

void UdwCancelOneLocalNotification(NSString *identifier) {
    if (SYSTEM_VERSION_GREATER_THAN_OR_EQUAL_TO(@"10.0")) {
            UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
            [center removePendingNotificationRequestsWithIdentifiers:@[identifier]];
        } else {
            UIApplication *app = [UIApplication sharedApplication];
            NSArray *eventArray = [app scheduledLocalNotifications];
            for (int i=0; i<[eventArray count]; i++)
            {
                UILocalNotification* oneEvent = [eventArray objectAtIndex:i];
                NSDictionary *userInfoCurrent = oneEvent.userInfo;
                NSString *uid=[NSString stringWithFormat:@"%@",[userInfoCurrent valueForKey:identifier]];
                if ([uid isEqualToString:identifier])
                {
                    [app cancelLocalNotification:oneEvent];
                    break;
                }
            }
        }
}

void UdwCancelAllLocalNotification() {
    if (SYSTEM_VERSION_GREATER_THAN_OR_EQUAL_TO(@"10.0")) {
        UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
        [center removeAllPendingNotificationRequests];
    } else {
        [[UIApplication sharedApplication] cancelAllLocalNotifications];
    }
}

void udwGetNotificationStatus(void (^statusCallback)(udwNotificationAuthorizationStatus status)) {
    if (@available(iOS 10.0, *)) {
        UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
        [center getNotificationSettingsWithCompletionHandler:^(UNNotificationSettings * _Nonnull settings) {
            udwRunOnMainThread(^{
                    switch (settings.authorizationStatus) {
                    case UNAuthorizationStatusNotDetermined:
                        statusCallback(udwNotificationAuthorizationStatusNotDetermined);
                        break;
                    case UNAuthorizationStatusDenied:
                        statusCallback(udwNotificationAuthorizationStatusDenied);
                        break;
                    case UNAuthorizationStatusAuthorized:
                        statusCallback(udwNotificationAuthorizationStatusAuthorized);
                        break;
                    default:
                        break;
                }
            });
        }];
    } else {
        UIUserNotificationSettings *settings = [[UIApplication sharedApplication] currentUserNotificationSettings];
        if (settings.types == UIUserNotificationTypeNone) {
            statusCallback(udwNotificationAuthorizationStatusNotDetermined);
        } else{
            statusCallback(udwNotificationAuthorizationStatusAuthorized);
        }
    }
}

void udwRegisterNotificationPermission(){
    if (@available(iOS 10.0, *)) {
        UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
        [center requestAuthorizationWithOptions:(UNAuthorizationOptionAlert + UNAuthorizationOptionSound + UNAuthorizationOptionBadge)completionHandler:^(BOOL granted, NSError * _Nullable error) {

        }];
    } else {
        if ([[UIApplication sharedApplication] respondsToSelector:@selector(registerUserNotificationSettings:)]) {
            UIUserNotificationType type =  UIUserNotificationTypeAlert | UIUserNotificationTypeBadge | UIUserNotificationTypeSound;
            UIUserNotificationSettings *settings = [UIUserNotificationSettings settingsForTypes:type
                                                                                     categories:nil];
            [[UIApplication sharedApplication] registerUserNotificationSettings:settings];
        }
    }
}

NSString* udwDeviceTokenGetString(NSData *deviceToken){
    NSUInteger length = deviceToken.length;
    if (length == 0) {
        return nil;
    }
    const unsigned char *buffer = deviceToken.bytes;
    NSMutableString *hexString  = [NSMutableString stringWithCapacity:(length * 2)];
    for (int i = 0; i < length; ++i) {
        [hexString appendFormat:@"%02x", buffer[i]];
    }
    return [hexString copy];
}
