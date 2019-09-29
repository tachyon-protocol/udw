package udwDebug_test

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwDebug"
	"github.com/tachyon-protocol/udw/udwTest"
	"strings"
	"testing"
)

func TestGetCurrentAllStackString(t *testing.T) {
	s := testFnNamemfr54rmg5s()
	fmt.Println(s)
	udwTest.Ok(strings.Contains(s, "testFnNamemfr54rmg5s"))
}

func testFnNamemfr54rmg5s() string {
	return udwDebug.GetCurrentAllStackString(0)
}
