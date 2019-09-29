package udwType

import (
	"fmt"
	"reflect"
	"strconv"
)

type SliceType struct {
	reflectTypeGetterImp
	getElemByStringEditorabler
	elemType UdwType
}

func (t *SliceType) init() (err error) {
	if t.elemType != nil {
		return
	}
	t.elemType, err = TypeOf(t.GetReflectType().Elem())
	return
}
func (t *SliceType) GetElemByString(v reflect.Value, k string) (ev reflect.Value, et UdwType, err error) {
	if err = t.init(); err != nil {
		return
	}
	et = t.elemType
	ev, err = arrayGetSubValueByString(v, k)
	if err != nil {
		return
	}
	return
}
func (t *SliceType) SaveByPath(v *reflect.Value, path Path, value string) (err error) {
	if err = t.init(); err != nil {
		return
	}
	if len(path) > 0 && path[0] == "" {
		if v.CanSet() {
			v.Set(
				reflect.Append(*v, reflect.New(t.elemType.GetReflectType()).Elem()),
			)
		} else {
			*v = reflect.Append(*v, reflect.New(t.elemType.GetReflectType()).Elem())
		}
		path[0] = strconv.Itoa(v.Len() - 1)
		if value == "" {
			return
		}
	}
	return t.getElemByStringEditorabler.SaveByPath(v, path, value)
}

func (t *SliceType) DeleteByPath(v *reflect.Value, path Path) (err error) {
	if err = t.init(); err != nil {
		return
	}
	if len(path) > 1 {
		return passThougthDeleteByPath(t, v, path)
	} else if len(path) == 0 {
		return fmt.Errorf("[SliceType.DeleteByPath] delete from slice with no path.")
	}
	i, err := arrayParseKey(*v, path[0])
	if err != nil {
		return err
	}
	if v.CanSet() {
		v.Set(
			reflect.AppendSlice(v.Slice(0, i), v.Slice(i+1, v.Len())),
		)
	} else {
		*v = reflect.AppendSlice(v.Slice(0, i), v.Slice(i+1, v.Len()))
	}
	return nil
}
