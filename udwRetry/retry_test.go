package udwRetry

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestRetry(t *testing.T) {
	debug := true
	i := -10
	err := Run(2, func() {
		i++
		if i == 1 {
			return
		}
		panic(i)
	})
	if debug {
		fmt.Println(err)
		return
	}
	udwErr.PanicIfError(err)
	udwTest.Equal(i, 1)
}

func TestInterval(t *testing.T) {
	return

}
