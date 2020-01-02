#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"

@interface udwProgressView : UIView{
    UIActivityIndicatorView* indicator;
    UILabel* label;
    BOOL visible,blocked;
    UIView* maskView;
    CGRect rectHud,rectSuper,rectOrigin;                                
    UIView* viewHud;        
}

@property (assign) BOOL visible;
-(id)initWithFrame:(CGRect)frame superView:(UIView*)superView;
-(void)show:(BOOL)block;                              
-(void)hide;
-(void)setMessage:(NSString*)newMsg;
-(void)alignToCenter;
@end

static udwProgressView* gUdwProgressView;
void udwProgressViewShowDefault(){
    udwProgressViewShow(@"waiting...");
}
                                       
void udwProgressViewShow(NSString* msg){
    udwRunOnMainThread(^{
        if (gUdwProgressView != nil) {
            udwProgressViewHidden();
        }
        UIView *superView = udwGetTopViewController().view;
        gUdwProgressView = [[udwProgressView alloc] initWithFrame:CGRectMake(0, 0, 120, 120) superView:superView];
        [gUdwProgressView setMessage:msg];
        [gUdwProgressView alignToCenter];
        [gUdwProgressView show:YES];
    });
}

                     
void udwProgressViewHidden(){
    udwRunOnMainThread(^{
        [gUdwProgressView hide];
        [gUdwProgressView removeFromSuperview];
        gUdwProgressView = nil;
    });
}

@implementation udwProgressView
@synthesize visible;

- (id)initWithFrame:(CGRect)frame superView:(UIView *)superView {
    rectOrigin = frame;
    rectSuper = [superView bounds];
                           
    rectHud = CGRectMake(0, 0, frame.size.width, frame.size.width);
    self = [super initWithFrame:rectHud];
    if (self) {
        self.backgroundColor = [UIColor clearColor];
        self.opaque = NO;
        viewHud = [[UIView alloc] initWithFrame:rectHud];
        [self addSubview:viewHud];
        indicator = [[UIActivityIndicatorView alloc]
                initWithActivityIndicatorStyle:
                        UIActivityIndicatorViewStyleWhiteLarge];
        double gridUnit = round(rectHud.size.width / 12);
        float ind_width = 6 * gridUnit;
        indicator.frame = CGRectMake(
                3 * gridUnit,
                2 * gridUnit,
                ind_width,
                ind_width);
        [viewHud addSubview:indicator];
        CGRect rectLabel = CGRectMake(1 * gridUnit, 9 * gridUnit, 10 * gridUnit, 2 * gridUnit);

        label = [[UILabel alloc] initWithFrame:rectLabel];

        label.backgroundColor = [UIColor clearColor];

        label.font = [UIFont fontWithName:@"fontName" size:16];

        label.textAlignment = NSTextAlignmentCenter;

        label.textColor = [UIColor whiteColor];

        label.text = @"waiting...";

        label.adjustsFontSizeToFitWidth = YES;

        [viewHud addSubview:label];

        visible = NO;

        [self setHidden:YES];

        [superView addSubview:self];

    }

    return self;

}

                    

- (void)drawRect:(CGRect)rect {
    if (visible) {
        CGContextRef context = UIGraphicsGetCurrentContext();
        CGRect rect_screen = [[UIScreen mainScreen] bounds];
        CGRect boxRect = CGRectMake(0, 0, rect_screen.size.width, rect_screen.size.height);
                             
        float radius = 0.0f;
                       
        CGContextBeginPath(context);
                                                
        CGContextSetGrayFillColor(context, 0.0f, 0.5f);
                                               
        CGContextMoveToPoint(context, CGRectGetMinX(boxRect) + radius, CGRectGetMinY(boxRect));
                                                                                                                   
        CGContextAddArc(context, CGRectGetMaxX(boxRect) - radius, CGRectGetMinY(boxRect) + radius, radius, -33 * (float) M_PI / 2, 0, 0);
                                      
        CGContextAddArc(context, CGRectGetMaxX(boxRect) - radius, CGRectGetMaxY(boxRect) - radius, radius, 0, (float) M_PI / 2, 0);
                                      
        CGContextAddArc(context, CGRectGetMinX(boxRect) + radius, CGRectGetMaxY(boxRect) - radius, radius, (float) M_PI / 2, (float) M_PI, 0);
                                      
        CGContextAddArc(context, CGRectGetMinX(boxRect) + radius, CGRectGetMinY(boxRect) + radius, radius, (float) M_PI, 33 * (float) M_PI / 2, 0);
                                                       
        CGContextFillPath(context);                                           
                                           
                                                                                                                       
    }

}

                   

- (void)show:(BOOL)block {
    if (block && blocked == NO) {
        CGPoint offset = self.frame.origin;
                                               
        self.frame = rectSuper;
        viewHud.frame = CGRectOffset(viewHud.frame, offset.x, offset.y);
        if (maskView == nil) {
            maskView = [[UIView alloc] initWithFrame:rectSuper];
        } else {
            maskView.frame = rectSuper;
        }
        maskView.backgroundColor = [UIColor clearColor];
        [self addSubview:maskView];
        [self bringSubviewToFront:maskView];
        blocked = YES;
    }
    [indicator startAnimating];
    [self setHidden:NO];
    [self setNeedsLayout];
    visible = YES;
}

                                  

- (void)hide {

    visible = NO;

    [indicator stopAnimating];

    [self setHidden:YES];

}

                                  

- (void)setMessage:(NSString *)newMsg {

    label.text = newMsg;

}

                                  

- (void)alignToCenter {

    CGRect rect_screen = [[UIScreen mainScreen]bounds];
    
    CGPoint centerSuper = {rectSuper.size.width / 2, rectSuper.size.height / 2};

    CGPoint centerSelf = { self.frame.size.width / 2, self.frame.size.height / 2};

    CGPoint offset = {centerSuper.x-centerSelf.x , centerSuper.y - centerSelf.y};

    CGRect newRect = CGRectOffset(self.frame,  offset.x,offset.y);

    [self setFrame:newRect];

    rectHud = newRect;
    
    CGFloat scale_screen = [UIScreen mainScreen].scale;
                                                                                           
                                                     

}
@end
