package udwGoParser

import (
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoReader"
)

func (gofile *file) readGoType(r *udwGoReader.Reader) {
	r.MustReadMatch([]byte("type"))
	r.ReadAllSpace()
	b := r.ReadByte()
	if b == '(' {
		for {
			r.ReadAllSpace()
			b = r.ReadByte()
			if b == ')' {
				return
			}
			r.UnreadByte()
			name := readIdentifier(r)
			r.ReadAllSpace()
			typ := gofile.readType(r)
			gofile.namedTypeList = append(gofile.namedTypeList, &NamedType{
				PkgImportPath: gofile.pkg.pkgImportPath,
				Name:          string(name),
				UnderType:     typ,
			})
		}
	} else {
		r.UnreadByte()
		name := readIdentifier(r)
		r.ReadAllSpace()
		typ := gofile.readType(r)
		gofile.namedTypeList = append(gofile.namedTypeList, &NamedType{
			PkgImportPath: gofile.pkg.pkgImportPath,
			Name:          string(name),
			UnderType:     typ,
		})
		return
	}
}

func (gofile *file) readGoVar(r *udwGoReader.Reader) {
	r.MustReadMatch([]byte("var"))
	r.ReadAllSpace()
	b := r.ReadByte()
	if b == '(' {
		readMatchSmallParantheses(r)
		return
	}
	r.UnreadByte()
	readIdentifier(r)
	r.ReadAllSpace()
	b = r.ReadByte()
	if b == '=' {
		r.ReadAllSpace()
	}
	for {
		if b == '"' || b == '`' {
			r.UnreadByte()
			mustReadGoString(r)
		}
		if b == '\'' {
			r.UnreadByte()
			mustReadGoChar(r)
		}
		if b == '\n' {
			return
		}
		if b == '{' {
			readMatchBigParantheses(r)

		}
		if b == '(' {
			readMatchSmallParantheses(r)
		}
		if r.IsEof() {
			return
		}
		b = r.ReadByte()
	}
}

func (gofile *file) readGoConst(r *udwGoReader.Reader) {
	r.MustReadMatch([]byte("const"))
	r.ReadAllSpace()
	b := r.ReadByte()
	if b == '(' {
		readMatchSmallParantheses(r)
		return
	}
	for {
		if b == '"' || b == '`' {
			r.UnreadByte()
			mustReadGoString(r)
		}
		if b == '\'' {
			r.UnreadByte()
			mustReadGoChar(r)
		}
		if b == '\n' {
			return
		}
		if r.IsEof() {
			return
		}
		b = r.ReadByte()
	}
}
