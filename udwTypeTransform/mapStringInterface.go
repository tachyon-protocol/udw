package udwTypeTransform

import (
	"fmt"
	"reflect"
	"strconv"
)

func GetStringInPlistObj(objI interface{}, keyList ...string) string {
	for _, key := range keyList {
		if objI == nil {
			return ""
		}
		objM, ok := objI.(map[string]interface{})
		if !ok {
			return ""
		}
		objI = objM[key]
	}
	switch objS := objI.(type) {
	case string:
		return objS
	case uint64:
		return strconv.FormatUint(objS, 10)
	default:
		return ""
	}
}

func SetStringInPlistObj(objI interface{}, keyList []string, toSetValue string) bool {
	var m map[string]interface{}
	var ok bool
	if len(keyList) == 0 {
		panic("[SetStringInPlistObj] len(keyList)==0")
	}
	if len(keyList) >= 2 {
		for _, key := range keyList[:len(keyList)-1] {
			if objI == nil {
				panic("[SetStringInPlistObj] objI==nil,key[" + key + "]")
			}
			objM, ok := objI.(map[string]interface{})
			if !ok {
				panic(fmt.Sprintf("[SetStringInPlistObj]1 objI[%T] is not map[string]interface{},key[%s]", objI, key))
			}
			objI = objM[key]
			if objI == nil {
				objI = map[string]interface{}{}
				objM[key] = objI
			}
		}
	}
	m, ok = objI.(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("[SetStringInPlistObj]2 objI[%T] is not map[string]interface{}", objI))
	}
	oldValue, ok := m[keyList[len(keyList)-1]]
	if !ok {
		m[keyList[len(keyList)-1]] = toSetValue
		return true
	}
	oldValueS, ok := oldValue.(string)
	if !ok {
		panic(fmt.Sprintf("[SetStringInPlistObj]4 oldValue[%T] is not string,key[%s]", oldValue, keyList[len(keyList)-1]))
	}
	if oldValueS == toSetValue {
		return false
	}
	m[keyList[len(keyList)-1]] = toSetValue
	return true
}

func SetStringListInPlistObj(objI interface{}, keyList []string, toSetValue []string) bool {
	var m map[string]interface{}
	var ok bool
	if len(keyList) == 0 {
		panic("[SetStringInPlistObj] len(keyList)==0")
	}
	if len(keyList) >= 2 {
		for _, key := range keyList[:len(keyList)-1] {
			if objI == nil {
				panic("[SetStringInPlistObj] objI==nil,key[" + key + "]")
			}
			objM, ok := objI.(map[string]interface{})
			if !ok {
				panic(fmt.Sprintf("[SetStringInPlistObj]1 objI[%T] is not map[string]interface{},key[%s]", objI, key))
			}
			objI = objM[key]
			if objI == nil {
				objI = map[string]interface{}{}
				objM[key] = objI
			}
		}
	}
	m, ok = objI.(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("[SetStringInPlistObj]2 objI[%T] is not map[string]interface{}", objI))
	}
	toSetValueObjList := make([]interface{}, len(toSetValue))
	for i := range toSetValue {
		toSetValueObjList[i] = toSetValue[i]
	}
	oldValue, ok := m[keyList[len(keyList)-1]]
	if !ok {
		m[keyList[len(keyList)-1]] = toSetValueObjList
		return true
	}
	oldValueS, ok := oldValue.([]interface{})
	if !ok {
		m[keyList[len(keyList)-1]] = toSetValueObjList
		return true
	}

	if reflect.DeepEqual(oldValueS, toSetValueObjList) {
		return false
	}
	m[keyList[len(keyList)-1]] = toSetValueObjList
	return true
}

func DeleteStringInPlistObj(objI interface{}, keyList ...string) bool {
	var ok bool
	if len(keyList) == 0 {
		panic("[DeleteStringInPlistObj] len(keyList)==0")
	}
	if len(keyList) >= 2 {
		for _, key := range keyList[:len(keyList)-1] {
			if objI == nil {
				return false
			}
			objM, ok := objI.(map[string]interface{})
			if !ok {
				panic(fmt.Sprintf("[SetStringInPlistObj]1 objI[%T] is not map[string]interface{},key[%s]", objI, key))
			}
			objI = objM[key]
		}
	}
	if objI == nil {
		return false
	}
	m, ok := objI.(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("[DeleteStringInPlistObj]2 objI[%T] is not map[string]interface{}", objI))
	}
	_, ok = m[keyList[len(keyList)-1]]
	if !ok {
		return false
	}
	delete(m, keyList[len(keyList)-1])
	return true
}

func GetInterfaceInPlistObj(objI interface{}, keyList ...string) interface{} {
	for _, key := range keyList {
		if objI == nil {
			return ""
		}
		objM, ok := objI.(map[string]interface{})
		if !ok {
			return ""
		}
		objI = objM[key]
	}
	return objI
}

func GetStringListInPlistObj(objI interface{}, keyList ...string) []string {
	for _, key := range keyList {
		if objI == nil {
			return nil
		}
		objM, ok := objI.(map[string]interface{})
		if !ok {
			return nil
		}
		objI = objM[key]
	}
	objS, ok := objI.([]interface{})
	if !ok {
		return nil
	}
	outList := []string{}
	for _, objI := range objS {
		thisS, ok := objI.(string)
		if !ok {
			return nil
		}
		outList = append(outList, thisS)
	}
	return outList
}

func GetMapInPlistObj(objI interface{}, keyList ...string) map[string]interface{} {
	for _, key := range keyList {
		if objI == nil {
			return nil
		}
		objM, ok := objI.(map[string]interface{})
		if !ok {
			return nil
		}
		objI = objM[key]
	}
	objS, ok := objI.(map[string]interface{})
	if !ok {
		return nil
	}
	return objS
}
