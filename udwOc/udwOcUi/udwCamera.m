#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"


                              
AVAuthorizationStatus udwGetCameraStatus() {
    return [AVCaptureDevice authorizationStatusForMediaType:AVMediaTypeVideo];
}

void udwRequestCameraPermission (void (^statusCallback)(bool allow)) {
    [AVCaptureDevice requestAccessForMediaType:AVMediaTypeVideo completionHandler:^(BOOL granted) {
            statusCallback(granted);
        }];
}
