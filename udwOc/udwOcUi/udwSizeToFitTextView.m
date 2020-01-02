#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"

@implementation udwClickTextModel
@end
udwClickTextModel* udwNewClickTextModel(NSString *clickStr,UIFont *clickFont,UIColor *clickStrColor, void(^onClick)()) {
    udwClickTextModel *clickText = [udwClickTextModel new];
    clickText.clickStr = clickStr;
    clickText.clickFont = clickFont;
    clickText.clickStrColor = clickStrColor;
    clickText.onClick = onClick;
    return clickText;
}

@interface udwSizeToFitTextView()
@property (nonatomic,assign)CGFloat lineOffset;
@property (nonatomic,assign)CGRect originalFrame;
@property (nonatomic,strong)NSArray<udwClickTextModel *> *clickTextList;
@end

@implementation udwSizeToFitTextView

- (instancetype _Nonnull)initWithFrame:(CGRect)frame font:(UIFont *_Nonnull)font textColor:(UIColor *_Nonnull)textColor context:(NSString *_Nonnull)context lineOffset:(CGFloat)lineOffset attributedStringCallback:(void (^ __nullable)(NSMutableAttributedString * _Nullable attributedString))attributedStringCallback{
    self = [super initWithFrame:frame];
    if (self) {
        self.originalFrame = frame;
        self.lineOffset = lineOffset;
        self.font = font;
        CGSize size = [self sizeThatFits:CGSizeMake(self.frame.size.width, MAXFLOAT)];
        CGRect newFrame = frame;
        newFrame.size.height = size.height;
        self.frame = newFrame;
        self.textColor = textColor;
        self.textContainerInset = UIEdgeInsetsZero;
        self.textContainer.lineFragmentPadding = 0;
        self.text = context;
        NSMutableAttributedString *attributedString = [[NSMutableAttributedString alloc] initWithString:context];
        NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc] init];
        [paragraphStyle setLineSpacing:lineOffset];
        [attributedString addAttribute:NSParagraphStyleAttributeName value:paragraphStyle range:NSMakeRange(0, [context length])];
        [attributedString addAttribute:NSFontAttributeName value:font range:NSMakeRange(0, [context length])];
        [attributedString addAttribute:NSForegroundColorAttributeName value:textColor range:NSMakeRange(0, [context length])];
        if (attributedStringCallback) {
           attributedStringCallback(attributedString);
         }
        [paragraphStyle setAlignment:NSTextAlignmentCenter];
        [self setAttributedText:attributedString];
        self.textAlignment = NSTextAlignmentLeft;
        self.backgroundColor = [UIColor clearColor];
        self.scrollEnabled = false;
        self.editable = false;
        self.dataDetectorTypes = UIDataDetectorTypeAll;
        [self sizeToFit];
    }
    return self;
}

- (instancetype _Nonnull)initWithFrame:(CGRect)frame font:(UIFont *_Nonnull)font textColor:(UIColor *_Nonnull)textColor context:(NSString *_Nonnull)context lineOffset:(CGFloat)lineOffset attributedStringCallback:(void (^ __nullable)(NSMutableAttributedString * _Nullable attributedString))attributedStringCallback clickTextArr:(NSArray<udwClickTextModel *> *_Nullable)clickTextArr{
    self = [super initWithFrame:frame];
    if (self) {
        self.originalFrame = frame;
        self.lineOffset = lineOffset;
        self.clickTextList = clickTextArr;
        self.font = font;
        CGSize size = [self sizeThatFits:CGSizeMake(self.frame.size.width, MAXFLOAT)];
        CGRect newFrame = frame;
        newFrame.size.height = size.height;
        self.frame = newFrame;
        self.textColor = textColor;
        self.textContainerInset = UIEdgeInsetsZero;
        self.textContainer.lineFragmentPadding = 0;
        self.text = context;
        NSMutableAttributedString *attributedString = [[NSMutableAttributedString alloc] initWithString:context];
        NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc] init];
        [paragraphStyle setLineSpacing:lineOffset];
        [attributedString addAttribute:NSParagraphStyleAttributeName value:paragraphStyle range:NSMakeRange(0, [context length])];
        [attributedString addAttribute:NSFontAttributeName value:font range:NSMakeRange(0, [context length])];
        [attributedString addAttribute:NSForegroundColorAttributeName value:textColor range:NSMakeRange(0, [context length])];
         if (attributedStringCallback) {
         attributedStringCallback(attributedString);
         }
        for (udwClickTextModel *clickText in clickTextArr) {
            NSRange clickRange = [context rangeOfString:clickText.clickStr];
            if (clickRange.location != NSNotFound) {
                [attributedString addAttribute:clickText.clickStr value:@(true) range:clickRange];
                [attributedString addAttribute:NSFontAttributeName value:clickText.clickFont range:clickRange];
                [attributedString addAttribute:NSForegroundColorAttributeName value:clickText.clickStrColor range:clickRange];
            }
        }
        [paragraphStyle setAlignment:NSTextAlignmentCenter];
        [self setAttributedText:attributedString];
        self.textAlignment = NSTextAlignmentLeft;
        self.backgroundColor = [UIColor clearColor];
        self.scrollEnabled = false;
        self.editable = false;
        self.dataDetectorTypes = UIDataDetectorTypeNone;
        self.selectable = false;
        UITapGestureRecognizer *tap = [[UITapGestureRecognizer alloc] initWithTarget:self action:@selector(textTapped:)];
        [self addGestureRecognizer:tap];
        [self sizeToFit];
    }
    return self;
}

- (void)textTapped:(UITapGestureRecognizer *)recognizer
{
    UITextView *textView = (UITextView *)recognizer.view;
    NSLayoutManager *layoutManager = textView.layoutManager;
    CGPoint location = [recognizer locationInView:textView];
    location.x -= textView.textContainerInset.left;
    location.y -= textView.textContainerInset.top;
    NSUInteger characterIndex;
    characterIndex = [layoutManager characterIndexForPoint:location
                                           inTextContainer:textView.textContainer
                  fractionOfDistanceBetweenInsertionPoints:NULL];
    if (characterIndex < textView.textStorage.length) {
        for (udwClickTextModel *clickText in self.clickTextList) {
            NSRange range;
            id value = [textView.attributedText attribute:clickText.clickStr atIndex:characterIndex effectiveRange:&range];
            if (value) {
                clickText.onClick();
            }
        }
    }
}

- (void)changeTextWithNewContext:(NSString *)context {
    CGSize size = [self sizeThatFits:CGSizeMake(self.originalFrame.size.width, MAXFLOAT)];
    CGRect newFrame = self.originalFrame;
    newFrame.size.height = size.height;
    self.frame = newFrame;
                              
    NSMutableAttributedString *attributedString = [[NSMutableAttributedString alloc] initWithString:context];
    NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc] init];
    [paragraphStyle setLineSpacing:self.lineOffset];
    [attributedString addAttribute:NSParagraphStyleAttributeName value:paragraphStyle range:NSMakeRange(0, [context length])];
    [attributedString addAttribute:NSFontAttributeName value:self.font range:NSMakeRange(0, [context length])];
    [self setAttributedText:attributedString];
    [self sizeToFit];
}

@end