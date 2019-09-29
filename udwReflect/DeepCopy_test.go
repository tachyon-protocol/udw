package udwReflect_test

import (
	"github.com/tachyon-protocol/udw/udwReflect"
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"testing"
)

func Test_DeepCopy(t *testing.T) {
	originS := "test"
	changeS := "change"
	type A struct {
		Data string
	}
	type Test1 struct {
		M map[string]*A
	}
	t1 := Test1{
		M: map[string]*A{
			originS: &A{
				Data: originS,
			},
		},
	}
	to1 := udwReflect.DeepCopy(t1).(Test1)
	udwTest.Equal(to1.M[originS].Data, originS)
	to1.M[originS].Data = changeS
	udwTest.Equal(to1.M[originS].Data, changeS)
	udwTest.Equal(t1.M[originS].Data, originS)

	t2 := &A{
		Data: originS,
	}
	to2 := udwReflect.DeepCopy(t2).(*A)
	udwTest.Equal(to2.Data, originS)
	to2.Data = changeS
	udwTest.Equal(to2.Data, changeS)
	udwTest.Equal(t2.Data, originS)
	t3 := []map[string]map[string]*A{
		map[string]map[string]*A{
			originS: map[string]*A{
				originS: &A{
					Data: originS,
				},
			},
		},
	}

	to3 := udwReflect.DeepCopy(t3).([]map[string]map[string]*A)
	udwTest.Equal(len(to3), 1)
	to3 = append(to3, map[string]map[string]*A{})
	udwTest.Equal(len(to3), 2)
	udwTest.Equal(len(t3), 1)
	udwTest.Equal(to3[0][originS][originS].Data, originS)
	to3[0][originS][originS].Data = changeS
	udwTest.Equal(to3[0][originS][originS].Data, changeS)
	udwTest.Equal(t3[0][originS][originS].Data, originS)

	t4 := map[string][]string{
		originS: []string{
			originS,
		},
	}
	to4 := udwReflect.DeepCopy(t4).(map[string][]string)
	udwTest.Equal(len(to4), 1)
	to4[changeS] = []string{}
	udwTest.Equal(len(to4), 2)
	udwTest.Equal(len(t4), 1)
	udwTest.Equal(len(to4[originS]), 1)
	to4[originS] = append(to4[originS], originS)
	udwTest.Equal(len(to4[originS]), 2)
	udwTest.Equal(len(t4[originS]), 1)
	udwTest.Equal(to4[originS][0], originS)
	to4[originS][0] = changeS
	udwTest.Equal(to4[originS][0], changeS)
	udwTest.Equal(t4[originS][0], originS)

	type Test5 struct {
		Lock sync.Mutex
		Data string
		gg   string
		Once sync.Once
	}
	t5 := &Test5{
		Data: originS,
		gg:   originS,
	}
	to5 := udwReflect.DeepCopy(t5).(*Test5)
	udwTest.Equal(to5.Data, originS)
	udwTest.Equal(to5.gg, "")

	type test6 struct {
		Data  string
		inner string
	}
	t6 := map[string]*test6{
		originS: &test6{
			Data:  originS,
			inner: originS,
		},
	}
	to6 := udwReflect.DeepCopy(t6).(map[string]*test6)
	udwTest.Equal(to6[originS].Data, originS)
	udwTest.Equal(to6[originS].inner, "")
	to6[originS].Data = changeS
	udwTest.Equal(to6[originS].Data, changeS)
	udwTest.Equal(t6[originS].Data, originS)
	udwTest.Equal(len(to6), 1)
	to6[changeS] = &test6{}
	udwTest.Equal(len(to6), 2)
	udwTest.Equal(len(t6), 1)

	t7 := map[string][]test6{
		originS: []test6{
			test6{
				Data: originS,
			},
			test6{
				Data: originS,
			},
		},
		"new": []test6{
			test6{
				Data: originS,
			},
		},
	}
	to7 := udwReflect.DeepCopy(t7).(map[string][]test6)
	udwTest.Equal(len(to7["new"]), 1)
	udwTest.Equal(len(to7), 2)
	to7[changeS] = []test6{}
	udwTest.Equal(len(to7), 3)
	udwTest.Equal(len(t7), 2)
	udwTest.Equal(len(to7[originS]), 2)
	to7[originS] = append(to7[originS], test6{})
	udwTest.Equal(len(to7[originS]), 3)
	udwTest.Equal(len(t7[originS]), 2)

	type test8 struct {
		Next *test8
		Id   string
	}
	t8 := test8{
		Id: originS,
		Next: &test8{
			Id: originS,
		},
	}
	to8 := udwReflect.DeepCopy(t8).(test8)
	udwTest.Equal(to8.Id, originS)
	to8.Id = changeS
	udwTest.Equal(to8.Id, changeS)
	udwTest.Equal(t8.Id, originS)
	udwTest.Equal(to8.Next.Id, originS)
	to8.Next.Id = changeS
	udwTest.Equal(to8.Next.Id, changeS)
	udwTest.Equal(t8.Next.Id, originS)

}
