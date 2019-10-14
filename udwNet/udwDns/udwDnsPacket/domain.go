package udwDnsPacket

import (
	"github.com/tachyon-protocol/udw/udwStrings"
	"strings"
)

func NormalizeDomain(domain string) string {
	return udwStrings.TrimOneChar(domain, '.')
}

func IsInDomainSuffixSet(domain string, m map[string]struct{}) bool {
	if len(m) == 0 || len(domain) == 0 {
		return false
	}
	domain = strings.TrimRight(domain, ".")
	thisSearchDomain := domain
	for {
		thisPos := strings.LastIndexByte(thisSearchDomain, '.')
		var thisPart string
		if thisPos == -1 {
			thisPart = domain
		} else {
			thisPart = domain[thisPos+1:]
		}
		_, ok := m[thisPart]
		if ok {
			return true
		}
		if thisPos <= 0 {
			break
		}
		thisSearchDomain = domain[:thisPos-1]
	}
	return false
}
