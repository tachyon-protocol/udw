package udwRand

import (
	"reflect"
	"strconv"
)

func PermStringSlice(slice []string) (output []string) {
	thisLen := len(slice)
	output = make([]string, thisLen)
	permSlice := getGlobalRand().Perm(thisLen)
	for i := 0; i < thisLen; i++ {
		output[i] = slice[permSlice[i]]
	}
	return
}

func PermStringSliceAndReturnPart(slice []string, partLen int) (output []string) {
	output = make([]string, partLen)
	permSlice := getGlobalRand().Perm(len(slice))
	for i := 0; i < partLen; i++ {
		output[i] = slice[permSlice[i]]
	}
	return
}

func MustPermInterfaceSlice(obj interface{}) (output interface{}) {
	inV := reflect.ValueOf(obj)
	if inV.Kind() != reflect.Slice {
		panic("[MustPermInterfaceSlice] inV.Kind()!=reflect.Slice")
	}
	thisLen := inV.Len()
	outputV := reflect.MakeSlice(inV.Type(), thisLen, thisLen)

	permSlice := getGlobalRand().Perm(thisLen)
	for i := 0; i < thisLen; i++ {
		outputV.Index(i).Set(inV.Index(permSlice[i]))
	}
	return outputV.Interface()
}

func MustPermInterfaceSliceInPlace(obj interface{}) {
	inV := reflect.ValueOf(obj)
	if inV.Kind() != reflect.Slice {
		panic("[MustPermInterfaceSlice] inV.Kind()!=reflect.Slice")
	}
	thisLen := inV.Len()
	if thisLen <= 1 {
		return
	}
	inVElem := inV.Type().Elem()
	for i := 0; i < thisLen; i++ {
		j := getGlobalRand().Intn(thisLen)
		if i == j {
			continue
		}
		tmp := reflect.New(inVElem).Elem()
		tmp.Set(inV.Index(j))
		inV.Index(j).Set(inV.Index(i))
		inV.Index(i).Set(tmp)
	}
	return
}

func Perm(num int) (out []int) {
	return getGlobalRand().Perm(num)
}

func (r *UdwRand) Perm(n int) []int {
	r.locker.Lock()
	o := r.Rand.Perm(n)
	r.locker.Unlock()
	return o
}

func (r *UdwRand) PermIntSlice(slice []int) (output []int) {
	thisLen := len(slice)
	output = make([]int, thisLen)
	permSlice := r.Perm(thisLen)
	for i := 0; i < thisLen; i++ {
		output[i] = slice[permSlice[i]]
	}
	return
}

func (r *UdwRand) PermNoAllocNoLock(randomIndexArray []int) {
	if len(randomIndexArray) <= 1 {
		return
	}
	randomIndexArray[0] = 0
	for i := 1; i < len(randomIndexArray); i++ {
		j := r.Rand.Intn(i + 1)
		randomIndexArray[i] = randomIndexArray[j]
		randomIndexArray[j] = i
	}
}

func (r *UdwRand) ShuffleIntArrayNoAllocNoLock(arr []int, ShuffleNum int) {
	dataSize := len(arr)
	if ShuffleNum >= dataSize {
		panic("ShuffleIntArray fail " + strconv.Itoa(ShuffleNum) + " " + strconv.Itoa(dataSize))
	}
	for i := 0; i < ShuffleNum; i++ {
		j := int(r.Rand.Uint32()%uint32(dataSize-i)) + i
		tmp := arr[i]
		arr[i] = arr[j]
		arr[j] = tmp
	}
}

func PermNoAlloc(randomIndexArray []int) {
	r := getGlobalRand()
	r.locker.Lock()
	r.PermNoAllocNoLock(randomIndexArray)
	r.locker.Unlock()
}

func PermFromCacheNoAlloc(cache []int, start int, end int, output []int) {
	i := 0
	for _, index := range cache {
		_index := index + start
		if _index > end {
			continue
		}
		output[i] = _index
		i++
		if i == len(output) {
			return
		}
	}
}
