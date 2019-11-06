// +build windows

package udwWindowsRegistry

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"syscall"
	"unsafe"
)

func GetStringByPath(path string) (out string, err error) {
	item, err := getByPath(path)
	if err != nil {
		return "", err
	}
	return item.getString()
}

func MustGetStringByPath(path string) (out string) {
	out, err := GetStringByPath(path)
	if err != nil {
		panic(err)
	}
	return out
}

func GetDirectoryOrFileNameListOneLevel(path string) (nameList []string, err error) {
	key, errMsg := getKeyObjByPath(path, READ)
	if errMsg != "" {
		err = errors.New(errMsg)
		return
	}
	nameList, err = key.ReadSubKeyNames(-1)
	key.Close()
	if err != nil {
		err = fmt.Errorf("[GetDirectoryOrFileListOneLevel] ReadSubKeyNames fail path[%s] err[%s]", path, err)
	}
	return nameList, err
}

func MustGetDirectoryOrFileNameListOneLevel(path string) (nameList []string) {
	nameList, err := GetDirectoryOrFileNameListOneLevel(path)
	if err != nil {
		panic(err)
	}
	return nameList
}

type registryItem struct {
	data    []byte
	valtype uint32
}

func (item registryItem) getString() (string, error) {
	switch item.valtype {
	case SZ, EXPAND_SZ:
	default:
		return "", ErrUnexpectedType
	}
	u := (*[1 << 29]uint16)(unsafe.Pointer(&item.data[0]))[:]
	return syscall.UTF16ToString(u), nil
}

func (item registryItem) getUInt32() (uint32, error) {
	switch item.valtype {
	case DWORD:
	default:
		return 0, ErrUnexpectedType
	}
	return binary.LittleEndian.Uint32(item.data), nil
}

func getByPath(path string) (item registryItem, err error) {
	pathPartList := strings.Split(path, `\`)
	if len(pathPartList) <= 2 {
		return registryItem{}, fmt.Errorf("[getByPath] need at last 2 part of the path [%s]", path)
	}
	keyObj, errMsg := getKeyObjByPath(strings.Join(pathPartList[0:len(pathPartList)-1], `\`), READ)
	if errMsg != "" {
		err = errors.New(errMsg)
		return
	}
	defer keyObj.Close()
	data, typ, err := keyObj.getValue(pathPartList[len(pathPartList)-1], make([]byte, 64))
	if err != nil {
		err = fmt.Errorf("[GetByPath] getValue fail path [%s], err [%s]", path, err)
		return
	}
	return registryItem{
		data:    data,
		valtype: typ,
	}, nil
}

func getKeyObjByPath(path string, access uint32) (key Key, errMsg string) {
	pathPartList := strings.Split(path, `\`)
	if len(pathPartList) <= 1 {
		return 0, "[getKeyObjByPath] need at last 1 part of the path [" + path + "]"
	}
	firstKey, errMsg := getFirstKeyByString(pathPartList[0])
	if errMsg != "" {
		return 0, errMsg
	}
	toOpenPath := strings.Join(pathPartList[1:], `\`)
	key, err := OpenKey(firstKey, toOpenPath, access)
	if err != nil {
		return 0, "[getKeyObjByPath] OpenKey fail path [" + path + "], err [" + err.Error() + "]"
	}
	return key, ""
}

func getFirstKeyByString(firstKey string) (first Key, errMsg string) {
	switch firstKey {
	case "LOCAL_MACHINE", "HKEY_LOCAL_MACHINE", "HKLM":
		return LOCAL_MACHINE, ""
	case "CURRENT_CONFIG", "HKEY_CURRENT_CONFIG":
		return CURRENT_CONFIG, ""
	case "CLASSES_ROOT", "HKEY_CLASSES_ROOT":
		return CLASSES_ROOT, ""
	case "CURRENT_USER", "HKEY_CURRENT_USER":
		return CURRENT_USER, ""
	case "USERS", "HKEY_USERS":
		return USERS, ""
	default:
		return 0, "[getFirstKeyByString] unknown first part of the path [" + firstKey + "]"
	}
}

func MustSetStringByPath(path string, value string) {
	mustValueByPath(path, value)
}

func MustSetDWordByPath(path string, value uint32) {
	mustValueByPath(path, value)
}

func mustValueByPath(path string, value interface{}) {
	pathPartList := strings.Split(path, `\`)
	if len(pathPartList) <= 2 {
		panic("[MustSetStringByPath] need at last 2 part of the path [" + path + "]")
	}
	var keyObj Key
	var errMsg string
	var err error
	keyObj, errMsg = getKeyObjByPath(strings.Join(pathPartList[0:len(pathPartList)-1], `\`), READ|WRITE)
	if errMsg != "" {
		if ErrorMsgIsNotFound(errMsg) {
			firstKey, errMsg := getFirstKeyByString(pathPartList[0])
			if errMsg != "" {
				panic(errMsg)
			}
			keyObj, _, err = CreateKey(firstKey, strings.Join(pathPartList[1:len(pathPartList)-1], `\`), READ|WRITE)
			if err != nil {
				panic(err)
			}
		} else {
			panic(errMsg)
		}
	}
	defer keyObj.Close()
	setValueFunc := func(_part string) {
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.String:
			err = keyObj.SetStringValue(_part, fmt.Sprint(value))
			if err != nil {
				panic(err)
			}
		case reflect.Uint32:
			err = keyObj.SetDWordValue(_part, value.(uint32))
			if err != nil {
				panic(err)
			}
		default:
			panic("error type")
		}
	}
	lastPart := pathPartList[len(pathPartList)-1]
	data, typ, err := keyObj.getValue(lastPart, make([]byte, 64))
	if err != nil {
		errMsg = err.Error()
		if ErrorMsgIsNotFound(errMsg) {
			setValueFunc(lastPart)
			return
		} else {
			panic(errMsg)
		}

	}
	item := registryItem{
		data:    data,
		valtype: typ,
	}
	isExistSame := func() bool {
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.String:
			s, err := item.getString()
			if err != nil {
				panic(err)
			}
			return s == fmt.Sprint(value)
		case reflect.Uint32:
			s, err := item.getUInt32()
			if err != nil {
				panic(err)
			}
			return s == value.(uint32)
		default:
			panic("error type")
		}
	}
	if !isExistSame() {
		setValueFunc(lastPart)
	}
}

func MustDeleteByPath(path string) {
	pathPartList := strings.Split(path, `\`)
	if len(pathPartList) <= 2 {
		panic("[MustDeleteByPath] need at last 2 part of the path [" + path + "]")
	}
	firstKey, errMsg := getFirstKeyByString(pathPartList[0])
	if errMsg != "" {
		panic(errMsg)
	}
	err := DeleteKey(firstKey, strings.Join(pathPartList[1:], `\`))
	if err != nil {
		if IsErrorNotExist(err) {
			return
		}
		panic("[MustDeleteByPath] [" + path + "] " + err.Error())
	}
}
