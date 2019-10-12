package udwGoWriter

type WriteTestWrapFunctionReq struct {
	PkgPath  string
	FuncName string
	FuncBody []byte
	GoFile   *GoFileContext
}

func WriteTestWrapFunction(req WriteTestWrapFunctionReq) {
	req.GoFile.Buf.WriteString(`//udw go test -v ` + req.PkgPath + ` -run Test_` + req.FuncName + `
func Test_` + req.FuncName + `(t *testing.T) {
	`)
	req.GoFile.Buf.Write(req.FuncBody)
	req.GoFile.Buf.WriteString("\n}\n")
}
