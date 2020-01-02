#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"

static UIView *udwLogAlertBackgroundView;
static NSString *udwLogAlertMsg;
static UITextView *udwLogAlertTextView;
static BOOL udwLogAlertIsShow;
static BOOL udwLogAlertUserClickClose;

void udwLogTextViewScrollToBottom(){
        if(udwLogAlertTextView.text.length > 0) {
            NSRange bottom = NSMakeRange(udwLogAlertTextView.text.length -1, 1);
            [udwLogAlertTextView scrollRangeToVisible:bottom];
        }
}

void udwLogAlertHidden(){
    udwRunOnMainThread(^{
        udwLogAlertBackgroundView.hidden = true;
        [udwLogAlertBackgroundView removeFromSuperview];
        udwLogAlertBackgroundView = nil;
        udwLogAlertTextView.hidden = true;
        [udwLogAlertTextView removeFromSuperview];
        udwLogAlertTextView = nil;
        udwLogAlertIsShow = false;
        udwLogAlertUserClickClose = true;
    });
}

void UdwLogAlert(NSString *msg){
    udwRunOnMainThread(^{
        if (udwLogAlertUserClickClose) {
            return;
        }
        if (udwStringIsNilOrEmpty(udwLogAlertMsg)) {
            udwLogAlertMsg = @"";
        }
        udwLogAlertMsg = [NSString stringWithFormat:@"%@%@",udwLogAlertMsg,msg];
        if (udwLogAlertIsShow) {
            udwLogAlertTextView.text = udwLogAlertMsg;
            udwLogTextViewScrollToBottom();
            return;
        }
        UIView * backView = [[UIView alloc] initWithFrame:[UIScreen mainScreen].bounds];
        UIWindow* window = [UIApplication sharedApplication].keyWindow;
        if (!window) {
            window = [[UIApplication sharedApplication].windows objectAtIndex:0];
        }
        [window addSubview:backView];
        backView.backgroundColor = UIColorFromRGBA(0x000000,0.4);
        udwLogAlertBackgroundView = backView;

        CGFloat buttonHeight = 44;
        CGFloat offset = 20;
        udwLogAlertTextView = [[UITextView alloc] initWithFrame:CGRectMake(0,offset,SCREEN_WIDTH,SCREEN_HEIGHT - buttonHeight-offset-1)];
        udwLogAlertTextView.text = udwLogAlertMsg;
        udwLogAlertTextView.scrollEnabled = YES;
        udwLogAlertTextView.contentSize = [udwLogAlertTextView sizeThatFits:udwLogAlertTextView.frame.size];
        udwLogAlertTextView.editable = NO;
        [udwLogAlertBackgroundView addSubview:udwLogAlertTextView];
        udwLogTextViewScrollToBottom();

        UIButton *closeButton = [UIButton buttonWithType:UIButtonTypeCustom];
        closeButton.frame = CGRectMake(10, SCREEN_HEIGHT - buttonHeight, SCREEN_WIDTH/2 - 20, buttonHeight);
        closeButton.backgroundColor = [UIColor blueColor];
        closeButton.layer.cornerRadius = 10;
        closeButton.layer.masksToBounds = true;
        [closeButton setTitle:@"Close" forState:UIControlStateNormal];
        [closeButton udwActionHandleEvents:UIControlEventTouchUpInside action:^{
            udwLogAlertHidden();
        }];
        UIButton *copyButton = [UIButton buttonWithType:UIButtonTypeCustom];
        copyButton.frame = CGRectMake(10+SCREEN_WIDTH/2, SCREEN_HEIGHT - buttonHeight, SCREEN_WIDTH/2 - 20, buttonHeight);
        copyButton.backgroundColor = [UIColor blueColor];
        copyButton.layer.cornerRadius = 10;
        copyButton.layer.masksToBounds = true;
        [copyButton setTitle:@"Copy" forState:UIControlStateNormal];
        [copyButton udwActionHandleEvents:UIControlEventTouchUpInside action:^{
            UIPasteboard *pasteboard = [UIPasteboard generalPasteboard];
            pasteboard.string = udwLogAlertMsg;
            udwAlertOk(@"copy done",@"copy done");
        }];
        [udwLogAlertBackgroundView addSubview:closeButton];
        [udwLogAlertBackgroundView addSubview:copyButton];
        udwLogAlertIsShow = true;
    });
}

void UdwLogAlertClean(){
    udwRunOnMainThread(^{
        udwLogAlertMsg = @"";
        if (udwLogAlertIsShow) {
            udwLogAlertTextView.text = udwLogAlertMsg;
            return;
        }
    });
}
