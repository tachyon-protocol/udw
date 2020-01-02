#import "udwOcUi.h"

#define FLEXIBLE_W SCREEN_WIDTH /375
#define FLEXIBLE_H SCREEN_HEIGHT /667

                              
CGFloat udwFixViewSizeWithOriginalFloatAndIsFlexibleHeight(CGFloat originalFloat,BOOL isHeight) {
  if (isHeight) {
    return originalFloat * FLEXIBLE_H;
  }
  return originalFloat * FLEXIBLE_W;
}

                                             
CGFloat udwFixViewSizeWithViewWidth(CGFloat width,CGSize size){
  if (size.width == 0) {
      return 0;
  }
  return width*size.height/size.width;
}

                                            
CGFloat udwFixViewSizeWithViewHeight(CGFloat height,CGSize size){
  if (size.height == 0) {
      return 0;
  }
  return height*size.width/size.height;
}

void udwLabelFixSize(UILabel *label){
  CGSize maxSize = [label sizeThatFits:CGSizeMake(MAXFLOAT, label.frame.size.height)];
  label.frame = CGRectMake(label.frame.origin.x, label.frame.origin.y, maxSize.width, maxSize.height);
}

CGRect udwFrameChangeX(CGRect frame,CGFloat newX){
  CGRect newFrame = frame;
  newFrame.origin.x = newX;
  return newFrame;
}

CGRect udwFrameChangeY(CGRect frame,CGFloat newY){
  CGRect newFrame = frame;
  newFrame.origin.y = newY;
  return newFrame;
}

CGRect udwFrameChangeWidth(CGRect frame,CGFloat width){
  CGRect newFrame = frame;
  newFrame.size.width = width;
  return newFrame;
}

CGRect udwFrameChangeHeight(CGRect frame,CGFloat height){
  CGRect newFrame = frame;
  newFrame.size.height = height;
  return newFrame;
}

CGRect udwFrameChangeSize(CGRect frame,CGSize newSize){
  CGRect newFrame = frame;
  newFrame.size = newSize;
  return newFrame;
}

void udwSetAdjustsFontForWidth(UIView *view){
    for(UIView * subView in view.subviews) {
        if([subView isKindOfClass:[UILabel class]]) {
            UILabel * label = (UILabel *)subView;
            label.adjustsFontSizeToFitWidth = TRUE;
        } else {
            udwSetAdjustsFontForWidth(subView);
        }
    }
}