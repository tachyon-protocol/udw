// +build ios macAppStore

#import "udwOcFoundation.h"
#import <mach/mach.h>

void report_memory(void) {
  struct task_basic_info info;
  mach_msg_type_number_t size = sizeof(info);
  kern_return_t kerr = task_info(mach_task_self(),
                                 TASK_BASIC_INFO,
                                 (task_info_t)&info,
                                 &size);
  if( kerr == KERN_SUCCESS ) {
    NSLog(@"Memory in use (in bytes): %lld", (int64_t)info.resident_size);
  } else {
    NSLog(@"Error with task_info(): %s", mach_error_string(kerr));
  }
}