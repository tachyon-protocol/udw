// +build ios macAppStore

#import "udwOcFoundation.h"

#define ARC4RANDOM_MAX      0x100000000

double udwRandomBetween0And1(){
    return (double)arc4random() / ARC4RANDOM_MAX;
}