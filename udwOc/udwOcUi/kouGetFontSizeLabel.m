#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"

@implementation kouGetFontSizeLabel

                                 
- (CGFloat)fontSize{
    if (self.adjustsFontSizeToFitWidth) {
        UIFont *currentFont = self.font;
        CGFloat originalFontSize = currentFont.pointSize;
        CGSize currentSize = [self.text sizeWithAttributes:@{NSFontAttributeName : currentFont}];
        while (currentSize.width > self.frame.size.width && currentFont.pointSize > (originalFontSize * self.minimumScaleFactor)) {
            currentFont = [currentFont fontWithSize:currentFont.pointSize - 1];
            currentSize = [self.text sizeWithAttributes:@{NSFontAttributeName : currentFont}];
        }
        return currentFont.pointSize;
    }
    return self.font.pointSize;
}

@end