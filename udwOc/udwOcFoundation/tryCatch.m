// +build ios macAppStore

#import "udwOcFoundation.h"

void udwOcTryCatch(void (^tryCallback)(),void (^catchCallback)(NSException *exception),void (^finallyCallback)()){
  @try {
  tryCallback();
  } @catch (NSException *exception) {
  catchCallback(exception);
  } @finally {
  finallyCallback();
  }
}

