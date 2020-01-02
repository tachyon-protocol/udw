// +build ios macAppStore

#import "udwOcFoundation.h"

NSString* udwGetClassNameFromInstance(id obj){
    return NSStringFromClass([obj class]);
}
