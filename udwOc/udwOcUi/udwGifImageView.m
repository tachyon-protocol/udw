#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"
#import <QuartzCore/QuartzCore.h>
#import <ImageIO/ImageIO.h>
#import <MobileCoreServices/MobileCoreServices.h>

#ifndef NS_DESIGNATED_INITIALIZER
#if __has_attribute(objc_designated_initializer)
#define NS_DESIGNATED_INITIALIZER __attribute((objc_designated_initializer))
#else
#define NS_DESIGNATED_INITIALIZER
#endif
#endif

extern const NSTimeInterval kudwPrivateAnimatedImageDelayTimeIntervalMinimum;

  
                                                                                                                                                   
                                                                                                                   
                                                                                                                                                                                                                                                    
                                                                                                                         
  
@interface udwPrivateAnimatedImage : NSObject

@property (nonatomic, strong, readonly) UIImage *posterImage;                                                                                
@property (nonatomic, assign, readonly) CGSize size;                                

@property (nonatomic, assign, readonly) NSUInteger loopCount;                                                
@property (nonatomic, strong, readonly) NSDictionary *delayTimesForIndexes;                                                 
@property (nonatomic, assign, readonly) NSUInteger frameCount;                                                          

@property (nonatomic, assign, readonly) NSUInteger frameCacheSizeCurrent;                                                                                                 
@property (nonatomic, assign) NSUInteger frameCacheSizeMax;                                                                    

                                                                                 
                                                                                                                                      
                                                                                                                       
- (UIImage *)imageLazilyCachedAtIndex:(NSUInteger)index;

                                                                                
+ (CGSize)sizeForImage:(id)image;

                                                                                                                                                          
- (instancetype)initWithAnimatedGIFData:(NSData *)data;
                                                                                         
- (instancetype)initWithAnimatedGIFData:(NSData *)data optimalFrameCacheSize:(NSUInteger)optimalFrameCacheSize predrawingEnabled:(BOOL)isPredrawingEnabled NS_DESIGNATED_INITIALIZER;
+ (instancetype)animatedImageWithGIFData:(NSData *)data;

@property (nonatomic, strong, readonly) NSData *data;                                                         

@end

typedef NS_ENUM(NSUInteger, udwPrivateLogLevel) {
    udwPrivateLogLevelNone = 0,
    udwPrivateLogLevelError,
    udwPrivateLogLevelWarn,
    udwPrivateLogLevelInfo,
    udwPrivateLogLevelDebug,
    udwPrivateLogLevelVerbose
};

@interface udwPrivateAnimatedImage (Logging)

+ (void)setLogBlock:(void (^)(NSString *logString, udwPrivateLogLevel logLevel))logBlock logLevel:(udwPrivateLogLevel)logLevel;
+ (void)logStringFromBlock:(NSString *(^)(void))stringBlock withLevel:(udwPrivateLogLevel)level;

@end

#define udwPrivateLog(logLevel, format, ...) [udwPrivateAnimatedImage logStringFromBlock:^NSString *{ return [NSString stringWithFormat:(format), ## __VA_ARGS__]; } withLevel:(logLevel)]

@interface udwPrivateWeakProxy : NSProxy

+ (instancetype)weakProxyForObject:(id)targetObject;

@end


                                                                    
#ifndef BYTE_SIZE
#define BYTE_SIZE 8                     
#endif

#define MEGABYTE (1024 * 1024)

                                                                                                                                                          
const NSTimeInterval kudwPrivateAnimatedImageDelayTimeIntervalMinimum = 0.02;

                                                                                                                 
                                                                                                                                                    
typedef NS_ENUM(NSUInteger, udwPrivateAnimatedImageDataSizeCategory) {
    udwPrivateAnimatedImageDataSizeCategoryAll = 10,                                                               
    udwPrivateAnimatedImageDataSizeCategoryDefault = 75,                                                                                                            
    udwPrivateAnimatedImageDataSizeCategoryOnDemand = 250,                                                                                     
    udwPrivateAnimatedImageDataSizeCategoryUnsupported                                                       
};

typedef NS_ENUM(NSUInteger, udwPrivateAnimatedImageFrameCacheSize) {
    udwPrivateAnimatedImageFrameCacheSizeNoLimit = 0,                                            
    udwPrivateAnimatedImageFrameCacheSizeLowMemory = 1,                                                                                  
    udwPrivateAnimatedImageFrameCacheSizeGrowAfterMemoryWarning = 2,                                                                                                                        
    udwPrivateAnimatedImageFrameCacheSizeDefault = 5                                                                                
};


#if defined(DEBUG) && DEBUG
@protocol udwPrivateAnimatedImageDebugDelegate <NSObject>
@optional
- (void)debug_animatedImage:(udwPrivateAnimatedImage *)animatedImage didUpdateCachedFrames:(NSIndexSet *)indexesOfFramesInCache;
- (void)debug_animatedImage:(udwPrivateAnimatedImage *)animatedImage didRequestCachedFrame:(NSUInteger)index;
- (CGFloat)debug_animatedImagePredrawingSlowdownFactor:(udwPrivateAnimatedImage *)animatedImage;
@end
#endif


@interface udwPrivateAnimatedImage ()

@property (nonatomic, assign, readonly) NSUInteger frameCacheSizeOptimal;                                                                                               
@property (nonatomic, assign, readonly, getter=isPredrawingEnabled) BOOL predrawingEnabled;                                                        
@property (nonatomic, assign) NSUInteger frameCacheSizeMaxInternal;                                                                                                    
@property (nonatomic, assign) NSUInteger requestedFrameIndex;                                       
@property (nonatomic, assign, readonly) NSUInteger posterImageFrameIndex;                                                     
@property (nonatomic, strong, readonly) NSMutableDictionary *cachedFramesForIndexes;
@property (nonatomic, strong, readonly) NSMutableIndexSet *cachedFrameIndexes;                            
@property (nonatomic, strong, readonly) NSMutableIndexSet *requestedFrameIndexes;                                                                   
@property (nonatomic, strong, readonly) NSIndexSet *allFramesIndexSet;                                                                   
@property (nonatomic, assign) NSUInteger memoryWarningCount;
@property (nonatomic, strong, readonly) dispatch_queue_t serialQueue;
@property (nonatomic, strong, readonly) __attribute__((NSObject)) CGImageSourceRef imageSource;

                                                                                           
                                                                                            
                                                          
@property (nonatomic, strong, readonly) udwPrivateAnimatedImage *weakProxy;

#if defined(DEBUG) && DEBUG
@property (nonatomic, weak) id<udwPrivateAnimatedImageDebugDelegate> debug_delegate;
#endif

@end


                                                                                                                                           
static NSHashTable *allAnimatedImagesWeak;

@implementation udwPrivateAnimatedImage

                        
                   

                                                                      
- (NSUInteger)frameCacheSizeCurrent
{
    NSUInteger frameCacheSizeCurrent = self.frameCacheSizeOptimal;

                                
    if (self.frameCacheSizeMax > udwPrivateAnimatedImageFrameCacheSizeNoLimit) {
        frameCacheSizeCurrent = MIN(frameCacheSizeCurrent, self.frameCacheSizeMax);
    }

    if (self.frameCacheSizeMaxInternal > udwPrivateAnimatedImageFrameCacheSizeNoLimit) {
        frameCacheSizeCurrent = MIN(frameCacheSizeCurrent, self.frameCacheSizeMaxInternal);
    }

    return frameCacheSizeCurrent;
}


- (void)setFrameCacheSizeMax:(NSUInteger)frameCacheSizeMax
{
    if (_frameCacheSizeMax != frameCacheSizeMax) {

                                                                                                                                            
        BOOL willFrameCacheSizeShrink = (frameCacheSizeMax < self.frameCacheSizeCurrent);

                           
        _frameCacheSizeMax = frameCacheSizeMax;

        if (willFrameCacheSizeShrink) {
            [self purgeFrameCacheIfNeeded];
        }
    }
}


                    

- (void)setFrameCacheSizeMaxInternal:(NSUInteger)frameCacheSizeMaxInternal
{
    if (_frameCacheSizeMaxInternal != frameCacheSizeMaxInternal) {

                                                                                                                                            
        BOOL willFrameCacheSizeShrink = (frameCacheSizeMaxInternal < self.frameCacheSizeCurrent);

                           
        _frameCacheSizeMaxInternal = frameCacheSizeMaxInternal;

        if (willFrameCacheSizeShrink) {
            [self purgeFrameCacheIfNeeded];
        }
    }
}


                         

+ (void)initialize
{
    if (self == [udwPrivateAnimatedImage class]) {
                                                                                   
        allAnimatedImagesWeak = [NSHashTable weakObjectsHashTable];

        [[NSNotificationCenter defaultCenter] addObserverForName:UIApplicationDidReceiveMemoryWarningNotification object:nil queue:nil usingBlock:^(NSNotification *note) {
                                                                                                                                                                     
            NSAssert([NSThread isMainThread], @"Received memory warning on non-main thread");
                                                                                                                                                   
                                                                                                                 
            NSArray *images = nil;
            @synchronized(allAnimatedImagesWeak) {
                images = [[allAnimatedImagesWeak allObjects] copy];
            }
                                                                                                    
            [images makeObjectsPerformSelector:@selector(didReceiveMemoryWarning:) withObject:note];
        }];
    }
}


- (instancetype)init
{
    udwPrivateAnimatedImage *animatedImage = [self initWithAnimatedGIFData:nil];
    if (!animatedImage) {
        udwPrivateLog(udwPrivateLogLevelError, @"Use `-initWithAnimatedGIFData:` and supply the animated GIF data as an argument to initialize an object of type `udwPrivateAnimatedImage`.");
    }
    return animatedImage;
}


- (instancetype)initWithAnimatedGIFData:(NSData *)data
{
    return [self initWithAnimatedGIFData:data optimalFrameCacheSize:0 predrawingEnabled:YES];
}

- (instancetype)initWithAnimatedGIFData:(NSData *)data optimalFrameCacheSize:(NSUInteger)optimalFrameCacheSize predrawingEnabled:(BOOL)isPredrawingEnabled
{
                                        
    BOOL hasData = ([data length] > 0);
    if (!hasData) {
        udwPrivateLog(udwPrivateLogLevelError, @"No animated GIF data supplied.");
        return nil;
    }

    self = [super init];
    if (self) {
                                                                                                                                                                       

                                                                              
                                                                                                          
        _data = data;
        _predrawingEnabled = isPredrawingEnabled;

                                              
        _cachedFramesForIndexes = [[NSMutableDictionary alloc] init];
        _cachedFrameIndexes = [[NSMutableIndexSet alloc] init];
        _requestedFrameIndexes = [[NSMutableIndexSet alloc] init];

                                                                                                                                     
        _imageSource = CGImageSourceCreateWithData((__bridge CFDataRef)data,
                                                   (__bridge CFDictionaryRef)@{(NSString *)kCGImageSourceShouldCache: @NO});
                                   
        if (!_imageSource) {
            udwPrivateLog(udwPrivateLogLevelError, @"Failed to `CGImageSourceCreateWithData` for animated GIF data %@", data);
            return nil;
        }

                                   
        CFStringRef imageSourceContainerType = CGImageSourceGetType(_imageSource);
        BOOL isGIFData = UTTypeConformsTo(imageSourceContainerType, kUTTypeGIF);
        if (!isGIFData) {
            udwPrivateLog(udwPrivateLogLevelError, @"Supplied data is of type %@ and doesn't seem to be GIF data %@", imageSourceContainerType, data);
            return nil;
        }

                          
                                                              
                                    
            
                                 
                          
                                         
                                 
                 
            
        NSDictionary *imageProperties = (__bridge_transfer NSDictionary *)CGImageSourceCopyProperties(_imageSource, NULL);
        _loopCount = [[[imageProperties objectForKey:(id)kCGImagePropertyGIFDictionary] objectForKey:(id)kCGImagePropertyGIFLoopCount] unsignedIntegerValue];

                                       
        size_t imageCount = CGImageSourceGetCount(_imageSource);
        NSUInteger skippedFrameCount = 0;
        NSMutableDictionary *delayTimesForIndexesMutable = [NSMutableDictionary dictionaryWithCapacity:imageCount];
        for (size_t i = 0; i < imageCount; i++) {
            @autoreleasepool {
                CGImageRef frameImageRef = CGImageSourceCreateImageAtIndex(_imageSource, i, NULL);
                if (frameImageRef) {
                    UIImage *frameImage = [UIImage imageWithCGImage:frameImageRef];
                                                                                                                                                                          
                    if (frameImage) {
                                           
                        if (!self.posterImage) {
                            _posterImage = frameImage;
                                                              
                            _size = _posterImage.size;
                                                                                                             
                            _posterImageFrameIndex = i;
                            [self.cachedFramesForIndexes setObject:self.posterImage forKey:@(self.posterImageFrameIndex)];
                            [self.cachedFrameIndexes addIndex:self.posterImageFrameIndex];
                        }

                                          
                                                                                                                                                                                                  
                                                    
                            
                                                
                                         
                                                 
                                                
                                          
                                                     
                                                              
                                 
                            

                        NSDictionary *frameProperties = (__bridge_transfer NSDictionary *)CGImageSourceCopyPropertiesAtIndex(_imageSource, i, NULL);
                        NSDictionary *framePropertiesGIF = [frameProperties objectForKey:(id)kCGImagePropertyGIFDictionary];

                                                                                                   
                        NSNumber *delayTime = [framePropertiesGIF objectForKey:(id)kCGImagePropertyGIFUnclampedDelayTime];
                        if (!delayTime) {
                            delayTime = [framePropertiesGIF objectForKey:(id)kCGImagePropertyGIFDelayTime];
                        }
                                                                                                                                                                
                        const NSTimeInterval kDelayTimeIntervalDefault = 0.1;
                        if (!delayTime) {
                            if (i == 0) {
                                udwPrivateLog(udwPrivateLogLevelInfo, @"Falling back to default delay time for first frame %@ because none found in GIF properties %@", frameImage, frameProperties);
                                delayTime = @(kDelayTimeIntervalDefault);
                            } else {
                                udwPrivateLog(udwPrivateLogLevelInfo, @"Falling back to preceding delay time for frame %zu %@ because none found in GIF properties %@", i, frameImage, frameProperties);
                                delayTime = delayTimesForIndexesMutable[@(i - 1)];
                            }
                        }
                                                                                                                                                                                                           
                                                                                                                                                                                             
                        if ([delayTime floatValue] < ((float)kudwPrivateAnimatedImageDelayTimeIntervalMinimum - FLT_EPSILON)) {
                            udwPrivateLog(udwPrivateLogLevelInfo, @"Rounding frame %zu's `delayTime` from %f up to default %f (minimum supported: %f).", i, [delayTime floatValue], kDelayTimeIntervalDefault, kudwPrivateAnimatedImageDelayTimeIntervalMinimum);
                            delayTime = @(kDelayTimeIntervalDefault);
                        }
                        delayTimesForIndexesMutable[@(i)] = delayTime;
                    } else {
                        skippedFrameCount++;
                        udwPrivateLog(udwPrivateLogLevelInfo, @"Dropping frame %zu because valid `CGImageRef` %@ did result in `nil`-`UIImage`.", i, frameImageRef);
                    }
                    CFRelease(frameImageRef);
                } else {
                    skippedFrameCount++;
                    udwPrivateLog(udwPrivateLogLevelInfo, @"Dropping frame %zu because failed to `CGImageSourceCreateImageAtIndex` with image source %@", i, _imageSource);
                }
            }
        }
        _delayTimesForIndexes = [delayTimesForIndexesMutable copy];
        _frameCount = imageCount;

        if (self.frameCount == 0) {
            udwPrivateLog(udwPrivateLogLevelInfo, @"Failed to create any valid frames for GIF with properties %@", imageProperties);
            return nil;
        } else if (self.frameCount == 1) {
                                                                            
            udwPrivateLog(udwPrivateLogLevelInfo, @"Created valid GIF but with only a single frame. Image properties: %@", imageProperties);
        } else {
                                                
        }

                                                                      
        if (optimalFrameCacheSize == 0) {
                                                                                                                                 
                                                                                          
            CGFloat animatedImageDataSize = CGImageGetBytesPerRow(self.posterImage.CGImage) * self.size.height * (self.frameCount - skippedFrameCount) / MEGABYTE;
            if (animatedImageDataSize <= udwPrivateAnimatedImageDataSizeCategoryAll) {
                _frameCacheSizeOptimal = self.frameCount;
            } else if (animatedImageDataSize <= udwPrivateAnimatedImageDataSizeCategoryDefault) {
                                                                                                                                                                                                                                                                                                               
                _frameCacheSizeOptimal = udwPrivateAnimatedImageFrameCacheSizeDefault;
            } else {
                                                                                                                               
                _frameCacheSizeOptimal = udwPrivateAnimatedImageFrameCacheSizeLowMemory;
            }
        } else {
                                      
            _frameCacheSizeOptimal = optimalFrameCacheSize;
        }
                                                                      
        _frameCacheSizeOptimal = MIN(_frameCacheSizeOptimal, self.frameCount);

                                                                                                                                       
        _allFramesIndexSet = [[NSIndexSet alloc] initWithIndexesInRange:NSMakeRange(0, self.frameCount)];

                                                          
        _weakProxy = (id)[udwPrivateWeakProxy weakProxyForObject:self];

                                                                                                                                         
                                                                                                             
        @synchronized(allAnimatedImagesWeak) {
            [allAnimatedImagesWeak addObject:self];
        }
    }
    return self;
}


+ (instancetype)animatedImageWithGIFData:(NSData *)data
{
    udwPrivateAnimatedImage *animatedImage = [[udwPrivateAnimatedImage alloc] initWithAnimatedGIFData:data];
    return animatedImage;
}


- (void)dealloc
{
    if (_weakProxy) {
        [NSObject cancelPreviousPerformRequestsWithTarget:_weakProxy];
    }

    if (_imageSource) {
        CFRelease(_imageSource);
    }
}


                             

                               
                                                                                                                                           
- (UIImage *)imageLazilyCachedAtIndex:(NSUInteger)index
{
                                                            
                                                                                                
    if (index >= self.frameCount) {
        udwPrivateLog(udwPrivateLogLevelWarn, @"Skipping requested frame %lu beyond bounds (total frame count: %lu) for animated image: %@", (unsigned long)index, (unsigned long)self.frameCount, self);
        return nil;
    }

                                                                                 
    self.requestedFrameIndex = index;
#if defined(DEBUG) && DEBUG
    if ([self.debug_delegate respondsToSelector:@selector(debug_animatedImage:didRequestCachedFrame:)]) {
        [self.debug_delegate debug_animatedImage:self didRequestCachedFrame:index];
    }
#endif

                                                                                                        
    if ([self.cachedFrameIndexes count] < self.frameCount) {
                                                                                                     
                                                                                                       
        NSMutableIndexSet *frameIndexesToAddToCacheMutable = [self frameIndexesToCache];
        [frameIndexesToAddToCacheMutable removeIndexes:self.cachedFrameIndexes];
        [frameIndexesToAddToCacheMutable removeIndexes:self.requestedFrameIndexes];
        [frameIndexesToAddToCacheMutable removeIndex:self.posterImageFrameIndex];
        NSIndexSet *frameIndexesToAddToCache = [frameIndexesToAddToCacheMutable copy];

                                                  
        if ([frameIndexesToAddToCache count] > 0) {
            [self addFrameIndexesToCache:frameIndexesToAddToCache];
        }
    }

                               
    UIImage *image = self.cachedFramesForIndexes[@(index)];

                                                              
    [self purgeFrameCacheIfNeeded];

    return image;
}


                                                                                                           
- (void)addFrameIndexesToCache:(NSIndexSet *)frameIndexesToAddToCache
{
                                                                                              
                                                                                 
    NSRange firstRange = NSMakeRange(self.requestedFrameIndex, self.frameCount - self.requestedFrameIndex);
    NSRange secondRange = NSMakeRange(0, self.requestedFrameIndex);
    if (firstRange.length + secondRange.length != self.frameCount) {
        udwPrivateLog(udwPrivateLogLevelWarn, @"Two-part frame cache range doesn't equal full range.");
    }

                                                                                                          
    [self.requestedFrameIndexes addIndexes:frameIndexesToAddToCache];

                                               
    if (!self.serialQueue) {
        _serialQueue = dispatch_queue_create("com.flipboard.framecachingqueue", DISPATCH_QUEUE_SERIAL);
    }

                                                                         
                                                                                                                 
    udwPrivateAnimatedImage * __weak weakSelf = self;
    dispatch_async(self.serialQueue, ^{
                                               
        void (^frameRangeBlock)(NSRange, BOOL *) = ^(NSRange range, BOOL *stop) {
                                                                                                                    
            for (NSUInteger i = range.location; i < NSMaxRange(range); i++) {
#if defined(DEBUG) && DEBUG
                CFTimeInterval predrawBeginTime = CACurrentMediaTime();
#endif
                UIImage *image = [weakSelf imageAtIndex:i];
#if defined(DEBUG) && DEBUG
                CFTimeInterval predrawDuration = CACurrentMediaTime() - predrawBeginTime;
                CFTimeInterval slowdownDuration = 0.0;
                if ([self.debug_delegate respondsToSelector:@selector(debug_animatedImagePredrawingSlowdownFactor:)]) {
                    CGFloat predrawingSlowdownFactor = [self.debug_delegate debug_animatedImagePredrawingSlowdownFactor:self];
                    slowdownDuration = predrawDuration * predrawingSlowdownFactor - predrawDuration;
                    [NSThread sleepForTimeInterval:slowdownDuration];
                }
                udwPrivateLog(udwPrivateLogLevelVerbose, @"Predrew frame %lu in %f ms for animated image: %@", (unsigned long)i, (predrawDuration + slowdownDuration) * 1000, self);
#endif
                                                                                                   
                                                                                                                                                                             
                if (image && weakSelf) {
                    dispatch_async(dispatch_get_main_queue(), ^{
                        weakSelf.cachedFramesForIndexes[@(i)] = image;
                        [weakSelf.cachedFrameIndexes addIndex:i];
                        [weakSelf.requestedFrameIndexes removeIndex:i];
#if defined(DEBUG) && DEBUG
                        if ([weakSelf.debug_delegate respondsToSelector:@selector(debug_animatedImage:didUpdateCachedFrames:)]) {
                            [weakSelf.debug_delegate debug_animatedImage:weakSelf didUpdateCachedFrames:weakSelf.cachedFrameIndexes];
                        }
#endif
                    });
                }
            }
        };

        [frameIndexesToAddToCache enumerateRangesInRange:firstRange options:0 usingBlock:frameRangeBlock];
        [frameIndexesToAddToCache enumerateRangesInRange:secondRange options:0 usingBlock:frameRangeBlock];
    });
}


+ (CGSize)sizeForImage:(id)image
{
    CGSize imageSize = CGSizeZero;

                           
    if (!image) {
        return imageSize;
    }

    if ([image isKindOfClass:[UIImage class]]) {
        UIImage *uiImage = (UIImage *)image;
        imageSize = uiImage.size;
    } else if ([image isKindOfClass:[udwPrivateAnimatedImage class]]) {
        udwPrivateAnimatedImage *animatedImage = (udwPrivateAnimatedImage *)image;
        imageSize = animatedImage.size;
    } else {
                                                                                       
        udwPrivateLog(udwPrivateLogLevelError, @"`image` isn't of expected types `UIImage` or `udwPrivateAnimatedImage`: %@", image);
    }

    return imageSize;
}


                              
                          

- (UIImage *)imageAtIndex:(NSUInteger)index
{
                                                                                                                                                                                                                                 
    CGImageRef imageRef = CGImageSourceCreateImageAtIndex(_imageSource, index, NULL);

                           
    if (!imageRef) {
        return nil;
    }

    UIImage *image = [UIImage imageWithCGImage:imageRef];
    CFRelease(imageRef);

                                                                                                                                                                                                         
    if (self.isPredrawingEnabled) {
        image = [[self class] predrawnImageFromImage:image];
    }

    return image;
}


                          

- (NSMutableIndexSet *)frameIndexesToCache
{
    NSMutableIndexSet *indexesToCache = nil;
                                                                                                                 
    if (self.frameCacheSizeCurrent == self.frameCount) {
        indexesToCache = [self.allFramesIndexSet mutableCopy];
    } else {
        indexesToCache = [[NSMutableIndexSet alloc] init];

                                                                                                                                        
                                                                                         
        NSUInteger firstLength = MIN(self.frameCacheSizeCurrent, self.frameCount - self.requestedFrameIndex);
        NSRange firstRange = NSMakeRange(self.requestedFrameIndex, firstLength);
        [indexesToCache addIndexesInRange:firstRange];
        NSUInteger secondLength = self.frameCacheSizeCurrent - firstLength;
        if (secondLength > 0) {
            NSRange secondRange = NSMakeRange(0, secondLength);
            [indexesToCache addIndexesInRange:secondRange];
        }
                                                                                                    
        if ([indexesToCache count] != self.frameCacheSizeCurrent) {
            udwPrivateLog(udwPrivateLogLevelWarn, @"Number of frames to cache doesn't equal expected cache size.");
        }

        [indexesToCache addIndex:self.posterImageFrameIndex];
    }

    return indexesToCache;
}


- (void)purgeFrameCacheIfNeeded
{
                                                                   
                                                                  
                                                                                                                       
    if ([self.cachedFrameIndexes count] > self.frameCacheSizeCurrent) {
        NSMutableIndexSet *indexesToPurge = [self.cachedFrameIndexes mutableCopy];
        [indexesToPurge removeIndexes:[self frameIndexesToCache]];
        [indexesToPurge enumerateRangesUsingBlock:^(NSRange range, BOOL *stop) {
                                                                                                                    
            for (NSUInteger i = range.location; i < NSMaxRange(range); i++) {
                [self.cachedFrameIndexes removeIndex:i];
                [self.cachedFramesForIndexes removeObjectForKey:@(i)];
                                                                                                                                                                  
#if defined(DEBUG) && DEBUG
                if ([self.debug_delegate respondsToSelector:@selector(debug_animatedImage:didUpdateCachedFrames:)]) {
                    dispatch_async(dispatch_get_main_queue(), ^{
                        [self.debug_delegate debug_animatedImage:self didUpdateCachedFrames:self.cachedFrameIndexes];
                    });
                }
#endif
            }
        }];
    }
}


- (void)growFrameCacheSizeAfterMemoryWarning:(NSNumber *)frameCacheSize
{
    self.frameCacheSizeMaxInternal = [frameCacheSize unsignedIntegerValue];
    udwPrivateLog(udwPrivateLogLevelDebug, @"Grew frame cache size max to %lu after memory warning for animated image: %@", (unsigned long)self.frameCacheSizeMaxInternal, self);

                                                                            
    const NSTimeInterval kResetDelay = 3.0;
    [self.weakProxy performSelector:@selector(resetFrameCacheSizeMaxInternal) withObject:nil afterDelay:kResetDelay];
}


- (void)resetFrameCacheSizeMaxInternal
{
    self.frameCacheSizeMaxInternal = udwPrivateAnimatedImageFrameCacheSizeNoLimit;
    udwPrivateLog(udwPrivateLogLevelDebug, @"Reset frame cache size max (current frame cache size: %lu) for animated image: %@", (unsigned long)self.frameCacheSizeCurrent, self);
}


                                                        

- (void)didReceiveMemoryWarning:(NSNotification *)notification
{
    self.memoryWarningCount++;

                                                                                                   
    [NSObject cancelPreviousPerformRequestsWithTarget:self.weakProxy selector:@selector(growFrameCacheSizeAfterMemoryWarning:) object:@(udwPrivateAnimatedImageFrameCacheSizeGrowAfterMemoryWarning)];
    [NSObject cancelPreviousPerformRequestsWithTarget:self.weakProxy selector:@selector(resetFrameCacheSizeMaxInternal) object:nil];

                                                                                                                                                                         
    udwPrivateLog(udwPrivateLogLevelDebug, @"Attempt setting frame cache size max to %lu (previous was %lu) after memory warning #%lu for animated image: %@", (unsigned long)udwPrivateAnimatedImageFrameCacheSizeLowMemory, (unsigned long)self.frameCacheSizeMaxInternal, (unsigned long)self.memoryWarningCount, self);
    self.frameCacheSizeMaxInternal = udwPrivateAnimatedImageFrameCacheSizeLowMemory;

                                                                                                                                                                   
      
                                                                            
                                                                      
                                                                                                        
                                                                                                                         
                                                                      
                                                                                                       
                                                                                           
                                     
                                                                                                                     
      
    const NSUInteger kGrowAttemptsMax = 2;
    const NSTimeInterval kGrowDelay = 2.0;
    if ((self.memoryWarningCount - 1) <= kGrowAttemptsMax) {
        [self.weakProxy performSelector:@selector(growFrameCacheSizeAfterMemoryWarning:) withObject:@(udwPrivateAnimatedImageFrameCacheSizeGrowAfterMemoryWarning) afterDelay:kGrowDelay];
    }

                                                                                                                                                                                                             
}


                           

                                                                                                                                     
                                                                                                          
                                                                                                                                              
                                                                              
+ (UIImage *)predrawnImageFromImage:(UIImage *)imageToPredraw
{
                                                                                                   
    CGColorSpaceRef colorSpaceDeviceRGBRef = CGColorSpaceCreateDeviceRGB();
                               
    if (!colorSpaceDeviceRGBRef) {
        udwPrivateLog(udwPrivateLogLevelError, @"Failed to `CGColorSpaceCreateDeviceRGB` for image %@", imageToPredraw);
        return imageToPredraw;
    }

                                                                                                                                                                    
                                                                                                                             
                                                                                                                                        
    size_t numberOfComponents = CGColorSpaceGetNumberOfComponents(colorSpaceDeviceRGBRef) + 1;              

                                                                                                                                               
    void *data = NULL;
    size_t width = imageToPredraw.size.width;
    size_t height = imageToPredraw.size.height;
    size_t bitsPerComponent = CHAR_BIT;

    size_t bitsPerPixel = (bitsPerComponent * numberOfComponents);
    size_t bytesPerPixel = (bitsPerPixel / BYTE_SIZE);
    size_t bytesPerRow = (bytesPerPixel * width);

    CGBitmapInfo bitmapInfo = kCGBitmapByteOrderDefault;

    CGImageAlphaInfo alphaInfo = CGImageGetAlphaInfo(imageToPredraw.CGImage);
                                                                                                                    
                                                                                                                                                      
    if (alphaInfo == kCGImageAlphaNone || alphaInfo == kCGImageAlphaOnly) {
        alphaInfo = kCGImageAlphaNoneSkipFirst;
    } else if (alphaInfo == kCGImageAlphaFirst) {
        alphaInfo = kCGImageAlphaPremultipliedFirst;
    } else if (alphaInfo == kCGImageAlphaLast) {
        alphaInfo = kCGImageAlphaPremultipliedLast;
    }
                                                                                                                                                                            
    bitmapInfo |= alphaInfo;

                                                                                                                                                                                                                                                            
                                                                                                                                                                                                                                              
    CGContextRef bitmapContextRef = CGBitmapContextCreate(data, width, height, bitsPerComponent, bytesPerRow, colorSpaceDeviceRGBRef, bitmapInfo);
    CGColorSpaceRelease(colorSpaceDeviceRGBRef);
                               
    if (!bitmapContextRef) {
        udwPrivateLog(udwPrivateLogLevelError, @"Failed to `CGBitmapContextCreate` with color space %@ and parameters (width: %zu height: %zu bitsPerComponent: %zu bytesPerRow: %zu) for image %@", colorSpaceDeviceRGBRef, width, height, bitsPerComponent, bytesPerRow, imageToPredraw);
        return imageToPredraw;
    }

                                                                                         
    CGContextDrawImage(bitmapContextRef, CGRectMake(0.0, 0.0, imageToPredraw.size.width, imageToPredraw.size.height), imageToPredraw.CGImage);
    CGImageRef predrawnImageRef = CGBitmapContextCreateImage(bitmapContextRef);
    UIImage *predrawnImage = [UIImage imageWithCGImage:predrawnImageRef scale:imageToPredraw.scale orientation:imageToPredraw.imageOrientation];
    CGImageRelease(predrawnImageRef);
    CGContextRelease(bitmapContextRef);

                               
    if (!predrawnImage) {
        udwPrivateLog(udwPrivateLogLevelError, @"Failed to `imageWithCGImage:scale:orientation:` with image ref %@ created with color space %@ and bitmap context %@ and properties and properties (scale: %f orientation: %ld) for image %@", predrawnImageRef, colorSpaceDeviceRGBRef, bitmapContextRef, imageToPredraw.scale, (long)imageToPredraw.imageOrientation, imageToPredraw);
        return imageToPredraw;
    }

    return predrawnImage;
}


                          

- (NSString *)description
{
    NSString *description = [super description];

    description = [description stringByAppendingFormat:@" size=%@", NSStringFromCGSize(self.size)];
    description = [description stringByAppendingFormat:@" frameCount=%lu", (unsigned long)self.frameCount];

    return description;
}


@end

                      

@implementation udwPrivateAnimatedImage (Logging)

static void (^_logBlock)(NSString *logString, udwPrivateLogLevel logLevel) = nil;
static udwPrivateLogLevel _logLevel;

+ (void)setLogBlock:(void (^)(NSString *logString, udwPrivateLogLevel logLevel))logBlock logLevel:(udwPrivateLogLevel)logLevel
{
    _logBlock = logBlock;
    _logLevel = logLevel;
}

+ (void)logStringFromBlock:(NSString *(^)(void))stringBlock withLevel:(udwPrivateLogLevel)level
{
    if (level <= _logLevel && _logBlock && stringBlock) {
        _logBlock(stringBlock(), level);
    }
}

@end


                                  

@interface udwPrivateWeakProxy ()

@property (nonatomic, weak) id target;

@end


@implementation udwPrivateWeakProxy

                       

                                                                         
                                                                    
+ (instancetype)weakProxyForObject:(id)targetObject
{
    udwPrivateWeakProxy *weakProxy = [udwPrivateWeakProxy alloc];
    weakProxy.target = targetObject;
    return weakProxy;
}


                                

- (id)forwardingTargetForSelector:(SEL)selector
{
                                                    
    return _target;
}


                                           
                                           

- (void)forwardInvocation:(NSInvocation *)invocation
{
                                                                                  
                                                                                                                       
                                                                                        
    void *nullPointer = NULL;
    [invocation setReturnValue:&nullPointer];
}


- (NSMethodSignature *)methodSignatureForSelector:(SEL)selector
{
                                                                      
                                                                                                                                         
                                                                                                                                                                             
                                                                                                                                                            
                                                                                                                                                                                     
    return [NSObject instanceMethodSignatureForSelector:@selector(init)];
}


@end

#if defined(DEBUG) && DEBUG
@protocol udwPrivateAnimatedImageViewDebugDelegate <NSObject>
@optional
- (void)debug_animatedImageView:(udwGifImageView *)animatedImageView waitingForFrame:(NSUInteger)index duration:(NSTimeInterval)duration;
@end
#endif


@interface udwGifImageView ()

                                                                  
@property (nonatomic, strong) udwPrivateAnimatedImage *animatedImage;
@property (nonatomic, strong, readwrite) UIImage *currentFrame;
@property (nonatomic, assign, readwrite) NSUInteger currentFrameIndex;

@property (nonatomic, assign) NSUInteger loopCountdown;
@property (nonatomic, assign) NSTimeInterval accumulator;
@property (nonatomic, strong) CADisplayLink *displayLink;

@property (nonatomic, assign) BOOL shouldAnimate;                                                                                                                                                     
@property (nonatomic, assign) BOOL needsDisplayWhenImageBecomesAvailable;

#if defined(DEBUG) && DEBUG
@property (nonatomic, weak) id<udwPrivateAnimatedImageViewDebugDelegate> debug_delegate;
#endif

@end


@implementation udwGifImageView
@synthesize runLoopMode = _runLoopMode;

                           

                                                                                                            
                                                                               
- (instancetype)initWithImage:(UIImage *)image
{
    self = [super initWithImage:image];
    if (self) {
        [self commonInit];
    }
    return self;
}

                                                                                                                                                            
- (instancetype)initWithImage:(UIImage *)image highlightedImage:(UIImage *)highlightedImage
{
    self = [super initWithImage:image highlightedImage:highlightedImage];
    if (self) {
        [self commonInit];
    }
    return self;
}

- (instancetype)initWithFrame:(CGRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        [self commonInit];
    }
    return self;
}

- (instancetype)initWithFrame:(CGRect)frame gifName:(NSString *)gifName
{
    self = [super initWithFrame:frame];
    if (self) {
        [self commonInit];
        NSString *GIFName= gifName;
        NSString *path = [[NSBundle mainBundle] pathForResource:GIFName ofType:@"gif"];
        NSData *data = [NSData dataWithContentsOfFile:path];
        udwPrivateAnimatedImage *image = [udwPrivateAnimatedImage animatedImageWithGIFData:data];
        [self setAnimatedImage:image];
    }
    return self;
}

- (instancetype)initWithCoder:(NSCoder *)aDecoder
{
    self = [super initWithCoder:aDecoder];
    if (self) {
        [self commonInit];
    }
    return self;
}

- (void)commonInit
{
    self.runLoopMode = [[self class] defaultRunLoopMode];
}


                        
                   

- (void)setAnimatedImage:(udwPrivateAnimatedImage *)animatedImage
{
    if (![_animatedImage isEqual:animatedImage]) {
        if (animatedImage) {
                                   
            super.image = nil;
                                                                                         
            super.highlighted = NO;
                                                                                                                                                                                
            [self invalidateIntrinsicContentSize];
        } else {
                                                                         
            [self stopAnimating];
        }

        _animatedImage = animatedImage;

        self.currentFrame = animatedImage.posterImage;
        self.currentFrameIndex = 0;
        if (animatedImage.loopCount > 0) {
            self.loopCountdown = animatedImage.loopCount;
        } else {
            self.loopCountdown = NSUIntegerMax;
        }
        self.accumulator = 0.0;

                                                                     
        [self updateShouldAnimate];
        if (self.shouldAnimate) {
            [self startAnimating];
        }

        [self.layer setNeedsDisplay];
    }
}


                         

- (void)dealloc
{
                                                        
    [_displayLink invalidate];
}


                                      
                                           

- (void)didMoveToSuperview
{
    [super didMoveToSuperview];

    [self updateShouldAnimate];
    if (self.shouldAnimate) {
        [self startAnimating];
    } else {
        [self stopAnimating];
    }
}


- (void)didMoveToWindow
{
    [super didMoveToWindow];

    [self updateShouldAnimate];
    if (self.shouldAnimate) {
        [self startAnimating];
    } else {
        [self stopAnimating];
    }
}

- (void)setAlpha:(CGFloat)alpha
{
    [super setAlpha:alpha];

    [self updateShouldAnimate];
    if (self.shouldAnimate) {
        [self startAnimating];
    } else {
        [self stopAnimating];
    }
}

- (void)setHidden:(BOOL)hidden
{
    [super setHidden:hidden];

    [self updateShouldAnimate];
    if (self.shouldAnimate) {
        [self startAnimating];
    } else {
        [self stopAnimating];
    }
}


                        

- (CGSize)intrinsicContentSize
{
                                                                                                      
    CGSize intrinsicContentSize = [super intrinsicContentSize];

                                                             
                                                                                                                                                                                                                                                                                                                     
                                                                                                                                                                                      
    if (self.animatedImage) {
        intrinsicContentSize = self.image.size;
    }

    return intrinsicContentSize;
}


                                           
                       

- (UIImage *)image
{
    UIImage *image = nil;
    if (self.animatedImage) {
                                             
        image = self.currentFrame;
    } else {
        image = super.image;
    }
    return image;
}


- (void)setImage:(UIImage *)image
{
    if (image) {
                                                                                
        self.animatedImage = nil;
    }

    super.image = image;
}


                             

- (NSTimeInterval)frameDelayGreatestCommonDivisor
{
                                                                                                                              
    const NSTimeInterval kGreatestCommonDivisorPrecision = 2.0 / kudwPrivateAnimatedImageDelayTimeIntervalMinimum;

    NSArray *delays = self.animatedImage.delayTimesForIndexes.allValues;

                                                                   
                                                                         
    NSUInteger scaledGCD = lrint([delays.firstObject floatValue] * kGreatestCommonDivisorPrecision);
    for (NSNumber *value in delays) {
        scaledGCD = gcd(lrint([value floatValue] * kGreatestCommonDivisorPrecision), scaledGCD);
    }

                                                           
    return scaledGCD / kGreatestCommonDivisorPrecision;
}


static NSUInteger gcd(NSUInteger a, NSUInteger b)
{
                                                           
    if (a < b) {
        return gcd(b, a);
    } else if (a == b) {
        return b;
    }

    while (true) {
        NSUInteger remainder = a % b;
        if (remainder == 0) {
            return b;
        }
        a = b;
        b = remainder;
    }
}


- (void)startAnimating
{
    if (self.animatedImage) {
                                          
        if (!self.displayLink) {
                                                                                                                               
                                                                                                                              
                                                                                                                  
                                                                                                    
            udwPrivateWeakProxy *weakProxy = [udwPrivateWeakProxy weakProxyForObject:self];
            self.displayLink = [CADisplayLink displayLinkWithTarget:weakProxy selector:@selector(displayDidRefresh:)];

            [self.displayLink addToRunLoop:[NSRunLoop mainRunLoop] forMode:self.runLoopMode];
        }

                                                                                                                                             
                                                                                                           
        const NSTimeInterval kDisplayRefreshRate = 60.0;        
        self.displayLink.frameInterval = MAX([self frameDelayGreatestCommonDivisor] * kDisplayRefreshRate, 1);

        self.displayLink.paused = NO;
    } else {
        [super startAnimating];
    }
}

- (void)setRunLoopMode:(NSString *)runLoopMode
{
    if (![@[NSDefaultRunLoopMode, NSRunLoopCommonModes] containsObject:runLoopMode]) {
        NSAssert(NO, @"Invalid run loop mode: %@", runLoopMode);
        _runLoopMode = [[self class] defaultRunLoopMode];
    } else {
        _runLoopMode = runLoopMode;
    }
}

- (void)stopAnimating
{
    if (self.animatedImage) {
        self.displayLink.paused = YES;
    } else {
        [super stopAnimating];
    }
}


- (BOOL)isAnimating
{
    BOOL isAnimating = NO;
    if (self.animatedImage) {
        isAnimating = self.displayLink && !self.displayLink.isPaused;
    } else {
        isAnimating = [super isAnimating];
    }
    return isAnimating;
}


                                        

- (void)setHighlighted:(BOOL)highlighted
{
                                                                                                                                               
    if (!self.animatedImage) {
        [super setHighlighted:highlighted];
    }
}


                              
                      

                                                                                                  
                                                                                                                        
- (void)updateShouldAnimate
{
    BOOL isVisible = self.window && self.superview && ![self isHidden] && self.alpha > 0.0;
    self.shouldAnimate = self.animatedImage && isVisible;
}


- (void)displayDidRefresh:(CADisplayLink *)displayLink
{
                                                                                            
                    
    if (!self.shouldAnimate) {
        udwPrivateLog(udwPrivateLogLevelWarn, @"Trying to animate image when we shouldn't: %@", self);
        return;
    }

    NSNumber *delayTimeNumber = [self.animatedImage.delayTimesForIndexes objectForKey:@(self.currentFrameIndex)];
                                                                                                                                          
    if (delayTimeNumber) {
        NSTimeInterval delayTime = [delayTimeNumber floatValue];
                                                                                               
        UIImage *image = [self.animatedImage imageLazilyCachedAtIndex:self.currentFrameIndex];
        if (image) {
            udwPrivateLog(udwPrivateLogLevelVerbose, @"Showing frame %lu for animated image: %@", (unsigned long)self.currentFrameIndex, self.animatedImage);
            self.currentFrame = image;
            if (self.needsDisplayWhenImageBecomesAvailable) {
                [self.layer setNeedsDisplay];
                self.needsDisplayWhenImageBecomesAvailable = NO;
            }

            self.accumulator += displayLink.duration * displayLink.frameInterval;

                                                                                                                              
            while (self.accumulator >= delayTime) {
                self.accumulator -= delayTime;
                self.currentFrameIndex++;
                if (self.currentFrameIndex >= self.animatedImage.frameCount) {
                                                                                                            
                    self.loopCountdown--;
                    if (self.loopCompletionBlock) {
                        self.loopCompletionBlock(self.loopCountdown);
                    }

                    if (self.loopCountdown == 0) {
                        [self stopAnimating];
                        return;
                    }
                    self.currentFrameIndex = 0;
                }
                                                                                                                             
                                                                                                                              
                self.needsDisplayWhenImageBecomesAvailable = YES;
            }
        } else {
            udwPrivateLog(udwPrivateLogLevelDebug, @"Waiting for frame %lu for animated image: %@", (unsigned long)self.currentFrameIndex, self.animatedImage);
#if defined(DEBUG) && DEBUG
            if ([self.debug_delegate respondsToSelector:@selector(debug_animatedImageView:waitingForFrame:duration:)]) {
                [self.debug_delegate debug_animatedImageView:self waitingForFrame:self.currentFrameIndex duration:(NSTimeInterval)displayLink.duration * displayLink.frameInterval];
            }
#endif
        }
    } else {
        self.currentFrameIndex++;
    }
}

+ (NSString *)defaultRunLoopMode
{
                                                                                                                                    
    return [NSProcessInfo processInfo].activeProcessorCount > 1 ? NSRunLoopCommonModes : NSDefaultRunLoopMode;
}


                                         
                                          

- (void)displayLayer:(CALayer *)layer
{
    layer.contents = (__bridge id)self.image.CGImage;
}


@end
