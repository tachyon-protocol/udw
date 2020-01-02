package udwQueryOnlyUrlFuzzLib

import (
	"github.com/tachyon-protocol/udw/udwQueryOnlyUrl"
)

func Fuzz(data []byte) int {
	s := string(data)
	obj := udwQueryOnlyUrl.ParseQueryUrlObj(s)
	if obj == nil {
		return 0
	}
	s2 := obj.Marshal()
	canReMarhsal(obj)
	obj.AddKv("a", "b")
	s4 := obj.Marshal()
	if s4 != s2+"&a=b" {
		panic("3")
	}
	canReMarhsal(obj)
	obj.AddKv(s, s)
	s5 := obj.Marshal()
	shouldS5 := s4 + "&" + udwQueryOnlyUrl.UrlvEncode(s) + "=" + udwQueryOnlyUrl.UrlvEncode(s)
	if s5 != shouldS5 {
		panic("4")
	}
	canReMarhsal(obj)
	return 1
}

func canReMarhsal(obj *udwQueryOnlyUrl.QueryUrlObj) {
	s2 := obj.Marshal()
	obj2 := udwQueryOnlyUrl.ParseQueryUrlObj(s2)
	if obj2 == nil {
		panic("obj2==nil")
	}
	s3 := obj2.Marshal()
	if s2 != s3 {
		panic("s2!=s3")
	}
}
