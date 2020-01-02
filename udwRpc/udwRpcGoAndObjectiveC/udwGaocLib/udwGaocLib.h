#include <Foundation/Foundation.h>

                               

                                     
                                                                 
                                          
                                                                   
typedef struct udwGaocGoBuffer {
  uint8_t *buf;
  size_t cap;
  size_t off;
                                               
} udwGaocGoBuffer;

NSString *udwGaocl_c_readString(udwGaocGoBuffer *seq);
void udwGaocl_c_writeString(udwGaocGoBuffer *seq, NSString *s);

size_t udwGaocl_c_readInt(udwGaocGoBuffer *seq);
void udwGaocl_c_writeInt(udwGaocGoBuffer *seq, size_t v);

int64_t udwGaocl_c_readInt64(udwGaocGoBuffer *seq);
void udwGaocl_c_writeInt64(udwGaocGoBuffer *seq, int64_t v);

bool udwGaocl_c_readBool(udwGaocGoBuffer *seq);
void udwGaocl_c_writeBool(udwGaocGoBuffer *seq, bool v);

void udwGaocl_c_init(udwGaocGoBuffer* seq);
void udwGaocl_c_free(udwGaocGoBuffer* seq);

double udwGaocl_c_readFloat64(udwGaocGoBuffer *seq);
void udwGaocl_c_writeFloat64(udwGaocGoBuffer *seq,double v);

float udwGaocl_c_readFloat32(udwGaocGoBuffer *seq);
void udwGaocl_c_writeFloat32(udwGaocGoBuffer *seq,float v);

NSData* udwGaocl_c_readByteSlice(udwGaocGoBuffer *seq);
void udwGaocl_c_writeByteSlice(udwGaocGoBuffer *seq,NSData* v);