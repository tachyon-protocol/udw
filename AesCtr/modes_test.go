package AesCtr

import (
	"crypto/cipher"
	"testing"
)

type testInterface interface {
	InAESPackage() bool
}

type testBlock struct{}

func (*testBlock) BlockSize() int      { return 0 }
func (*testBlock) Encrypt(a, b []byte) {}
func (*testBlock) Decrypt(a, b []byte) {}
func (*testBlock) NewGCM(int) (cipher.AEAD, error) {
	return &testAEAD{}, nil
}
func (*testBlock) NewCBCEncrypter([]byte) cipher.BlockMode {
	return &testBlockMode{}
}
func (*testBlock) NewCBCDecrypter([]byte) cipher.BlockMode {
	return &testBlockMode{}
}
func (*testBlock) NewCTR([]byte) cipher.Stream {
	return &testStream{}
}

type testAEAD struct{}

func (*testAEAD) NonceSize() int                         { return 0 }
func (*testAEAD) Overhead() int                          { return 0 }
func (*testAEAD) Seal(a, b, c, d []byte) []byte          { return []byte{} }
func (*testAEAD) Open(a, b, c, d []byte) ([]byte, error) { return []byte{}, nil }
func (*testAEAD) InAESPackage() bool                     { return true }

func TestGCMAble(t *testing.T) {
	b := cipher.Block(&testBlock{})
	if _, ok := b.(gcmAble); !ok {
		t.Fatalf("testBlock does not implement the gcmAble interface")
	}
	aead, err := cipher.NewGCM(b)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if _, ok := aead.(testInterface); !ok {
		t.Fatalf("cipher.NewGCM did not use gcmAble interface")
	}
}

type testBlockMode struct{}

func (*testBlockMode) BlockSize() int          { return 0 }
func (*testBlockMode) CryptBlocks(a, b []byte) {}
func (*testBlockMode) InAESPackage() bool      { return true }

func TestCBCEncAble(t *testing.T) {
	b := cipher.Block(&testBlock{})
	if _, ok := b.(cbcEncAble); !ok {
		t.Fatalf("testBlock does not implement the cbcEncAble interface")
	}
	bm := cipher.NewCBCEncrypter(b, []byte{})
	if _, ok := bm.(testInterface); !ok {
		t.Fatalf("cipher.NewCBCEncrypter did not use cbcEncAble interface")
	}
}

func TestCBCDecAble(t *testing.T) {
	b := cipher.Block(&testBlock{})
	if _, ok := b.(cbcDecAble); !ok {
		t.Fatalf("testBlock does not implement the cbcDecAble interface")
	}
	bm := cipher.NewCBCDecrypter(b, []byte{})
	if _, ok := bm.(testInterface); !ok {
		t.Fatalf("cipher.NewCBCDecrypter did not use cbcDecAble interface")
	}
}

type testStream struct{}

func (*testStream) XORKeyStream(a, b []byte) {}
func (*testStream) InAESPackage() bool       { return true }

func TestCTRAble(t *testing.T) {
	b := cipher.Block(&testBlock{})
	if _, ok := b.(ctrAble); !ok {
		t.Fatalf("testBlock does not implement the ctrAble interface")
	}
	s := cipher.NewCTR(b, []byte{})
	if _, ok := s.(testInterface); !ok {
		t.Fatalf("cipher.NewCTR did not use ctrAble interface")
	}
}
