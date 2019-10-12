package udwGoParser

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoReader"
)

func MustParseGoFunc(s string) *FuncOrMethodDeclaration {
	gofile := &file{
		importMap:      map[string]bool{},
		aliasImportMap: map[string]string{},
		pkg: &Package{
			program:       NewProgramFromDefault(),
			pkgImportPath: "main",
		},
	}
	r := udwGoReader.NewReaderWithBuf([]byte(s))
	return gofile.readGoFunc(r)
}

type MustParseGoFuncDeclarationRequest struct {
	Pkg                      *Package
	ImportList               []string
	GoFuncDeclarationContent string
}

func MustParseGoFuncDeclaration(req MustParseGoFuncDeclarationRequest) *FuncOrMethodDeclaration {
	if req.Pkg == nil {
		req.Pkg = &Package{
			program:       NewProgramFromDefault(),
			pkgImportPath: "main",
		}
	}
	gofile := &file{
		importMap:      map[string]bool{},
		aliasImportMap: map[string]string{},
	}
	gofile.pkg = req.Pkg
	for _, importPath := range req.ImportList {
		gofile.addImport(importPath, "")
	}
	r := udwGoReader.NewReaderWithBuf([]byte(req.GoFuncDeclarationContent))
	return gofile.readGoFunc(r)
}

func (gofile *file) readGoFunc(r *udwGoReader.Reader) *FuncOrMethodDeclaration {
	var isVariadic bool
	funcDecl := &FuncOrMethodDeclaration{}
	r.MustReadMatch([]byte("func"))

	r.ReadAllSpace()
	b := r.ReadByte()
	if b == '(' {
		r.UnreadByte()
		receiver, isVariadic := gofile.readParameters(r)
		if len(receiver) != 1 {
			panic(fmt.Errorf("%s receiver must have one parameter", r.GetFileLineInfo()))
		}
		funcDecl.ReceiverType = receiver[0].Type
		if isVariadic {
			panic("[gofile.readGoFunc] found isVariadic(...) in receiver")
		}
		r.ReadAllSpace()

	} else {
		r.UnreadByte()
	}
	id := readIdentifier(r)
	funcDecl.Name = string(id)
	if funcDecl.Name == "" {
		panic(fmt.Errorf("%s need function name", r.GetFileLineInfo()))
	}
	r.ReadAllSpace()
	b = r.ReadByte()
	if b != '(' {
		panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
	}
	r.UnreadByte()

	funcDecl.InParameter, isVariadic = gofile.readParameters(r)
	funcDecl.IsVariadic = isVariadic
	r.ReadAllSpaceWithoutLineBreak()
	if r.IsEof() {
		return funcDecl
	}
	b = r.ReadByte()
	if b == '\n' {
		return funcDecl
	} else if b != '{' {
		r.UnreadByte()
		funcDecl.OutParameter, isVariadic = gofile.readParameters(r)
		if isVariadic {
			panic("[gofile.readGoFunc] found isVariadic(...) in outParameter")
		}
		r.ReadAllSpaceWithoutLineBreak()
		if r.IsEof() {
			return funcDecl
		}
		b = r.ReadByte()
	}
	if b == '\n' {
		return funcDecl
	}
	if b != '{' {
		panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
	}

	readMatchBigParantheses(r)
	return funcDecl
}

func (gofile *file) readParameters(r *udwGoReader.Reader) (output []FuncParameter, isVariadic bool) {
	b := r.ReadByte()
	if b != '(' {

		r.UnreadByte()
		return []FuncParameter{
			{
				Type: gofile.readType(r),
			},
		}, false
	}
	parameterPartList := []*astParameterPart{}
	lastPart := &astParameterPart{}
	for {
		r.ReadAllSpace()
		b := r.ReadByte()
		if b == ')' || b == ',' {
			if lastPart.partList[0].originByte != nil {
				parameterPartList = append(parameterPartList, lastPart)
				lastPart = &astParameterPart{}
			}
			if b == ')' {
				break
			}
			if b == ',' {
				continue
			}
		}

		r.UnreadByte()
		if r.IsMatchAfter([]byte("...")) {
			r.MustReadMatch([]byte("..."))
			isVariadic = true
		}
		startPos := r.Pos()
		typ := gofile.readType(r)
		buf := r.BufToCurrent(startPos)

		hasSet := false
		for i := range lastPart.partList {
			if lastPart.partList[i].originByte == nil {
				lastPart.partList[i].originByte = buf
				lastPart.partList[i].typ = typ
				hasSet = true
				break
			}
		}
		if !hasSet {
			panic(r.GetFileLineInfo() + " unexcept func parameterList.")
		}
	}

	output = make([]FuncParameter, len(parameterPartList))
	onlyHavePart1Num := 0
	for i := range parameterPartList {
		if parameterPartList[i].partList[1].originByte == nil {
			onlyHavePart1Num++
		}
	}

	if onlyHavePart1Num == len(parameterPartList) {
		for i := range parameterPartList {
			output[i].Type = parameterPartList[i].partList[0].typ
		}
		return output, isVariadic
	}

	for i, parameterPart := range parameterPartList {
		output[i].Name = string(parameterPart.partList[0].originByte)
		if parameterPart.partList[1].typ != nil {
			output[i].Type = parameterPart.partList[1].typ
		}
	}

	for i := range parameterPartList {
		if output[i].Type == nil {
			for j := i + 1; j < len(parameterPartList); j++ {
				if output[j].Type != nil {
					output[i].Type = output[j].Type
				}
			}
		}
	}
	return output, isVariadic
}

type astParameterPart struct {
	partList [2]struct {
		originByte []byte
		typ        Type
	}
	isVariadic bool
}
