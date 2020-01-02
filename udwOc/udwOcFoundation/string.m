// +build ios macAppStore

#import "udwOcFoundation.h"

NSData *udwStringToNSData(NSString *s){
    return [s dataUsingEncoding:NSUTF8StringEncoding];
}

NSString *udwNSDataToString(NSData *data){
    return [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
}

NSString* udwItoA64(int64_t num){
    return [NSString stringWithFormat:@"%lld",num];
}

int64_t udwAtoI64(NSString* s,NSError **err){
    int64_t out;
    BOOL ret = [[NSScanner scannerWithString:s] scanLongLong:&out];
    if (ret){
        return out;
    }else{
        *err = udwNSError(@"[jbff7y8fuj] can not parse int");
        return 0;
    }
}

BOOL udwStringIsNilOrEmpty(NSString *s){
    return s==nil || s.length==0;
}

                                                                                                                               
                                       
NSString* udwCStringToNSString(char *s){
    return [[NSString alloc] initWithCString:s encoding:NSUTF8StringEncoding];
}

                                                                                                                    
char * udwNSStringToCString(NSString* str){
    return (char *)([str UTF8String]);
}

bool udwStringIsEqual(NSString *a,NSString* b){
    if (udwStringIsNilOrEmpty(a) && udwStringIsNilOrEmpty(b)){
        return true;
    }
    if (a==nil || b==nil){
        return false;
    }
    return [a isEqualToString:b];
}

BOOL udwStringContainString(NSString* superString, NSString *subString) {
    if ([superString rangeOfString: subString].location == NSNotFound) {
        return false;
    }
    return true;
}

NSString* udwStringAdd2(NSString* a,NSString* b){
    if (a==nil){
        return b;
    }
    if (b==nil){
        return a;
    }
    return [a stringByAppendingString:b];
}