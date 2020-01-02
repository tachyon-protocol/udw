// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
#include "textflag.h"
DATA bswapMask<>+0x00(SB)/8, $0x08090a0b0c0d0e0f
DATA bswapMask<>+0x08(SB)/8, $0x0001020304050607
GLOBL bswapMask<>(SB), (NOPTR+RODATA), $16
// func fillEightBlocks(nr int, xk *uint32, dst, counter *byte)
TEXT ·fillEightBlocks(SB),0,$112-32
#define BSWAP X2
#define aesRound AESENC X1, X8; AESENC X1, X9; AESENC X1, X10; AESENC X1, X11; \
                 AESENC X1, X12; AESENC X1, X13; AESENC X1, X14; AESENC X1, X15;
#define increment(i) ADDQ $1, R9; ADCQ $0, R8; \
                     MOVQ R9, (i*16)(SP); MOVQ R8, (i*16+8)(SP);
	MOVQ nr+0(FP), CX
	MOVQ xk+8(FP), AX
	MOVQ dst+16(FP), DX
	MOVQ counter+24(FP), BX
	MOVOU 0(AX), X1
	MOVOU bswapMask<>(SB), BSWAP
	ADDQ $16, AX
	MOVOU 0(BX), X8
	MOVQ 0(BX), R8
	MOVQ 8(BX), R9
	BSWAPQ R8
	BSWAPQ R9
	increment(0)
	increment(1)
	increment(2)
	increment(3)
	increment(4)
	increment(5)
	increment(6)
	ADDQ $1, R9
	ADCQ $0, R8
	BSWAPQ R8
	BSWAPQ R9
	MOVQ R8, 0(BX)
	MOVQ R9, 8(BX)
	MOVOU 0(SP), X9
	MOVOU 16(SP), X10
	MOVOU 32(SP), X11
	MOVOU 48(SP), X12
	MOVOU 64(SP), X13
	MOVOU 80(SP), X14
	MOVOU 96(SP), X15
	PSHUFB BSWAP, X9
	PSHUFB BSWAP, X10
	PSHUFB BSWAP, X11
	PSHUFB BSWAP, X12
	PSHUFB BSWAP, X13
	PSHUFB BSWAP, X14
	PSHUFB BSWAP, X15
	PXOR X1, X8
	PXOR X1, X9
	PXOR X1, X10
	PXOR X1, X11
	PXOR X1, X12
	PXOR X1, X13
	PXOR X1, X14
	PXOR X1, X15
	SUBQ $12, CX
	JE Lenc196
	JB Lenc128
Lenc256:
	MOVOU 0(AX), X1
	aesRound
	MOVOU 16(AX), X1
	aesRound
	ADDQ $32, AX
Lenc196:
	MOVOU 0(AX), X1
	aesRound
	MOVOU 16(AX), X1
	aesRound
	ADDQ $32, AX
Lenc128:
	MOVOU 0(AX), X1
	aesRound
	MOVOU 16(AX), X1
	aesRound
	MOVOU 32(AX), X1
	aesRound
	MOVOU 48(AX), X1
	aesRound
	MOVOU 64(AX), X1
	aesRound
	MOVOU 80(AX), X1
	aesRound
	MOVOU 96(AX), X1
	aesRound
	MOVOU 112(AX), X1
	aesRound
	MOVOU 128(AX), X1
	aesRound
	MOVOU 144(AX), X1
	AESENCLAST X1, X8
	AESENCLAST X1, X9
	AESENCLAST X1, X10
	AESENCLAST X1, X11
	AESENCLAST X1, X12
	AESENCLAST X1, X13
	AESENCLAST X1, X14
	AESENCLAST X1, X15
	MOVOU X8, 0(DX)
	MOVOU X9, 16(DX)
	MOVOU X10, 32(DX)
	MOVOU X11, 48(DX)
	MOVOU X12, 64(DX)
	MOVOU X13, 80(DX)
	MOVOU X14, 96(DX)
	MOVOU X15, 112(DX)
	RET
// func xorBytes(dst, a, b []byte) int
TEXT ·xorBytes(SB),NOSPLIT,$0
	MOVQ dst_base+0(FP), DI
	MOVQ a_base+24(FP), SI
	MOVQ a_len+32(FP), R8
	MOVQ b_base+48(FP), BX
	MOVQ b_len+56(FP), R9
	CMPQ R8, R9
	JLE skip
	MOVQ R9, R8
skip:
	MOVQ R8, ret+72(FP)
	XORQ CX, CX
	CMPQ R8, $16
	JL tail
loop:
	MOVOU (SI)(CX*1), X1
	MOVOU (BX)(CX*1), X2
	PXOR X1, X2
	MOVOU X2, (DI)(CX*1)
	ADDQ $16, CX
	SUBQ $16, R8
	CMPQ R8, $16
	JGE loop
tail:
	CMPQ R8, $0
	JE done
	MOVBLZX (SI)(CX*1), R9
	MOVBLZX (BX)(CX*1), R10
	XORL R10, R9
	MOVB R9B, (DI)(CX*1)
	INCQ CX
	DECQ R8
	JMP tail
done:
	RET
