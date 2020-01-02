// +build ios macAppStore

#import <Foundation/Foundation.h>

NSString* UdwOcUserDefaultsKvGet(NSString* key){
    return [[NSUserDefaults standardUserDefaults] objectForKey:key];
}

void UdwOcUserDefaultsKvSet(NSString* key,NSString* value){
    [[NSUserDefaults standardUserDefaults] setObject:value forKey:key];
    BOOL succ = [[NSUserDefaults standardUserDefaults] synchronize];
    if (!succ){
        NSLog(@"[UdwOcUserDefaultsKvSet] [[NSUserDefaults standardUserDefaults] synchronize] !succ");
    }
}

void UdwUserDefaultsKvSetForC(const char* key, const char* value) {
    NSString* _key = [[NSString alloc] initWithCString:key encoding:NSUTF8StringEncoding];
    NSString* _value = [[NSString alloc] initWithCString:value encoding:NSUTF8StringEncoding];
    [[NSUserDefaults standardUserDefaults] setObject:_value forKey:_key];
    BOOL succ = [[NSUserDefaults standardUserDefaults] synchronize];
    if (!succ){
       NSLog(@"[UdwOcUserDefaultsKvSet] [[NSUserDefaults standardUserDefaults] synchronize] !succ");
    }
}