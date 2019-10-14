package udwGoParser

import "github.com/tachyon-protocol/udw/udwGoSource/udwGoReader"

func (gofile *file) readImport(r *udwGoReader.Reader) {
	r.MustReadMatch([]byte("import"))
	r.ReadAllSpace()
	if r.IsEof() {
		panic(r.GetFileLineInfo() + " unexcept EOF ")
	}
	b := r.ReadByte()
	if b == '(' {
		for {
			r.ReadAllSpace()
			b := r.ReadByte()
			if b == ')' {
				return
			} else {
				r.UnreadByte()
				gofile.readImportSpec(r)
			}
		}
	} else {
		r.UnreadByte()
		gofile.readImportSpec(r)
	}
}

func (gofile *file) readImportSpec(r *udwGoReader.Reader) {
	r.ReadAllSpace()
	b := r.ReadByte()

	if b == '"' || b == '`' {
		r.UnreadByte()
		gofile.addImport(string(mustReadGoString(r)), "")
	} else if b == '.' {
		r.ReadAllSpace()
		gofile.addImport(string(mustReadGoString(r)), ".")
	} else {
		r.UnreadByte()
		alias := readIdentifier(r)
		r.ReadAllSpace()
		gofile.addImport(string(mustReadGoString(r)), string(alias))
	}
}
