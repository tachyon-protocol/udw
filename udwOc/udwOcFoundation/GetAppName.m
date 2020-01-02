// +build ios macAppStore

#include "udwOcFoundation.h"

NSString* UdwOcGetAppName(){
    return [[NSBundle mainBundle] infoDictionary][(NSString *) kCFBundleNameKey];
}