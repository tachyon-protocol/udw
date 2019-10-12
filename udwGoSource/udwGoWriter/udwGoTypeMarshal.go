package udwGoWriter

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal"
	"path"
	"reflect"
)

func MustWriteObjectToMainPackage(obj interface{}) string {
	return udwGoTypeMarshal.MustWriteObjectToMainPackage(obj)
}

func (file *GoFileContext) MustWriteObject(obj interface{}) string {
	buf := &bytes.Buffer{}
	file.MustWriteObjectToBuf(obj, buf)
	return buf.String()
}

func (file *GoFileContext) MustWriteObjectToSelfBuf(obj interface{}) {
	file.MustWriteObjectToBuf(obj, &file.Buf)
	return
}

func (file *GoFileContext) MustWriteObjectToBuf(obj interface{}, _buf *bytes.Buffer) {
	udwGoTypeMarshal.MustGoTypeMarshalWithBuf(udwGoTypeMarshal.MustGoTypeMarshalContext{
		AddTypeNameWithImportPackage: func(pkgPath string, name string, buf *bytes.Buffer) {
			if file.pkgImportPath == pkgPath {
				buf.WriteString(name)
				return
			} else {
				file.AddImportPath(pkgPath)
				buf.WriteString(path.Base(pkgPath) + "." + name)
				return
			}
		},
	}, obj, _buf)
	return
}

func (file *GoFileContext) MustWriteTypeByReflect(typ reflect.Type) string {
	buf := &bytes.Buffer{}
	udwGoTypeMarshal.WriteReflectTypeNameToGoFile(udwGoTypeMarshal.MustGoTypeMarshalContext{
		AddTypeNameWithImportPackage: func(pkgPath string, name string, buf *bytes.Buffer) {
			if file.pkgImportPath == pkgPath {
				buf.WriteString(name)
				return
			} else {
				file.AddImportPath(pkgPath)
				buf.WriteString(path.Base(pkgPath) + "." + name)
				return
			}
		},
	}, typ, buf)
	return buf.String()
}
