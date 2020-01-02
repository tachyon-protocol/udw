#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"

@interface udwSizeToFitLabel()
@property (nonatomic,assign)CGFloat lineOffset;
@property (nonatomic,assign)CGRect originalFrame;
@end

@implementation udwSizeToFitLabel

                     
- (instancetype)initWithFrame:(CGRect)frame font:(UIFont *)font string:(NSString *)context offset:(CGFloat)lineOffset{
    self = [super initWithFrame:frame];
    if (self) {
        self.originalFrame = frame;
        self.lineOffset = lineOffset;
        self.font = font;
        self.numberOfLines = 0;
        self.lineBreakMode = NSLineBreakByWordWrapping;
        CGSize size = [self sizeThatFits:CGSizeMake(self.frame.size.width, MAXFLOAT)];
        CGRect newFrame = frame;
        newFrame.size.height = size.height;
        self.frame = newFrame;
                                  
        NSMutableAttributedString *attributedString = [[NSMutableAttributedString alloc] initWithString:context];
        NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc] init];
        [paragraphStyle setLineSpacing:lineOffset];
        [attributedString addAttribute:NSParagraphStyleAttributeName value:paragraphStyle range:NSMakeRange(0, [context length])];
        [paragraphStyle setAlignment:NSTextAlignmentCenter];
        [self setAttributedText:attributedString];
        self.textAlignment = NSTextAlignmentLeft;
        self.textColor = [UIColor whiteColor];
        [self sizeToFit];
    }
    return self;
}

- (instancetype)initWithFrame:(CGRect)frame font:(UIFont *)font textColor:(UIColor *)textColor context:(NSString *)context lineOffset:(CGFloat)lineOffset differentColor:(UIColor *)differentColor colorRange:(NSRange)colorRange{
    self = [super initWithFrame:frame];
    if (self) {
        self.originalFrame = frame;
        self.lineOffset = lineOffset;
        self.font = font;
        self.numberOfLines = 0;
        self.lineBreakMode = NSLineBreakByWordWrapping;
        CGSize size = [self sizeThatFits:CGSizeMake(self.frame.size.width, MAXFLOAT)];
        CGRect newFrame = frame;
        newFrame.size.height = size.height;
        self.frame = newFrame;
        self.textColor = textColor;
                                  
        NSMutableAttributedString *attributedString = [[NSMutableAttributedString alloc] initWithString:context];
        NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc] init];
        [paragraphStyle setLineSpacing:lineOffset];
        [attributedString addAttribute:NSParagraphStyleAttributeName value:paragraphStyle range:NSMakeRange(0, [context length])];
        [attributedString addAttribute:NSForegroundColorAttributeName
                                 value:differentColor
                                 range:colorRange];
        [paragraphStyle setAlignment:NSTextAlignmentCenter];
        [self setAttributedText:attributedString];
        self.textAlignment = NSTextAlignmentLeft;
        [self sizeToFit];
    }
    return self;
}

- (instancetype _Nonnull)initWithFrame:(CGRect)frame font:(UIFont *_Nonnull)font textColor:(UIColor *_Nonnull)textColor context:(NSString *_Nonnull)context lineOffset:(CGFloat)lineOffset attributedStringCallback:(void (^ __nullable)(NSMutableAttributedString * _Nullable attributedString))attributedStringCallback{
    self = [super initWithFrame:frame];
    if (self) {
        self.originalFrame = frame;
        self.lineOffset = lineOffset;
        self.font = font;
        self.numberOfLines = 0;
        self.lineBreakMode = NSLineBreakByWordWrapping;
        CGSize size = [self sizeThatFits:CGSizeMake(self.frame.size.width, MAXFLOAT)];
        CGRect newframe = frame;
        newframe.size.height = size.height;
        self.frame = newframe;
        self.textColor = textColor;
                                  
        NSMutableAttributedString *attributedString = [[NSMutableAttributedString alloc] initWithString:context];
        NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc] init];
        [paragraphStyle setLineSpacing:lineOffset];
        [attributedString addAttribute:NSParagraphStyleAttributeName value:paragraphStyle range:NSMakeRange(0, [context length])];
        if (attributedStringCallback !=nil) {
        attributedStringCallback(attributedString);
        }
        [paragraphStyle setAlignment:NSTextAlignmentCenter];
        [self setAttributedText:attributedString];
        self.textAlignment = NSTextAlignmentLeft;
        [self sizeToFit];
    }
    return self;
}

- (void)changeTextWithNewContext:(NSString *)context{
    [self changeTextWithNewContext:context withUnderLine:false];
}

- (void)changeTextWithNewContext:(NSString *)context withUnderLine:(BOOL)withUnderLine{
    CGSize size = [self sizeThatFits:CGSizeMake(self.originalFrame.size.width, MAXFLOAT)];
    CGRect newFrame = self.originalFrame;
    newFrame.size.height = size.height;
    self.frame = newFrame;
                              
    NSMutableAttributedString *attributedString = [[NSMutableAttributedString alloc] initWithString:context];
    NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc] init];
    [paragraphStyle setLineSpacing:self.lineOffset];
    [attributedString addAttribute:NSParagraphStyleAttributeName value:paragraphStyle range:NSMakeRange(0, [context length])];
    if (withUnderLine) {
        [attributedString addAttribute:NSUnderlineStyleAttributeName value:@(NSUnderlineStyleSingle) range:NSMakeRange(0, [context length])];
    }
    [self setAttributedText:attributedString];
    [self sizeToFit];
}

- (void)changeLabelWithModel:(udwSizeToFitLabelModel *)model {
    CGSize size = [self sizeThatFits:CGSizeMake(self.originalFrame.size.width, MAXFLOAT)];
    CGRect newFrame = self.originalFrame;
    newFrame.size.height = size.height;
    self.frame = newFrame;
    NSString *context = self.attributedText.string;
    if (!udwStringIsNilOrEmpty(model.labelStr)) {
        context = model.labelStr;
    }
    NSMutableAttributedString *attributedString = [[NSMutableAttributedString alloc] initWithString:context];
    if (model.labelLineH != 0) {
        self.lineOffset = model.labelLineH;
    }
    NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc] init];
    [paragraphStyle setLineSpacing:self.lineOffset];
    [attributedString addAttribute:NSParagraphStyleAttributeName value:paragraphStyle range:NSMakeRange(0, [self.attributedText length])];
    if (model.attributedStringCallback) {
    model.attributedStringCallback(attributedString);
    }
    [paragraphStyle setAlignment:NSTextAlignmentCenter];
    [self setAttributedText:attributedString];
    self.textAlignment = NSTextAlignmentLeft;
    [self sizeToFit];
}

@end

@implementation udwSizeToFitLabelModel
@end