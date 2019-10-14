package udwDnsPacket

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestDomain(ot *testing.T) {
	udwTest.Equal(IsInDomainSuffixSet("", nil), false)
	udwTest.Equal(IsInDomainSuffixSet("abc", nil), false)
	udwTest.Equal(IsInDomainSuffixSet("abc", map[string]struct{}{}), false)
	udwTest.Equal(IsInDomainSuffixSet("abc", map[string]struct{}{"abc": struct{}{}}), true)
	udwTest.Equal(IsInDomainSuffixSet("abc.", map[string]struct{}{"abc": struct{}{}}), true)
	udwTest.Equal(IsInDomainSuffixSet("google.jp", map[string]struct{}{"jp": struct{}{}}), true)
	udwTest.Equal(IsInDomainSuffixSet("google.jp.", map[string]struct{}{"jp": struct{}{}}), true)
	udwTest.Equal(IsInDomainSuffixSet("bing.com", map[string]struct{}{"bing.com": struct{}{}}), true)
	udwTest.Equal(IsInDomainSuffixSet("google.com", map[string]struct{}{"bing.com": struct{}{}}), false)
}
