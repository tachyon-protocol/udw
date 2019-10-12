package udwCache

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwPlatform"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
	"time"
)

func TestFileTtlCache(t *testing.T) {
	if udwPlatform.IsWindows() {

		return
	}
	udwFile.MustDelete(getFileTtlCachePath("test_file_ttl_cache"))
	_, err := FileTtlCache("test_file_ttl_cache", func() (b []byte, ttl time.Duration, err error) {
		return []byte("1"), time.Millisecond, fmt.Errorf("error")
	})
	udwTest.Equal(err.Error(), "error")
	b, err := FileTtlCache("test_file_ttl_cache", func() (b []byte, ttl time.Duration, err error) {
		return []byte("1"), time.Millisecond, nil
	})
	udwTest.Equal(b, []byte("1"))
	udwTest.Equal(err, nil)

	b, err = FileTtlCache("test_file_ttl_cache", func() (b []byte, ttl time.Duration, err error) {
		return []byte("2"), time.Millisecond, nil
	})
	udwTest.Equal(b, []byte("1"))
	udwTest.Equal(err, nil)

	time.Sleep(time.Millisecond)
	b, err = FileTtlCache("test_file_ttl_cache", func() (b []byte, ttl time.Duration, err error) {
		return []byte("2"), time.Millisecond, nil
	})
	udwTest.Equal(b, []byte("2"))
	udwTest.Equal(err, nil)
}
