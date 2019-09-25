package udwRand

import "fmt"

type LcgTransformer struct {
	Start uint64
	Range uint64
	A     uint64
	C     uint64
}

func (t LcgTransformer) GenerateInRange(i uint64) (output uint64) {
	if i >= t.Range {
		panic(fmt.Errorf("[LcgTransformer.GenerateInRange] i[%d]>=t.Range[%d]", i, t.Range))
	}
	return t.Start + (i*t.A+t.C)%t.Range
}

func (t LcgTransformer) Generate(i uint64) (output uint64) {
	return t.Start + (i*t.A+t.C)%t.Range
}
