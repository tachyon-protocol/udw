#import "udwOcUi.h"

void udwRegisterNibClass(){
    [kouLoadingButton class];
    [kouGetFontSizeLabel class];
    [kouPointInsideButton class];
}

void udwLogAllFontName(){
    for (NSString *fontFamilyName in [UIFont familyNames]) {
        for (NSString *fontName in [UIFont fontNamesForFamilyName:fontFamilyName]) {
            NSLog(@"font>>>>>%@",fontName);
        }
        NSLog(@"-----------------");
    }
}