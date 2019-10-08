package udwType

import (
	"fmt"
	"reflect"
	"strconv"
)

func arrayGetSubValueByString(v reflect.Value, k string) (reflect.Value, error) {
	i, err := arrayParseKey(v, k)
	if err != nil {
		return reflect.Value{}, nil
	}
	return v.Index(i), nil
}
func arrayParseKey(v reflect.Value, k string) (int, error) {
	i64, err := strconv.ParseInt(k, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("[arrayParseKey] index is not int k:%s", k)
	}
	i := int(i64)
	if i >= v.Len() || i < 0 {
		return 0, fmt.Errorf("[arrayParseKey] index is not of range k:%s,len:%d", k, v.Len())
	}
	return i, nil
}

type saveScaleFromStringer struct {
	SaveScaleInterface
	ReflectTypeGetter
}

func (t saveScaleFromStringer) FromString(value string) (reflect.Value, error) {
	rv := reflect.New(t.GetReflectType()).Elem()
	err := t.SaveScale(rv, value)
	if err != nil {
		return reflect.Value{}, err
	}
	return rv, nil
}

type reflectTypeGetterImp struct {
	reflect.Type
}

func (t reflectTypeGetterImp) GetReflectType() reflect.Type {
	return t.Type
}

type saveScaleEditabler struct {
	SaveScaleInterface
	ReflectTypeGetter
}

func (t saveScaleEditabler) SaveByPath(v *reflect.Value, path Path, value string) (err error) {
	if len(path) != 0 {
		return fmt.Errorf("[saveScaleEditabler.Save] get string with some path,path:%s type:%s", path, v.Type().Kind())
	}
	if !v.CanSet() {
		*v = reflect.New(t.GetReflectType()).Elem()
	}
	return t.SaveScale(*v, value)
}

func (t saveScaleEditabler) DeleteByPath(v *reflect.Value, path Path) (err error) {
	return fmt.Errorf("[saveScaleEditabler.Delete] scale type can not delete,path:%s type:%s", path, v.Type().Kind())
}

type getElemByStringEditorabler struct {
	GetElemByStringInterface
	ReflectTypeGetter
}

func (t getElemByStringEditorabler) SaveByPath(v *reflect.Value, path Path, value string) (err error) {
	if len(path) == 0 {
		return nil
	}
	ev, et, err := t.GetElemByString(*v, path[0])
	if err != nil {
		return err
	}
	oEv := ev
	err = et.SaveByPath(&oEv, path[1:], value)
	if err != nil {
		return err
	}

	if oEv == ev {
		return nil
	}
	if v.CanSet() {
		return nil
	}
	output := reflect.New(t.GetReflectType()).Elem()
	output.Set(*v)
	*v = output
	ev, _, err = t.GetElemByString(*v, path[0])
	if err != nil {
		return err
	}
	ev.Set(oEv)
	return nil
}

func passThougthDeleteByPath(t GetElemByStringAndReflectTypeGetterInterface, v *reflect.Value, path Path) (err error) {
	ev, et, err := t.GetElemByString(*v, path[0])
	if err != nil {
		return err
	}
	oEv := ev
	err = et.DeleteByPath(&oEv, path[1:])
	if err != nil {
		return err
	}
	if oEv == ev {
		return
	}
	if v.CanSet() {
		return
	}
	output := reflect.New(t.GetReflectType()).Elem()
	output.Set(*v)
	*v = output
	ev, _, err = t.GetElemByString(*v, path[0])
	if err != nil {
		return err
	}
	ev.Set(oEv)
	return nil
}
