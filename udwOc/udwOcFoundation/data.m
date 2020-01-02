// +build ios macAppStore

#import "udwOcFoundation.h"

NSData* udwNSDataSub(NSData *inData, unsigned long startPos,unsigned long endPos){
    return [inData subdataWithRange:NSMakeRange( startPos, endPos-startPos)];
}

NSData* udwNSDataFromBase64(NSString *s){
    return [[NSData alloc] initWithBase64EncodedData:[s dataUsingEncoding:NSUTF8StringEncoding] options:0];
}

                                                 
                                                                                       
   

NSString* udwNSDataToBase64(NSData *s){
    return [s base64EncodedStringWithOptions:0];
}

NSString* udwByteUnitToString(int64_t s){
    uint64_t si = (uint64_t) llabs(s);
    if (si>=1024ll*1024*1024*1024*1024){
        return [NSString stringWithFormat:@"%.2fPB",((double)s)/((double)1024*1024*1024*1024*1024)];
    }else if(si>=1024ll*1024*1024*1024){
        return [NSString stringWithFormat:@"%.2fTB",((double)s)/((double)1024*1024*1024*1024)];
    }else if(si>=1024ll*1024*1024){
        return [NSString stringWithFormat:@"%.2fGB",((double)s)/((double)1024*1024*1024)];
    }else if(si>=1024ll*1024){
        return [NSString stringWithFormat:@"%.2fMB",((double)s)/((double)1024*1024)];
    }else if(si>=1024){
        return [NSString stringWithFormat:@"%.2fKB",((double)s)/((double)1024)];
    }else{
        return [NSString stringWithFormat:@"%lldB",s];
    }
}