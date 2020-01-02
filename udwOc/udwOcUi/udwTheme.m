#import "udwOcUi.h"
#import "src/github.com/tachyon-protocol/udw/udwOc/udwOcFoundation/udwOcFoundation.h"
#import <objc/message.h>
#import <objc/runtime.h>

static NSString * const udwThemeChangingNotification = @"udwThemeChangingNotification";
static NSString * const udwThemeAllTags = @"udwThemeAllTags";
static NSString * const udwThemeCurrentTag = @"udwThemeCurrentTag";
static NSString * const udwThemeCustomInfoFlag = @"udwThemeCustomInfoFlag";

typedef udwThemeConfigModel *(^udwThemeConfigForSelectorAndColor)(NSString *tag,SEL sel,UIColor *color);
typedef udwThemeConfigModel *(^udwThemeConfigForSelectorAndValueArray)(NSString *tag,SEL sel,NSArray *values);
typedef udwThemeConfigModel *(^udwThemeConfigForSelectorAndImage)(NSString *tag,SEL sel,UIImage *image);
typedef udwThemeConfigModel *(^udwThemeConfigForSelectorAndValues)(NSString *tag,SEL sel,...);

@interface udwTheme ()
@property (nonatomic , copy)NSString *defaultTag;
@property (nonatomic , copy)NSString *currentTag;
@property (nonatomic , strong)NSMutableArray *allTags;
@end
@implementation udwTheme
+ (void)startTheme:(NSString *)tag{
    if (![[self shareTheme].allTags containsObject:tag]) {
        udwLog(@"warning >>>> udwTheme not contain current theme %@",[self shareTheme].allTags);
        return;
    }
    if (udwStringIsNilOrEmpty(tag)){
        return;
    }
    [udwTheme shareTheme].currentTag = tag;
    [[NSNotificationCenter defaultCenter] postNotificationName:udwThemeChangingNotification object:nil userInfo:nil];
}

+ (NSString *)currentThemeTag{
    return !udwStringIsNilOrEmpty([udwTheme shareTheme].currentTag) ? [udwTheme shareTheme].currentTag : [[NSUserDefaults standardUserDefaults] objectForKey:udwThemeCurrentTag];
}

+ (void)setDefaultTheme:(NSString *)tag{
    if (udwStringIsNilOrEmpty(tag)){
        return;
    }
    [self shareTheme].defaultTag = tag;
    if (udwStringIsNilOrEmpty([self shareTheme].currentTag) && udwStringIsNilOrEmpty([[NSUserDefaults standardUserDefaults] objectForKey:udwThemeCurrentTag])) {
        [self shareTheme].currentTag = tag;
    }
}

+ (NSArray *)allThemeTag{
    return [[self shareTheme].allTags copy];
}

                      
+ (udwTheme *)shareTheme{
    static udwTheme *themeManager = nil;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        themeManager = [[udwTheme alloc]init];
    });
    return themeManager;
}

+ (void)addTagToAllTags:(NSString *)tag{
    if (![[udwTheme shareTheme].allTags containsObject:tag]) {
        [[udwTheme shareTheme].allTags addObject:tag];
        [[NSUserDefaults standardUserDefaults] setObject:[udwTheme shareTheme].allTags forKey:udwThemeAllTags];
    }
}

- (void)setCurrentTag:(NSString *)currentTag{
    if (_currentTag == currentTag) {
        return;
    }
    _currentTag = currentTag;
    [[NSUserDefaults standardUserDefaults] setObject:currentTag forKey:udwThemeCurrentTag];
}

- (NSMutableArray *)allTags{
    if (!_allTags) {
        _allTags = [NSMutableArray arrayWithArray:[[NSUserDefaults standardUserDefaults] objectForKey:udwThemeAllTags]];
    }
    return _allTags;
}
@end

@interface udwThemeConfigModel ()
@property (nonatomic, weak)UIView *superView;
@property (nonatomic, copy)udwThemeConfigForSelectorAndColor udwAddSelectorAndColor;
@property (nonatomic, copy)udwThemeConfigForSelectorAndValueArray udwAddSelectorAndValueArray;
@property (nonatomic, copy)udwThemeConfigForSelectorAndImage udwAddSelectorAndImage;
@property (nonatomic, copy)udwThemeConfigForSelectorAndValues udwAddSelectorAndValues;

@property (nonatomic, copy)void(^modelUpdateCurrentThemeConfig)(void);
@property (nonatomic, copy)void(^modelConfigThemeChangBlock)(void);
@property (nonatomic, copy)udwThemeChangeThemeBlock modelChangBlock;
@property (nonatomic, copy)NSString *modelCurrentThemeTag;
@property (nonatomic, strong)NSMutableDictionary <NSString *,NSMutableDictionary *>*modelThemeBlockConfigInfo;                             
@property (nonatomic, strong)NSMutableDictionary <NSString *,NSMutableDictionary *>*modelThemeSelectorConfigInfo;                                                                  
@end
@implementation udwThemeConfigModel
- (udwThemeConfigForColor)udwAddBackGroundColor{
    __weak typeof(self) weakSelf = self;
    return ^(NSString *tag,UIColor *color){
        typeof(self) strongSelf = weakSelf;
        return strongSelf.udwAddSelectorAndColor(tag,@selector(setBackgroundColor:),color);
    };
}

- (udwThemeConfigForBackgroundImage)udwAddImage{
    __weak typeof(self) weakSelf = self;
    return ^(NSString *tag,UIImage *image){
        typeof(self) strongSelf = weakSelf;
        return strongSelf.udwAddSelectorAndImage(tag,@selector(setImage:),image);
    };
}

- (udwThemeConfigForColor)udwAddLabelTitleColor{
    __weak typeof(self) weakSelf = self;
    return ^(NSString *tag,UIColor *color){
        typeof(self) strongSelf = weakSelf;
        return strongSelf.udwAddSelectorAndColor(tag,@selector(setTextColor:),color);
    };
}

- (udwThemeConfigForColorAndState)udwAddButtonTitleColor{
    __weak typeof(self) weakSelf = self;
    return ^(NSString *tag,UIColor *color , UIControlState state){
        typeof(self) strongSelf = weakSelf;
        return strongSelf.udwAddSelectorAndValues(tag,@selector(setTitleColor:forState:),color,@(state),nil);
    };
}

- (udwThemeConfigForImageAndState)udwAddButtonImage{
    __weak typeof(self) weakSelf = self;
    return ^(NSString *tag,UIImage *image,UIControlState state){
        typeof(self) strongSelf = weakSelf;
        return strongSelf.udwAddSelectorAndValues(tag,@selector(setImage:forState:),image,@(state),nil);
    };
}

- (udwThemeConfigForImageAndState)udwAddButtonBackgroundImage{
    __weak typeof(self) weakSelf = self;
    return ^(NSString *tag,UIImage *image,UIControlState state){
        typeof(self) strongSelf = weakSelf;
        return strongSelf.udwAddSelectorAndValues(tag,@selector(setBackgroundImage:forState:),image,@(state),nil);
    };
}

                             
- (udwThemeConfigForSelectorAndColor)udwAddSelectorAndColor{
    __weak typeof(self) weakSelf = self;
    return ^(NSString *tag,SEL sel,UIColor *color){
        typeof(self) strongSelf = weakSelf;
        if (color) {
            strongSelf.udwAddSelectorAndValueArray(tag,sel,@[color]);
        }
        return strongSelf;
    };
}

- (udwThemeConfigForSelectorAndImage)udwAddSelectorAndImage{
    __weak typeof(self) weakSelf = self;
    return ^(NSString *tag,SEL sel,UIImage *image){
        typeof(self) strongSelf = weakSelf;
        if (image) {
            strongSelf.udwAddSelectorAndValueArray(tag,sel,@[image]);
        }
        return strongSelf;
    };
}

- (udwThemeConfigForSelectorAndValues)udwAddSelectorAndValues{
    __weak typeof(self) weakSelf = self;
    return ^(NSString *tag, SEL sel,...){
        typeof(self) strongSelf = weakSelf;
        if (!sel) {
            return strongSelf;
        }
        NSMutableArray *array = [NSMutableArray array];
        va_list argsList;
        va_start(argsList, sel);
        id arg;
        while ((arg = va_arg(argsList, id))) {
            [array addObject:arg];
        }
        va_end(argsList);
        return strongSelf.udwAddSelectorAndValueArray(tag,sel,array);
    };
}

- (udwThemeConfigForSelectorAndValueArray)udwAddSelectorAndValueArray{
    __weak typeof(self) weakSelf = self;
    return ^(NSString *tag,SEL sel,NSArray *values){
        typeof(self) strongSelf = weakSelf;
        if (udwStringIsNilOrEmpty(tag)) {
            return strongSelf;
        };
        if (!sel) {
            return strongSelf;
        }
        [udwTheme addTagToAllTags:tag];
        NSString *key = NSStringFromSelector(sel);
        NSMutableDictionary *info = strongSelf.modelThemeSelectorConfigInfo[key];
        if (!info) {
            info = [NSMutableDictionary dictionary];
        }
        NSMutableArray *valuesArray = info[tag];
        if (!valuesArray) {
            valuesArray = [NSMutableArray array];
        }
        [[valuesArray copy] enumerateObjectsUsingBlock:^(NSArray *valueArray, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([valueArray isEqualToArray:values]) {
                [valuesArray removeObject:valueArray];
            }
        }];
        if (values && values.count >0){
            [valuesArray addObject:values];
        }
        [info setObject:valuesArray forKey:tag];
        [strongSelf.modelThemeSelectorConfigInfo setObject:info forKey:key];
        [strongSelf updateCurrentThemeConfigHandleWithTag:tag];
        return strongSelf;
    };
}

- (void)updateCurrentThemeConfigHandleWithTag:(NSString *)tag{
    if ([[udwTheme currentThemeTag] isEqualToString:tag]) {
        udwRunOnMainThread(^{
            if (self.modelUpdateCurrentThemeConfig) {
                self.modelUpdateCurrentThemeConfig();
            }
        });
    }
}

                
- (udwThemeConfigForCustomConfigBlock)udwAddCustomConfig{
    __weak typeof(self) weakSelf = self;
    return ^(NSString *tag,udwThemeConfigBlock configBlock){
        typeof(self) strongSelf = weakSelf;
        if (configBlock) {
            [udwTheme addTagToAllTags:tag];
            NSMutableDictionary *info = strongSelf.modelThemeBlockConfigInfo[tag];
            if (!info) {
                info = [NSMutableDictionary dictionary];
            }
            [info setObject:udwThemeCustomInfoFlag forKey:configBlock];
            [strongSelf.modelThemeBlockConfigInfo setObject:info forKey:tag];
            [strongSelf updateCurrentThemeConfigHandleWithTag:tag];
        }

        return strongSelf;
    };
}

               
- (udwThemeConfigForThemeChangBlock)udwThemeChangBlock{
    __weak typeof(self) weakSelf = self;
    return ^(udwThemeChangeThemeBlock changingBlock){
        typeof(self) strongSelf = weakSelf;
        if (changingBlock) {
            strongSelf.modelChangBlock = changingBlock;
            if (strongSelf.modelConfigThemeChangBlock) {
                strongSelf.modelConfigThemeChangBlock();
            }
        }
        return strongSelf;
    };
}

                     
- (NSMutableDictionary *)modelThemeBlockConfigInfo{
    if (!_modelThemeBlockConfigInfo) {
        _modelThemeBlockConfigInfo = [NSMutableDictionary dictionary];
    }
    return _modelThemeBlockConfigInfo;
}

- (NSMutableDictionary *)modelThemeSelectorConfigInfo{
    if (!_modelThemeSelectorConfigInfo) {
        _modelThemeSelectorConfigInfo = [NSMutableDictionary dictionary];
    }
    return _modelThemeSelectorConfigInfo;
}

- (void)dealloc{
    [[NSNotificationCenter defaultCenter] removeObserver:self];
    objc_removeAssociatedObjects(self);
    _modelCurrentThemeTag = nil;
    _modelThemeBlockConfigInfo = nil;
    _modelThemeSelectorConfigInfo = nil;
    [[NSNotificationCenter defaultCenter] removeObserver:self.superView name:udwThemeChangingNotification object:nil];
    objc_removeAssociatedObjects(self.superView);
    udwLog(@"----------dealloc %@----------",NSStringFromClass([self class]));
}
@end

static void * udwThemeModelKey = &udwThemeModelKey;
@implementation UIView (udwThemeConfig)
@dynamic udwThemeModel;
- (udwThemeConfigModel *)udwThemeModel{
    udwThemeConfigModel *model = objc_getAssociatedObject(self, udwThemeModelKey);
    if (!model) {
        model = [[udwThemeConfigModel alloc] init];
        model.superView = self;
        objc_setAssociatedObject(self, udwThemeModelKey, model , OBJC_ASSOCIATION_RETAIN_NONATOMIC);
        [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(udwThemeChangeThemeConfigNotify:) name:udwThemeChangingNotification object:nil];
        __weak typeof(self) weakSelf = self;
        model.modelUpdateCurrentThemeConfig = ^{
            typeof(self) strongSelf = weakSelf;
            if (strongSelf) {
                [strongSelf changeThemeConfig];
            }
        };
        model.modelConfigThemeChangBlock = ^{
            typeof(self) strongSelf = weakSelf;
            if (strongSelf) {
                strongSelf.udwThemeModel.modelChangBlock([udwTheme currentThemeTag], strongSelf);
            }
        };
    }
    return model;
}

- (void)setUdwThemeModel:(udwThemeConfigModel *)udwThemeModel{
    if (self) {
        objc_setAssociatedObject(self, udwThemeModelKey, udwThemeModel , OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    }
}

- (BOOL)isChangeTheme{
    return (!self.udwThemeModel.modelCurrentThemeTag || ![self.udwThemeModel.modelCurrentThemeTag isEqualToString:[udwTheme currentThemeTag]]) ? true : false;
}

- (void)udwThemeChangeThemeConfigNotify:(NSNotification *)notify{
    udwRunOnMainThread(^{
        if ([self isChangeTheme]) {
            if (self.udwThemeModel.modelChangBlock) {
                self.udwThemeModel.modelChangBlock([udwTheme currentThemeTag] , self);
            }
            [CATransaction begin];
            [CATransaction setDisableActions:true];
            [self changeThemeConfig];
            [CATransaction commit];
        }
    });
}

- (void)changeThemeConfig{
    self.udwThemeModel.modelCurrentThemeTag = [udwTheme currentThemeTag];
    NSString *tag = [udwTheme currentThemeTag];
            
    for (id blockKey in self.udwThemeModel.modelThemeBlockConfigInfo[tag]) {
        id value = self.udwThemeModel.modelThemeBlockConfigInfo[tag][blockKey];
        if ([value isEqualToString:udwThemeCustomInfoFlag]) {
            udwThemeConfigBlock callback = (udwThemeConfigBlock)blockKey;
            if (callback) {
                callback(self);
            }
        }
    }
               
    for (NSString *selector in self.udwThemeModel.modelThemeSelectorConfigInfo) {
        NSDictionary *info = self.udwThemeModel.modelThemeSelectorConfigInfo[selector];
        NSArray *valuesArray = info[tag];
        for (NSArray *values in valuesArray) {
            SEL sel = NSSelectorFromString(selector);
            if (![self respondsToSelector:sel]) {
                continue;
            }
            NSMethodSignature * sig = [self methodSignatureForSelector:sel];
            if (!sig) {
                [self doesNotRecognizeSelector:sel];
            }
            NSInvocation *inv = [NSInvocation invocationWithMethodSignature:sig];
            if (!inv) {
                [self doesNotRecognizeSelector:sel];
            }
            [inv setTarget:self];
            [inv setSelector:sel];
            if (sig.numberOfArguments == values.count + 2) {
                [values enumerateObjectsUsingBlock:^(id  _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
                    NSInteger index = idx + 2;
                    [self setInv:inv Sig:sig Obj:obj Index:index];
                }];
                [inv invoke];
            } else {
                NSAssert(true, @"values.count does not match the count of parameters and method parameters");
            }
        }
    }
}

- (void)setInv:(NSInvocation *)inv Sig:(NSMethodSignature *)sig Obj:(id)obj Index:(NSInteger)index{
    if (sig.numberOfArguments <= index) {
        return;
    }
    char *type = (char *)[sig getArgumentTypeAtIndex:index];
    while (*type == 'r' ||         
           *type == 'n' ||      
           *type == 'N' ||         
           *type == 'o' ||       
           *type == 'O' ||          
           *type == 'R' ||         
           *type == 'V') {          
        type++;                         
    }
    BOOL unsupportedType = NO;
    switch (*type) {
        case 'v':           
        case 'B':           
        case 'c':                  
        case 'C':                    
        case 's':            
        case 'S':                     
        case 'i':                             
        case 'I':                                       
        case 'l':                  
        case 'L':                           
        {                                                 
            int value = [obj intValue];
            [inv setArgument:&value atIndex:index];
        } break;

        case 'q':                                                 
        case 'Q':                                                                    
        {
            long long value = [obj longLongValue];
            [inv setArgument:&value atIndex:index];
        } break;

        case 'f':                             
        {                                         
            double value = [obj doubleValue];
            float valuef = value;
            [inv setArgument:&valuef atIndex:index];
        } break;

        case 'd':                              
        {
            double value = [obj doubleValue];
            [inv setArgument:&value atIndex:index];
        } break;
        case '*':          
        case '^':           
        {
            if ([obj isKindOfClass:UIColor.class]) obj = (id)[obj CGColor];                 
            if ([obj isKindOfClass:UIImage.class]) obj = (id)[obj CGImage];                 
            void *value = (__bridge void *)obj;
            [inv setArgument:&value atIndex:index];
        } break;
        case '@':      
        {
            id value = obj;
            [inv setArgument:&value atIndex:index];
        } break;
        case '{':          
        {
            if (strcmp(type, @encode(CGPoint)) == 0) {
                CGPoint value = [obj CGPointValue];
                [inv setArgument:&value atIndex:index];
            } else if (strcmp(type, @encode(CGSize)) == 0) {
                CGSize value = [obj CGSizeValue];
                [inv setArgument:&value atIndex:index];
            } else if (strcmp(type, @encode(CGRect)) == 0) {
                CGRect value = [obj CGRectValue];
                [inv setArgument:&value atIndex:index];
            } else if (strcmp(type, @encode(CGVector)) == 0) {
                CGVector value = [obj CGVectorValue];
                [inv setArgument:&value atIndex:index];
            } else if (strcmp(type, @encode(CGAffineTransform)) == 0) {
                CGAffineTransform value = [obj CGAffineTransformValue];
                [inv setArgument:&value atIndex:index];
            } else if (strcmp(type, @encode(CATransform3D)) == 0) {
                CATransform3D value = [obj CATransform3DValue];
                [inv setArgument:&value atIndex:index];
            } else if (strcmp(type, @encode(NSRange)) == 0) {
                NSRange value = [obj rangeValue];
                [inv setArgument:&value atIndex:index];
            } else if (strcmp(type, @encode(UIOffset)) == 0) {
                UIOffset value = [obj UIOffsetValue];
                [inv setArgument:&value atIndex:index];
            } else if (strcmp(type, @encode(UIEdgeInsets)) == 0) {
                UIEdgeInsets value = [obj UIEdgeInsetsValue];
                [inv setArgument:&value atIndex:index];
            } else {
                unsupportedType = true;
            }
        } break;
        case '(':         
        {
            unsupportedType = true;
        } break;

        case '[':         
        {
            unsupportedType = true;
        } break;

        default:          
        {
            unsupportedType = true;
        } break;
    }
    NSAssert(unsupportedType == false, @"Methods of parameter type does not support");
}
@end

