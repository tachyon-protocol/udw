// +build ios macAppStore

#import "udwOcFoundation.h"

NSError *udwNSError(NSString * inS){
    if (inS==nil){
        return nil;
    }
    return [NSError errorWithDomain:@"udw" code:-1 userInfo:@{ NSLocalizedDescriptionKey : inS }];
}

NSError *udwNSErrorWithCode(NSString * inS,NSInteger code){
    if (inS==nil){
        return nil;
    }
    return [NSError errorWithDomain:@"udw" code:code userInfo:@{ NSLocalizedDescriptionKey : inS }];
}

NSError *udwNSErrorf(NSString * format, ...){
    va_list args;
    va_start(args, format);
    NSString *s = [[NSString alloc] initWithFormat:format arguments:args];
    va_end(args);
    return udwNSError(s);
}

NSString* udwNSErrorToString(NSError* error){
    if (error==nil){
        return nil;
    }
    return [NSString stringWithFormat:@"%@ (code=%d)",[error localizedDescription],(int)error.code];
}