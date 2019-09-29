package udwType

import (
	"fmt"
	"reflect"
)

type Context struct {
	RootType  UdwType
	RootValue reflect.Value
}

func NewContext(ptr interface{}) (*Context, error) {
	rt := reflect.TypeOf(ptr)
	if rt.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("NewContext need a ptr,but get %T", ptr)
	}
	context := &Context{}
	context.RootValue = reflect.ValueOf(ptr)
	et, err := TypeOf(rt)
	if err != nil {
		return nil, err
	}
	context.RootType = et
	return context, nil
}
func (m *Context) SaveByPath(path Path, value string) (err error) {
	oEv := m.RootValue
	err = m.RootType.SaveByPath(&oEv, path, value)
	if err != nil {
		return
	}
	if oEv != m.RootValue {
		err = fmt.Errorf("[context.SaveByPath] can not save")
		return
	}
	return nil
}

func (m *Context) DeleteByPath(path Path) (err error) {
	oEv := m.RootValue
	err = m.RootType.DeleteByPath(&oEv, path)
	if err != nil {
		return
	}
	if oEv != m.RootValue {
		err = fmt.Errorf("[context.DeleteByPath] can not save")
		return
	}
	return nil
}

func (m *Context) GetElemByPath(p Path) (v reflect.Value, t UdwType, err error) {
	t = m.RootType
	v = m.RootValue
	for _, ps := range p {
		getter, ok := t.(GetElemByStringInterface)
		if ok == false {
			err = fmt.Errorf("[udwType.context.GetElemByPath] some stuff in path not gettable path:%s", p)
			return
		}
		v, t, err = getter.GetElemByString(v, ps)
		if err != nil {
			return
		}
	}
	return
}
