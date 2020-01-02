#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"
@interface kouLoadingButton()
@property (nonatomic,strong)UIActivityIndicatorView *loadingView;
@end

@implementation kouLoadingButton

-(void) LoadingStatus:(BOOL) loading {
    [self bringSubviewToFront:self.loadingView];
    if (loading) {
        [self.loadingView startAnimating];
        [self setEnabled:NO];
    }else {
        [self.loadingView stopAnimating];
        [self setEnabled:YES];
    }
}

- (void)setLoadingColor:(UIColor *)loadingColor{
  if (loadingColor == _loadingColor) {
    return;
  }
      _loadingColor = loadingColor;
  self.loadingView.color = loadingColor;
}

- (UIActivityIndicatorView *)loadingView{
  if (_loadingView == nil) {
    _loadingView = [[UIActivityIndicatorView alloc] initWithFrame:CGRectMake(self.frame.size.width-self.frame.size.height, 0, self.frame.size.height, self.frame.size.height)];
    [_loadingView setHidesWhenStopped:YES];
    [self addSubview:_loadingView];
    }
    return _loadingView;
}

@end

