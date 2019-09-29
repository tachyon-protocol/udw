package udwType

import (
	"fmt"
	"reflect"
)

type MapType struct {
	reflectTypeGetterImp
	KeyStringConverter StringConverterInterface
	KeyType            UdwType
	ElemType           UdwType
}

func (t *MapType) Init() (err error) {
	if t.KeyStringConverter != nil {
		return
	}
	t.KeyType, err = TypeOf(t.GetReflectType().Key())
	if err != nil {
		return err
	}
	var ok bool
	t.KeyStringConverter, ok = t.KeyType.(StringConverterInterface)
	if !ok {
		return fmt.Errorf(
			"mapType key type not implement stringConverterType,key: %s",
			t.KeyType.GetReflectType().Kind().String(),
		)
	}
	t.ElemType, err = TypeOf(t.GetReflectType().Elem())
	if err != nil {
		return err
	}
	return nil
}

func (t *MapType) GetElemByString(v reflect.Value, k string) (ev reflect.Value, et UdwType, err error) {
	err = t.Init()
	if err != nil {
		return
	}
	vk, err := t.KeyStringConverter.FromString(k)
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	ev = v.MapIndex(vk)
	if !ev.IsValid() {
		err = fmt.Errorf("[mapType.getSubValueByString] map key not found k:%s", k)
		return
	}
	et = t.ElemType
	return
}

func (t *MapType) SaveByPath(v *reflect.Value, path Path, value string) (err error) {
	err = t.Init()
	if err != nil {
		return
	}
	if len(path) == 0 {
		return fmt.Errorf("[mapType.save] get map with no path, value:%s", value)
	}

	if v.IsNil() {
		if v.CanSet() {
			v.Set(reflect.MakeMap(t.GetReflectType()))
		} else {
			*v = reflect.MakeMap(t.GetReflectType())
		}
	}
	vk, err := t.KeyStringConverter.FromString(path[0])
	if err != nil {
		return err
	}
	saveElemV := v.MapIndex(vk)
	KeyNotExist := false
	if !saveElemV.IsValid() {
		saveElemV = reflect.New(t.ElemType.GetReflectType()).Elem()
		KeyNotExist = true
	}
	oElemV := saveElemV
	err = t.ElemType.SaveByPath(&saveElemV, path[1:], value)
	if err != nil {
		return err
	}
	if KeyNotExist {
		v.SetMapIndex(vk, saveElemV)
	}
	if oElemV != saveElemV {
		v.SetMapIndex(vk, saveElemV)
	}
	return nil
}

func (t *MapType) DeleteByPath(v *reflect.Value, path Path) (err error) {
	err = t.Init()
	if err != nil {
		return
	}
	if len(path) > 1 {
		vk, err := t.KeyStringConverter.FromString(path[0])
		if err != nil {
			return err
		}
		ev := v.MapIndex(vk)
		et := t.ElemType
		oEv := ev
		err = et.DeleteByPath(&ev, path[1:])
		if err != nil {
			return err
		}
		if oEv == ev {
			return nil
		}
		v.SetMapIndex(vk, ev)
		return nil
	} else if len(path) == 0 {
		return fmt.Errorf("[MapType.DeleteByPath] delete map with no path.")
	}

	vk, err := t.KeyStringConverter.FromString(path[0])
	if err != nil {
		return err
	}
	v.SetMapIndex(vk, reflect.Value{})
	return nil
}
