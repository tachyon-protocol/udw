package udwGoParser

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoReader"
	"path/filepath"
	"strings"
	"unicode"
)

func parseFile(pkgPath string, path string, pkg *Package) *file {
	gofile := &file{
		importMap:      map[string]bool{},
		aliasImportMap: map[string]string{},
		pkg:            pkg,
	}
	content := udwFile.MustReadFile(path)
	if bytes.Contains(content, []byte("// +build ignore\n")) {
		return gofile
	}

	posFile := udwGoReader.NewPosFile(path, content)
	content = goSourceRemoveComment(content, posFile)
	r := udwGoReader.NewReader(content, posFile)
	r.ReadAllSpace()
	r.MustReadMatch(tokenPackage)
	r.ReadAllSpace()
	buf := r.ReadUntilByte('\n')
	gofile.packageName = strings.TrimSpace(string(buf))
	if pkg.pkgImportPath == "" {
		pkg.pkgImportPath = pkg.pkgPath
	}
	if !strings.Contains(filepath.Base(path), "_test") {
		if pkg.pkgName != "" && pkg.pkgName != gofile.packageName {
			panic(fmt.Sprintln("two packageName in same directory", pkg.pkgPath, path, pkg.pkgName, gofile.packageName))
		}
		pkg.pkgName = gofile.packageName
		if pkg.IsMain() {
			pkg.pkgImportPath = "main"
		} else {
			pkg.pkgImportPath = pkg.pkgPath
		}
	}
	for {
		if r.IsEof() {
			return gofile
		}
		r.ReadAllSpace()
		if r.IsMatchAfter(tokenImport) {
			gofile.readImport(r)
			continue
		}
		break
	}
	for {
		switch {
		case r.IsEof():
			return gofile
		case r.IsMatchAfter(tokenFunc):
			funcDecl := gofile.readGoFunc(r)
			if funcDecl.GetKind() == DefinedFunc {
				gofile.funcList = append(gofile.funcList, funcDecl)
			} else {
				gofile.methodList = append(gofile.methodList, funcDecl)
			}
		case r.IsMatchAfter(tokenType):

			gofile.readGoType(r)
		case r.IsMatchAfter(tokenVar):
			gofile.readGoVar(r)
		case r.IsMatchAfter(tokenConst):
			gofile.readGoConst(r)

		case r.IsMatchAfter(tokenDoubleQuate) || r.IsMatchAfter(tokenGraveAccent):
			mustReadGoString(r)

		case r.IsMatchAfter(tokenSingleQuate):
			mustReadGoChar(r)
		default:
			r.ReadByte()
		}
	}
}

func readIdentifier(r *udwGoReader.Reader) []byte {
	buf := &bytes.Buffer{}
	if r.IsEof() {
		panic(r.GetFileLineInfo() + " unexcept EOF")
	}
	b := r.ReadRune()
	if b == '_' || unicode.IsLetter(b) {
		buf.WriteRune(b)
	} else {
		r.UnreadRune()
		return nil
	}
	for {
		if r.IsEof() {
			return buf.Bytes()
		}
		b := r.ReadRune()
		if b == '_' || unicode.IsLetter(b) || unicode.IsDigit(b) {
			buf.WriteRune(b)
		} else {
			r.UnreadRune()
			return buf.Bytes()
		}
	}
}

func readMatchBigParantheses(r *udwGoReader.Reader) []byte {
	return readMatchChar(r, '{', '}')
}

func readMatchMiddleParantheses(r *udwGoReader.Reader) []byte {
	return readMatchChar(r, '[', ']')
}

func readMatchSmallParantheses(r *udwGoReader.Reader) []byte {
	return readMatchChar(r, '(', ')')
}

func readMatchChar(r *udwGoReader.Reader, starter byte, ender byte) []byte {
	startPos := r.Pos()
	level := 1
	for {
		if r.IsEof() {
			panic(r.GetFileLineInfo() + " unexcept EOF")
		}
		b := r.ReadByte()
		if b == '"' || b == '`' {
			r.UnreadByte()
			mustReadGoString(r)
		} else if b == '\'' {
			r.UnreadByte()
			mustReadGoChar(r)
		} else if b == starter {
			level++
		} else if b == ender {
			level--
			if level == 0 {
				return r.BufToCurrent(startPos)
			}
		}
	}
}
