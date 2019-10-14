#include "textflag.h"

// func asmBulkAddSum(bufp uintptr,loopNum uintptr) uint64
// bufp 8(SP)
// out 0x18(SP)
TEXT ·asmBulkAddSum(SB),NOSPLIT,$0
    // var bufp uintptr -> BX
    MOVQ 8(SP),BX
    // var loopNum uintptr -> BX
    MOVQ 0x10(SP),CX
    // var sum uint64 -> AX
    XORQ AX,AX
    TESTQ CX,CX
    JE Finish_1
    ADDQ $0,AX
startFor_1:
    ADCL 0(BX),AX
    LEAQ 4(BX),BX
    DECQ CX
    JNE startFor_1
    ADCL $0,AX
Finish_1:
    MOVQ AX,0x18(SP)
    RET
