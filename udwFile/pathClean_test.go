package udwFile

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestPathClean(ot *testing.T) {
	for _, cas := range []struct {
		in  string
		out string
	}{
		{"/1.txt", "1.txt"},
		{"/1.txt/", "1.txt"},
		{"//1.txt////", "1.txt"},
		{"test/1.txt", "test/1.txt"},
		{"/../1.txt", "1.txt"},
		{"../1.txt", "1.txt"},
		{"1.txt/..", "1.txt"},
		{"1.txt/.", "1.txt"},
		{"..", ""},
		{".", ""},
		{"/test/../1.txt", "test/1.txt"},
		{"/test/..//../.../1.txt/", "test/1.txt"},
		{"/test/..//../test.../1.txt/../../", "test/test.../1.txt"},
		{"../abc/..", "abc"},
		{"../abc/../////", "abc"},
	} {
		udwTest.Equal(PathClean(cas.in), cas.out, cas.in)
	}
}
