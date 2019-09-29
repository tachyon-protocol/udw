package udwStrings

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"strconv"
	"testing"
)

func TestGetIntersection(t *testing.T) {
	intersection := GetIntersection([]string{"1", "2"}, []string{"2", "5"}, []string{"12", "2", "1"})
	udwTest.Ok(IsEqualIgnoreOrder(intersection, []string{"2"}))

	intersection = GetIntersection([]string{"1", "2"}, []string{"2", "5"}, []string{"12", "0", "1"})
	udwTest.Ok(IsEqualIgnoreOrder(intersection, []string{}))

	intersection = GetIntersection([]string{}, []string{"2", "5"}, []string{"12", "0", "1"})
	udwTest.Ok(IsEqualIgnoreOrder(intersection, []string{}))

	intersection = GetIntersection([]string{"2", "5"})
	udwTest.Ok(IsEqualIgnoreOrder(intersection, []string{}))
}

func TestGetIntersectionV2(t *testing.T) {
	intersectionV2 := GetIntersectionV2([]string{"1", "2"}, []string{"2", "5"}, []string{"12", "2", "1"})
	udwTest.Ok(IsEqualIgnoreOrder(intersectionV2, []string{"2"}))

	intersectionV2 = GetIntersectionV2([]string{"1", "2"}, []string{"2", "5"}, []string{"12", "0", "1"})
	udwTest.Ok(IsEqualIgnoreOrder(intersectionV2, []string{}))

	intersectionV2 = GetIntersectionV2([]string{}, []string{"2", "5"}, []string{"12", "0", "1"})
	udwTest.Ok(IsEqualIgnoreOrder(intersectionV2, []string{}))

	intersectionV2 = GetIntersectionV2([]string{"2", "5"})
	udwTest.Ok(IsEqualIgnoreOrder(intersectionV2, []string{}))
}

func TestGetIntersectionCompare(t *testing.T) {
	sample := [][]string{}

	for i := 101; i > 0; i-- {
		l := []string{}
		for j := i; j > 0; j-- {
			l = append(l, strconv.Itoa(j))
		}
		sample = append(sample, l)
	}
	repeat := 1 << 10
	udwTest.BenchmarkWithRepeatNum(repeat, func() {
		GetIntersection(sample...)
	})
	udwTest.BenchmarkWithRepeatNum(repeat, func() {
		GetIntersectionV2(sample...)
	})
	udwTest.BenchmarkWithRepeatNum(repeat, func() {
		GetIntersectionV3(sample...)
	})
	udwTest.BenchmarkWithRepeatNum(repeat, func() {
		GetIntersectionV4(sample...)
	})
}

func TestIsEqualCheckOrder(t *testing.T) {
	udwTest.Ok(
		IsEqualCheckOrder([]string{"0", "1", "2"}, []string{"0", "1", "2"}, []string{"0", "1", "2"}),
	)
	udwTest.Ok(
		!IsEqualCheckOrder([]string{"1", "0", "2"}, []string{"0", "1", "2"}, []string{"0", "1", "2"}),
	)
	udwTest.Ok(
		!IsEqualCheckOrder([]string{"1", "1", "0", "2"}, []string{"0", "1", "2"}, []string{"0", "1", "2"}),
	)
	udwTest.Ok(
		IsEqualCheckOrder([]string{"1", "1", "0", "2"}),
	)
	udwTest.Ok(
		IsEqualCheckOrder([]string{"1", "1", "0", "2"}, []string{"1", "1", "0", "2"}),
	)
	udwTest.Ok(
		!IsEqualCheckOrder([]string{}, []string{"1", "1", "0", "2"}),
	)
}

func TestIsEqualIgnoreOrder(t *testing.T) {
	udwTest.Ok(
		IsEqualIgnoreOrder([]string{"0", "1", "2"}, []string{"0", "1", "2"}, []string{"0", "1", "2"}),
	)
	udwTest.Ok(
		IsEqualIgnoreOrder([]string{"0", "1", "2"}),
	)
	udwTest.Ok(
		IsEqualIgnoreOrder([]string{"0", "1", "2"}, []string{"1", "0", "2"}),
	)
	udwTest.Ok(
		!IsEqualIgnoreOrder([]string{}, []string{"1", "1", "0", "2"}),
	)
}
