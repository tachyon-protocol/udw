package udwCryptoAesCtrV7

func Encrypt16NoRandom(psk *[32]byte, in []byte) (out []byte) {
	return EncryptRecord(RecordRequest{
		In:             in,
		Psk:            psk,
		RandomOverHead: 0,
		VerifyOverHead: 16,
	})
}

func Decrypt16NoRandom(psk *[32]byte, in []byte) (out []byte, errMsg string) {
	return DecryptRecord(RecordRequest{
		In:             in,
		Psk:            psk,
		RandomOverHead: 0,
		VerifyOverHead: 16,
	})
}
