
#import <Cocoa/Cocoa.h>

NS_ASSUME_NONNULL_BEGIN

@interface NSView (Add)

@property (nonatomic, strong) NSColor *ty_backgroundColor;

@property (nonatomic, assign) CGFloat cornerRadius;

                 
@property (nonatomic, assign) CGFloat width;
                 
@property (nonatomic, assign) CGFloat height;
                 
@property (nonatomic, assign) CGFloat originX;
                 
@property (nonatomic, assign) CGFloat originY;
                 
@property (nonatomic, assign) CGSize size;
                 
@property (nonatomic, assign) CGPoint origin;

   
                                        
                               

                                                                        
   
- (void)setAnchorPoint:(CGPoint)anchor;

@end

NS_ASSUME_NONNULL_END
