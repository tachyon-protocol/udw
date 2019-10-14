package udwGoParser

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoReader"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"unicode"
)

func (gofile *file) readType(r *udwGoReader.Reader) Type {
	id := readIdentifier(r)
	if len(id) == 0 {
		if r.IsMatchAfter([]byte("<-chan")) {
			r.MustReadMatch([]byte("<-chan"))
			r.ReadAllSpace()
			return &ChanType{
				Dir:  RecvDir,
				Elem: gofile.readType(r),
			}
		}
		b := r.ReadByte()
		if b == '*' {
			return &PointerType{
				Elem: gofile.readType(r),
			}
		} else if b == '[' {
			content := readMatchMiddleParantheses(r)
			if len(content) == 1 {
				return &SliceType{
					Elem: gofile.readType(r),
				}
			}
			if len(content) == 0 {
				panic(fmt.Errorf("%s expect ]", r.GetFileLineInfo()))
			}
			if len(content) > 1 {
				size, err := udwStrconv.ParseInt(string(content[:len(content)-1]))
				if err != nil {
					fmt.Printf("[udwGoParser] %s unexpect/unsupported array length: %s\n", r.GetFileLineInfo(), string(content[:len(content)-1]))
				}
				return &ArrayType{
					Elem: gofile.readType(r),
					Size: size,
				}
			}
		} else if b == '(' {
			typ := gofile.readType(r)
			r.MustReadMatch([]byte(")"))
			return typ
		} else {
			panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
		}
	}
	if bytes.Equal(id, []byte("struct")) {
		return gofile.readStruct(r)
	} else if bytes.Equal(id, []byte("interface")) {

		r.ReadAllSpace()
		b := r.ReadByte()
		if b != '{' {
			panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
		}
		readMatchBigParantheses(r)
		return InterfaceType{}
	} else if bytes.Equal(id, []byte("map")) {
		b := r.ReadByte()
		m := &MapType{}
		if b != '[' {
			panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
		}
		m.Key = gofile.readType(r)
		r.MustReadMatch([]byte("]"))
		m.Value = gofile.readType(r)
		return m
	} else if bytes.Equal(id, []byte("func")) {

		r.ReadAllSpace()
		b := r.ReadByte()
		if b != '(' {
			panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
		}
		readMatchSmallParantheses(r)
		r.ReadAllSpaceWithoutLineBreak()
		run := r.ReadRune()
		if run == '(' {
			readMatchSmallParantheses(r)
		} else if run == '\n' {
			return &FuncType{}
		} else if unicode.IsLetter(run) || run == '[' || run == '*' || run == '<' {
			r.UnreadRune()
			gofile.readType(r)
		} else {
			r.UnreadRune()
		}
		return &FuncType{}
	} else if bytes.Equal(id, []byte("chan")) {
		if r.IsMatchAfter([]byte("<-")) {
			r.MustReadMatch([]byte("<-"))
			r.ReadAllSpace()
			return &ChanType{
				Dir:  SendDir,
				Elem: gofile.readType(r),
			}
		} else {
			r.ReadAllSpace()
			return &ChanType{
				Dir:  BothDir,
				Elem: gofile.readType(r),
			}
		}
	} else {
		if r.IsEof() {
			return gofile.readTypeOneIdentifier(id)
		}
		b := r.ReadByte()
		if b == '.' {
			pkgPath := string(id)
			pkgPath, err := gofile.lookupFullPackagePath(pkgPath)
			if err != nil {
				fmt.Println(r.GetFileLineInfo(), err.Error())
			}
			id2 := readIdentifier(r)
			return &NamedType{
				PkgImportPath: pkgPath,
				Name:          string(id2),
				program:       gofile.pkg.program,
			}
		} else {
			r.UnreadByte()
			return gofile.readTypeOneIdentifier(id)
		}
	}

}

func (gofile *file) readTypeOneIdentifier(id []byte) Type {
	name := string(id)
	if getKindFromBuiltinType(name) != Invalid {
		return BuiltinType(name)
	} else {
		return &NamedType{
			PkgImportPath: gofile.pkg.pkgImportPath,
			Name:          string(id),
			program:       gofile.pkg.program,
		}
	}
}

func getTypeStructAnonymousName(typ Type) string {
	ntyp, ok := typ.(*NamedType)
	if ok {
		return ntyp.Name
	}
	ptyp, ok := typ.(*PointerType)
	if ok {
		return "*" + getTypeStructAnonymousName(ptyp.Elem)
	}
	btyp, ok := typ.(BuiltinType)
	if ok {
		return string(btyp)
	}
	panic(fmt.Errorf("[getTypeStructAnonymousName] unexpect type %T", typ))
}

func (gofile *file) readStruct(r *udwGoReader.Reader) *StructType {

	r.ReadAllSpace()
	b := r.ReadByte()
	if b != '{' {
		panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
	}
	lastReadBuf := []bytesAndType{}
	var lastTag []byte
	out := &StructType{}
	for {
		r.ReadAllSpaceWithoutLineBreak()
		b := r.ReadByte()
		if b == '}' {
			return out
		} else if b == '"' || b == '\'' || b == '`' {
			r.UnreadByte()
			lastTag = mustReadGoString(r)
		} else if b == ',' {
			continue
		} else if b == '\n' {
			if len(lastReadBuf) == 0 {
				continue
			} else if len(lastReadBuf) == 1 {
				typ := lastReadBuf[0].typ
				name := getTypeStructAnonymousName(typ)
				out.Field = append(out.Field, StructField{
					Name:             name,
					Elem:             typ,
					IsAnonymousField: true,
					Tag:              string(lastTag),
				})

				lastReadBuf = []bytesAndType{}
			} else if len(lastReadBuf) >= 2 {
				typ := lastReadBuf[len(lastReadBuf)-1].typ
				for i := range lastReadBuf[:len(lastReadBuf)-1] {
					out.Field = append(out.Field, StructField{
						Name:             string(lastReadBuf[i].originByte),
						Elem:             typ,
						IsAnonymousField: false,
						Tag:              string(lastTag),
					})
				}
				lastReadBuf = []bytesAndType{}
			}
			lastTag = nil
		} else {
			r.UnreadByte()
			startPos := r.Pos()
			typ := gofile.readType(r)
			lastReadBuf = append(lastReadBuf, bytesAndType{
				originByte: r.BufToCurrent(startPos),
				typ:        typ,
			})
		}
	}
}

type bytesAndType struct {
	originByte []byte
	typ        Type
}
