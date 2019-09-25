package udwRand

import (
	"fmt"
	mrand "math/rand"
	"sort"
	"sync"
	"time"
)

func NewCryptSeedMathRand() (r *mrand.Rand) {
	seed := int64(MustCryptoRandUint64())
	return mrand.New(mrand.NewSource(seed))
}

func MustNewCryptSeedUdwRand() (r *UdwRand) {
	mr := NewCryptSeedMathRand()
	return &UdwRand{Rand: mr}
}

func NewInt64SeedUdwRand(seed int64) (r *UdwRand) {
	mr := mrand.New(mrand.NewSource(seed))
	return &UdwRand{Rand: mr}
}

type UdwRand struct {
	*mrand.Rand
	locker sync.Mutex
}

func (r *UdwRand) MulitChoice(totalLength int, choiceNumber int) []int {
	if choiceNumber > totalLength || totalLength < 0 || choiceNumber < 0 {
		panic(fmt.Errorf("[UdwRand.MulitChoice] input error totalLength=%d choiceNumber=%d,need (choiceNumber<=totalLength&&totalLength>=0&&choiceNumber>=0) ",
			totalLength, choiceNumber))
	}
	return r.Perm(totalLength)[:choiceNumber]
}

func (r *UdwRand) MulitChoiceOriginOrder(n int, size int) []int {
	perm := r.MulitChoice(n, size)
	sort.IntSlice(perm).Sort()
	return perm
}

func (r *UdwRand) HappendBaseOnPossibility(possibility float64) bool {
	if possibility > 1+1e-10 || possibility < -1e-10 {
		panic(fmt.Errorf("[UdwRand.HappendBaseOnPossibility] possibility:%f > 1 or < 0", possibility))
	}
	out := r.Float64()
	ret := out < possibility
	return ret
}

func (r *UdwRand) IntBetween(min int, max int) int {
	if min > max {
		panic(fmt.Errorf("[UdwRand.IntBetween] min:%d<max:%d", min, max))
	} else if min == max {
		return min
	}
	r.locker.Lock()
	defer r.locker.Unlock()
	o := r.Rand.Intn(max-min+1) + min

	return o
}

func (r *UdwRand) Int63Between(min int64, max int64) int64 {
	if min > max {
		panic(fmt.Errorf("[UdwRand.Int63Between] min:%d<max:%d", min, max))
	} else if min == max {
		return min
	}
	r.locker.Lock()
	defer r.locker.Unlock()
	o := r.Int63n(max-min) + min
	return o
}

func (r *UdwRand) Float64Between(min float64, max float64) float64 {
	if min > max {
		panic(fmt.Errorf("[UdwRand.Float64Between] min:%f<max:%f", min, max))
	} else if min == max {
		return min
	}
	return r.Float64()*(max-min) + min
}

func (r *UdwRand) TimeDurationBetween(min time.Duration, max time.Duration) time.Duration {
	return time.Duration(r.Int63Between(int64(min), int64(max)))
}

func (r *UdwRand) ChoiceFromIntSlice(slice []int) int {
	return slice[r.Intn(len(slice))]
}

func (r *UdwRand) Intn(n int) int {
	r.locker.Lock()
	o := r.Rand.Intn(n)
	r.locker.Unlock()
	return o
}

func (r *UdwRand) Float32() float32 {
	r.locker.Lock()
	o := r.Rand.Float32()
	r.locker.Unlock()
	return o
}

func (r *UdwRand) Float64() float64 {
	r.locker.Lock()
	o := r.Rand.Float64()
	r.locker.Unlock()
	return o
}

var globalRand *UdwRand
var globalRandOnce sync.Once

func getGlobalRand() *UdwRand {
	globalRandOnce.Do(func() {
		globalRand = MustNewCryptSeedUdwRand()
	})
	return globalRand
}

func IntBetween(min int, max int) int {
	return getGlobalRand().IntBetween(min, max)
}

func Int63Between(min int64, max int64) int64 {
	return getGlobalRand().Int63Between(min, max)
}

func IntSliceBetween(start int, end int, size int) []int {
	if size <= 0 {
		return nil
	}
	n := end - start + 1
	if n < 0 {
		return nil
	}
	if n == 0 {
		if size == 1 {
			return []int{start}
		}
		return nil
	}
	nl := Perm(n)
	_nl := make([]int, 0, n)
	for _, i := range nl {
		_nl = append(_nl, i+start)
	}
	if size >= n {
		return _nl
	}
	return _nl[:size]
}

func TimeDurationBetween(min time.Duration, max time.Duration) time.Duration {
	return getGlobalRand().TimeDurationBetween(min, max)
}

func HappendBaseOnPossibility(possibility float64) bool {
	return getGlobalRand().HappendBaseOnPossibility(possibility)
}

func Intn(num int) (out int) {
	return getGlobalRand().Intn(num)
}

func Float32() float32 {
	return getGlobalRand().Float32()
}

func Float64() float64 {
	return getGlobalRand().Float64()
}

func Float64Between(min float64, max float64) float64 {
	return getGlobalRand().Float64Between(min, max)
}

func RandRunOneFromFuncList(funList []func()) {
	if len(funList) == 0 {
		return
	}
	i := Intn(len(funList))
	funList[i]()
}
