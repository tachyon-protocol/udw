#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"

UIColor *UIColorFromRGBHexStr(NSString *HexStr){
  unsigned result = 0;
  NSScanner *scanner = [NSScanner scannerWithString:HexStr];
  [scanner scanHexInt:&result];
  return UIColorFromRGB(result);
}