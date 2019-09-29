package udwDebug_test

import (
	"github.com/tachyon-protocol/udw/udwDebug"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestSprintln(ot *testing.T) {
	udwTest.Equal(udwDebug.Sprintln([]byte{0, 1}), "[udwDebug.Println] []byte{0x0, 0x1}\n")
	udwTest.Equal(udwDebug.Sprintln([]byte{}), "[udwDebug.Println] []byte{}\n")

	udwTest.Equal(udwDebug.Sprintln([]byte(nil)), "[udwDebug.Println] nil\n")
}
