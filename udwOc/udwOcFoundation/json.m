// +build ios macAppStore

#import "udwOcFoundation.h"

id UdwJSONDecodeWithString(NSString *inS,NSError **err){
    return [NSJSONSerialization JSONObjectWithData:[inS dataUsingEncoding:NSUTF8StringEncoding]
                                            options:0 error:err];
}

                                        
id UdwJSONDecode(NSData *data,NSError **err){
    return [NSJSONSerialization JSONObjectWithData:data
                                            options:NSJSONReadingAllowFragments error:err];
}

                                                                             
                                                   
NSData* UdwJSONEncode(id obj,NSError **err){
    if ([obj isKindOfClass:[NSString class]]){
        NSData* data = [NSJSONSerialization dataWithJSONObject:@[obj] options:0 error:err];
        if (*err!=nil){
            return nil;
        }
                                                     
        return [data subdataWithRange:NSMakeRange(1, data.length-2)];
    }

    return [NSJSONSerialization dataWithJSONObject:obj options:0 error:err];
}

                                                                             
NSString* UdwJSONEncodeToString(id obj,NSError **err){
    if (err!=nil) {
        return udwNSDataToString(UdwJSONEncode(obj, err));
    }else{
        NSError *err1;
        return udwNSDataToString(UdwJSONEncode(obj, &err1));
    }
}

