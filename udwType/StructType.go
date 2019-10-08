package udwType

import (
	"fmt"
	"reflect"
)

type StructType struct {
	reflectTypeGetterImp
	getElemByStringEditorabler
}

func (t *StructType) GetElemByString(v reflect.Value, k string) (ev reflect.Value, et UdwType, err error) {
	ev = v.FieldByName(k)
	if !ev.IsValid() {
		err = fmt.Errorf("field %s not find in struct", k)
		return
	}
	et, err = TypeOf(ev.Type())
	return
}

func (t *StructType) DeleteByPath(v *reflect.Value, path Path) (err error) {
	if len(path) > 1 {
		return passThougthDeleteByPath(t, v, path)
	}
	return fmt.Errorf("[StructType.DeleteByPath] can not delete f struct")
}
