#include "udwGaocLib.h"

#define LOG_INFO(...) NSLog(__VA_ARGS__);
#define udwIntSize (sizeof(size_t))

void udwGaocl_c_ensureAddSize(udwGaocGoBuffer *seq,size_t size);
static void udwGaocl_c_align(udwGaocGoBuffer* _buf,size_t align){
    size_t pad = _buf->off % align;
    if (pad > 0) {
    	_buf->off += align - pad;
    }
}

NSString *udwGaocl_c_readString(udwGaocGoBuffer *seq) {
  udwGaocl_c_align(seq,udwIntSize);
  size_t len = *(size_t*)(&seq->buf[seq->off]);
  seq->off += udwIntSize;
  if (len == 0) {                  
    return @"";
  }
  const void *buf = (const void *)(&seq->buf[seq->off]);
  seq->off += len;
  return [[NSString alloc] initWithBytes:buf
                                  length:len
                                encoding:NSUTF8StringEncoding];
}

void udwGaocl_c_writeString(udwGaocGoBuffer *seq, NSString *s) {
  size_t len = [s lengthOfBytesUsingEncoding:NSUTF8StringEncoding];
  if (len == 0 && s.length > 0) {
      LOG_INFO(@"unable to encode an NSString into UTF-8");
      return;
  }
  udwGaocl_c_align(seq,udwIntSize);
  udwGaocl_c_ensureAddSize(seq,udwIntSize+len);
  *(size_t*)(&seq->buf[seq->off]) = len;
  seq->off+=udwIntSize;
  NSUInteger used;
  [s getBytes: (char *)&seq->buf[seq->off]
         maxLength:(NSUInteger)len
        usedLength:&used
          encoding:NSUTF8StringEncoding
           options:0
             range:NSMakeRange(0, [s length])
    remainingRange:NULL];

  seq->off+=len;
  return;
}

size_t udwGaocl_c_readInt(udwGaocGoBuffer *seq) {
  udwGaocl_c_align(seq,udwIntSize);
  size_t out = *(size_t*)(&seq->buf[seq->off]);
  seq->off += udwIntSize;
  return out;
}

void udwGaocl_c_writeInt(udwGaocGoBuffer *seq,size_t v){
    udwGaocl_c_align(seq,udwIntSize);
    udwGaocl_c_ensureAddSize(seq,udwIntSize);
    *(size_t*)(&seq->buf[seq->off]) = v;
    seq->off+=udwIntSize;
}

int64_t udwGaocl_c_readInt64(udwGaocGoBuffer *seq) {
  udwGaocl_c_align(seq,8);
  int64_t out = *(int64_t*)(&seq->buf[seq->off]);
  seq->off += 8;
  return out;
}

void udwGaocl_c_writeInt64(udwGaocGoBuffer *seq,int64_t v){
    udwGaocl_c_align(seq,8);
    udwGaocl_c_ensureAddSize(seq,udwIntSize);
    *(int64_t*)(&seq->buf[seq->off]) = v;
    seq->off+=8;
}

bool udwGaocl_c_readBool(udwGaocGoBuffer *seq){
   uint8_t outB = seq->buf[seq->off];
   seq->off += 1;
   return (outB==1);
}
void udwGaocl_c_writeBool(udwGaocGoBuffer *seq, bool v){
    udwGaocl_c_ensureAddSize(seq,1);
    if (v==true){
        seq->buf[seq->off] = 1;
    }else{
        seq->buf[seq->off] = 0;
    }
    seq->off += 1;
}

double udwGaocl_c_readFloat64(udwGaocGoBuffer *seq) {
    udwGaocl_c_align(seq,8);
  double out = *(double*)(&seq->buf[seq->off]);
  seq->off += 8;
  return out;
}

void udwGaocl_c_writeFloat64(udwGaocGoBuffer *seq,double v){
    udwGaocl_c_align(seq,8);
    udwGaocl_c_ensureAddSize(seq,8);
    *(double*)(&seq->buf[seq->off]) = v;
    seq->off+=8;
}

float udwGaocl_c_readFloat32(udwGaocGoBuffer *seq) {
    udwGaocl_c_align(seq,4);
  double out = *(float*)(&seq->buf[seq->off]);
  seq->off += 4;
  return out;
}

void udwGaocl_c_writeFloat32(udwGaocGoBuffer *seq,float v){
    udwGaocl_c_align(seq,4);
    udwGaocl_c_ensureAddSize(seq,4);
    *(float*)(&seq->buf[seq->off]) = v;
    seq->off+=4;
}

NSData* udwGaocl_c_readByteSlice(udwGaocGoBuffer *seq) {
  udwGaocl_c_align(seq,udwIntSize);
  size_t len = *(size_t*)(&seq->buf[seq->off]);
  seq->off += udwIntSize;
  if (len == 0) {                  
    return [NSData dataWithBytes:nil length:0];
  }
  const void *buf = (const void *)(&seq->buf[seq->off]);
  seq->off += len;
  return [NSData dataWithBytes:buf length:len];
}

void udwGaocl_c_writeByteSlice(udwGaocGoBuffer *seq,NSData* v){
  udwGaocl_c_align(seq,udwIntSize);
  size_t len = v.length;
  udwGaocl_c_ensureAddSize(seq,udwIntSize+len);
  *(size_t*)(&seq->buf[seq->off]) = len;
  seq->off+=udwIntSize;
  memcpy((void*)(&seq->buf[seq->off]),v.bytes,len);

  seq->off+=len;
  return;
}