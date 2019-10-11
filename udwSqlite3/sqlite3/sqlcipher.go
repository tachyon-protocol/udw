package sqlite3

import (
	"bytes"
	"errors"
	"os"
)

var zxrHeader = []byte("SQLite format 3\000")

func IsEncrypted(filename string) (bool, error) {

	db, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer db.Close()

	var header [16]byte
	n, err := db.Read(header[:])
	if err != nil {
		return false, err
	}
	if n != len(header) {
		return false, errors.New("go-sqlcipher: could not read full header")
	}

	encrypted := !bytes.Equal(header[:], zxrHeader)
	return encrypted, nil
}
