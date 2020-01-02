// +build ios macAppStore

#import "udwOcFoundation.h"

NSString *udwNSTimeIntervalFormat(NSTimeInterval uptime){
    unsigned int seconds = (unsigned int)floor(uptime);
    return [NSString stringWithFormat:@"%02u:%02u:%02u", seconds / 3600, (seconds / 60) % 60, seconds % 60];
}

NSDate*udwGolangTimeToNsDate(NSString *timeS){
    if (timeS==nil){
        return nil;
    }
    NSDateFormatter *dateFormatter = [[NSDateFormatter alloc] init];
    NSDate *capturedStartDate;
                                                  
    [dateFormatter setDateFormat:@"yyyy-MM-dd'T'HH:mm:ss.SSSSSSSSSZZZZZ"];
    capturedStartDate = [dateFormatter dateFromString: timeS];
    if (capturedStartDate==nil){
                                            
        [dateFormatter setDateFormat:@"yyyy-MM-dd'T'HH:mm:ssZZZZZ"];
        capturedStartDate = [dateFormatter dateFromString: timeS];
    }
    if (capturedStartDate==nil){
        NSLog(@"[6yvyg54432] can not parse GolangTime %@",timeS);
        return nil;
    }
                                                                                                 
    if (capturedStartDate.timeIntervalSince1970<0){
                      
        return nil;
    }
    return capturedStartDate;
}

                                                     
           
                                                                                        
                                
BOOL udwTimeAfter(NSDate* after,NSDate* before){
    if (before==nil){
        return (after!=nil);
    }
    if (after==nil){
        return false;
    }
    return [after compare:before]==NSOrderedDescending;
}

NSDate* udwTimeNSDateNow(){
    return [NSDate date];
}

NSDate* udwTimeNSDateAddDuration(NSDate* in,NSTimeInterval dur){
    return [in dateByAddingTimeInterval:dur];
}

NSString *udwNSDateToLocalMysqlTime(NSDate* time){
    NSDateFormatter *df = [NSDateFormatter new];
    [df setDateFormat:@"yyyy-MM-dd HH:mm:ss"];
    df.timeZone = [NSTimeZone timeZoneForSecondsFromGMT:[NSTimeZone localTimeZone].secondsFromGMT];
    return [df stringFromDate:time];
}

BOOL udwTimeIsSameDay(NSDate* a,NSDate* b,NSTimeZone* timeZone){
    NSDateFormatter *df = [NSDateFormatter new];
    [df setDateFormat:@"yyyy-MM-dd"];
    df.timeZone = timeZone;
    return [[df stringFromDate:a] isEqualToString:[df stringFromDate:b]];
}

void udwSleep(NSTimeInterval dur){
    [NSThread sleepForTimeInterval:dur];
}

NSTimeInterval udwTimeSinceNow(NSDate* time){
    return [[NSDate new] timeIntervalSinceDate:time];
}

NSString *UdwTimeGetTimeZoneName(){
    NSTimeZone *zone = [NSTimeZone systemTimeZone];
    return zone.name;
}

NSInteger UdwTimeGetTimeZoneOffset(){
    NSTimeZone *zone = [NSTimeZone systemTimeZone];
    NSDate *date = [NSDate date];
    return [zone secondsFromGMTForDate:date];
}