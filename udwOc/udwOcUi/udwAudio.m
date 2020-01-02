#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"


                              
AVAuthorizationStatus udwGetAudioStatus() {
    return [AVCaptureDevice authorizationStatusForMediaType:AVMediaTypeAudio];
}

void udwRequestRecordPermission (void (^statusCallback)(bool allow)) {
    [[AVAudioSession sharedInstance] requestRecordPermission:^(BOOL granted) {
        statusCallback(granted);
    }];
}
