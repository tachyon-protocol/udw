// +build ios macAppStore

#import <Foundation/Foundation.h>
#import <sys/sysctl.h>


        
void udwAsyncRun(void (^handler)());
void udwRunOnMainThread(void (^handler)());
bool udwRunIsMainThread();
NSTimer *udwRunWithTimeIntervalOnMainThread(NSTimeInterval seconds,void(^handler)());
@interface udwSignalChan: NSObject
- (void) send;
- (void) recv;
@end
udwSignalChan* udwNewSignalChan();
      
id UdwJSONDecodeWithString(NSString *inS,NSError **err);
id UdwJSONDecode(NSData *data,NSError **err);
NSData* UdwJSONEncode(id obj,NSError **err);
NSString* UdwJSONEncodeToString(id obj,NSError **err);

       
NSError *udwNSError(NSString * inS);
NSError *udwNSErrorf(NSString * format, ...);
NSError *udwNSErrorWithCode(NSString * inS,NSInteger code);
NSString* udwNSErrorToString(NSError* error);

      
NSData *udwNSDataSub(NSData *inData, unsigned long startPos,unsigned long endPos);
NSString *udwByteUnitToString(int64_t s);
NSData *udwNSDataFromBase64(NSString *s);
NSString *udwNSDataToBase64(NSData *s);

        
NSData *udwStringToNSData(NSString *s);
NSString *udwNSDataToString(NSData *data);
NSString* udwItoA64(int64_t num);
int64_t udwAtoI64(NSString* s,NSError **err);
BOOL udwStringIsNilOrEmpty(NSString *s);
NSString* udwCStringToNSString(char *s);
char * udwNSStringToCString(NSString* str);
                                     
bool udwStringIsEqual(NSString *a,NSString* b);
BOOL udwStringContainString(NSString* superString, NSString *subString);

         
                                               
                                                    
                                                                  
                                             
                       

                
                              
                         

       
                                                            
                                       
                                                  

      
NSString *udwNSTimeIntervalFormat(NSTimeInterval uptime);
NSDate*udwGolangTimeToNsDate(NSString *timeS);
BOOL udwTimeAfter(NSDate*after,NSDate*before);
NSString *udwNSDateToLocalMysqlTime(NSDate* time);
BOOL udwTimeIsSameDay(NSDate* a,NSDate* b,NSTimeZone* timeZone);
NSDate* udwTimeNSDateNow();
NSDate* udwTimeNSDateAddDuration(NSDate* in,NSTimeInterval dur);
void udwSleep(NSTimeInterval dur);

                 
                                          
                               

        
                                                                          
                           
                           
                              
                              
                               
                                                  
                     
      
                                                         
                          

      
                                                                                           

               
NSString* UdwOcUserDefaultsKvGet(NSString* key);
void UdwOcUserDefaultsKvSet(NSString* key,NSString* value);
void UdwUserDefaultsKvSetForC(const char* key, const char* value);

NSString* UdwOcGetAppName();

           
void udwOcTryCatch(void (^tryCallback)(),void (^catchCallback)(NSException *exception),void (^finallyCallback)());

                                                                                             
#ifdef DEBUG
#define udwLog(...) NSLog(__VA_ARGS__)
#else
#define udwLog(...)
#endif

#define udwLogFunc udwLog(@"--------------%@:%d%s", [[NSString stringWithUTF8String:__FILE__] lastPathComponent], __LINE__, __func__)
#define udwLogFilePath udwLog(@"%s", __FILE__)

double udwRandomBetween0And1();

               
@interface udwWeakProxy : NSProxy
@property (nullable, nonatomic, weak, readonly) id target;
- (instancetype)initWithTarget:(id)target;
+ (instancetype)proxyWithTarget:(id)target;
@end