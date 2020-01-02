#include "udwGaocLib.h"

__attribute__((visibility("hidden"))) void udwGaocl_c_ensureAddSize(udwGaocGoBuffer *seq,size_t size){
    size_t needSize = seq->off+size;
    if (seq->cap >= needSize){
        return;
    }
    size_t toAllocSize = seq->cap;
    if (toAllocSize<64){
        toAllocSize = 64;
    }
    while(true){
        if (toAllocSize>=needSize+64){
            break;
        }
        toAllocSize = toAllocSize*2;
    }
                                              
                                                   
                                                                  
            
        seq->buf = (uint8_t*)realloc((void*)seq->buf,toAllocSize);
       
    seq->cap = toAllocSize;
}

void udwGaocl_c_init(udwGaocGoBuffer* seq){
                                   
                                          
}
void udwGaocl_c_free(udwGaocGoBuffer *seq){
    if (seq->buf!=nil){
                                                    
            free(seq->buf);
           
        seq->buf = nil;
    }
}