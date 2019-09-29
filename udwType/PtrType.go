package udwType

import (
	"fmt"
	"reflect"
)

type PtrType struct {
	reflectTypeGetterImp
	elemType UdwType
}

func (t *PtrType) init() (err error) {
	if t.elemType != nil {
		return
	}
	t.elemType, err = TypeOf(t.GetReflectType().Elem())
	return
}
func (t *PtrType) SaveByPath(inV *reflect.Value, path Path, value string) error {
	err := t.init()
	if err != nil {
		return err
	}

	if inV.IsNil() {
		if inV.CanSet() {
			inV.Set(reflect.New(t.GetReflectType().Elem()))
		} else {
			*inV = reflect.New(t.GetReflectType().Elem())
		}
	}

	elemV := inV.Elem()
	if len(path) >= 1 {
		return t.elemType.SaveByPath(&elemV, path[1:], value)
	}
	return nil
}
func (t *PtrType) GetElemByString(v reflect.Value, k string) (ev reflect.Value, et UdwType, err error) {
	err = t.init()
	if err != nil {
		return
	}
	if v.IsNil() {
		err = fmt.Errorf("[PtrType.GetElemByString] get null pointer k:%s", k)
		return
	}
	ev = v.Elem()
	et = t.elemType
	return
}
func (t *PtrType) DeleteByPath(v *reflect.Value, path Path) (err error) {
	if len(path) > 1 {
		return passThougthDeleteByPath(t, v, path)
	} else if len(path) == 0 {
		return fmt.Errorf("[PtrType.DeleteByPath] delete ptr with no path.")
	}
	nilPtr := reflect.Zero(t.GetReflectType())
	if v.CanSet() {
		v.Set(nilPtr)
	} else {
		*v = nilPtr
	}
	return nil
}
