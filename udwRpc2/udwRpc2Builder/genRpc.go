package udwRpc2Builder

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoParser"
	"go/format"
	"strconv"
)

type GenerateReq struct {
	RpcDefine      RpcService
	TargetPkgPath  string
	Prefix         string
	TargetFilePath string
	GoFmt          bool
}

func Generate(req GenerateReq) {
	ctx := &generateCtx{
		req: req,
	}
	initRpcDefine(&ctx.req.RpcDefine)
	ctx.goFileWriter = udwGoParser.NewGoFileContext(req.TargetPkgPath)
	ctx.goFileWriter.AddImportPath("github.com/tachyon-protocol/udw/udwRpc2")
	ctx.generateServer()
	ctx.generateClient()
	ctx.goFileWriter.MustWriteFile(req.TargetFilePath, ctx.goFileBuf.GetBytes())
	if req.GoFmt {
		content := udwFile.MustReadFile(req.TargetFilePath)
		out, err := format.Source(content)
		if err != nil {
			panic(err)
		}
		udwFile.MustWriteFileWithMkdir(req.TargetFilePath, out)
	}
}

type generateCtx struct {
	req          GenerateReq
	goFileWriter *udwGoParser.GoFileWriter
	goFileBuf    udwBytes.BufWriter
	nameIndex    int
}

func (ctx *generateCtx) newName() string {
	ctx.nameIndex++
	return "tmp_" + strconv.Itoa(ctx.nameIndex)
}

func (ctx *generateCtx) generateServer() {
	ctx.goFileBuf.WriteString_(`func ` + ctx.req.Prefix + `_RunServer(addr string) (closer func()){
	s:=Server{}
	sh:=udwRpc2.NewServerHub(udwRpc2.ServerReq{
		Addr: addr,
		Handler: func(ctx *udwRpc2.ReqCtx){
			var fnId uint64
			var errMsg string
			fnId,errMsg = ctx.GetReader().ReadUvarint()
			if errMsg!=""{
				return
			}
			panicErrMsg:=udwRpc2.PanicToErrMsg(func(){
			switch fnId {
`)
	for _, api := range ctx.req.RpcDefine.List {
		ctx.goFileBuf.WriteString_(`case ` + strconv.Itoa(api.Pos) + `:
`)
		inNameList := []string{}
		for _, parameter := range api.InputParameterList {
			n1 := ctx.newName()
			inNameList = append(inNameList, n1)
			if parameter.IsClientIp {
				ctx.goFileBuf.WriteString_(n1 + `:=udwRpc2.PeerIp{ctx.GetPeerIp()}
`)
				continue
			}
			goTypeString := mustRpcTypeToGoTypeName(ctx.goFileWriter, &parameter.Type)
			ctx.goFileBuf.WriteString_(`var ` + n1 + ` ` + goTypeString + `
errMsg = ctx.GetReader().ReadValue(&` + n1 + `)
if errMsg!=""{
	return
}
`)
		}
		ctx.goFileBuf.WriteString_(`errMsg = ctx.GetReader().ReadArrayEnd()
if errMsg!=""{
	return
}
`)
		outNameList := []string{}
		if len(api.OutputParameterList) > 0 {
			for i := range api.OutputParameterList {
				n1 := ctx.newName()
				outNameList = append(outNameList, n1)
				ctx.goFileBuf.WriteString_(n1)
				if i != len(api.OutputParameterList)-1 {
					ctx.goFileBuf.WriteString_(`,`)
				}
			}
			ctx.goFileBuf.WriteString_(`:=`)
		}
		ctx.goFileBuf.WriteString_(`s.` + api.Name + `(`)
		for _, name := range inNameList {
			ctx.goFileBuf.WriteString_(name + `,`)
		}
		ctx.goFileBuf.WriteString_(`)
ctx.GetWriter().WriteString("")
`)
		for i := range api.OutputParameterList {
			ctx.goFileBuf.WriteString_(`errMsg = ctx.GetWriter().WriteValue(` + outNameList[i] + `)
if errMsg!=""{
	return
}
`)
		}
		ctx.goFileBuf.WriteString_(`ctx.GetWriter().WriteArrayEnd()
errMsg = ctx.GetWriter().Flush()
if errMsg!=""{
	return
}
`)
	}
	ctx.goFileBuf.WriteString_(`default:
				}
			})
			if panicErrMsg!=""{
				ctx.GetWriter().WriteString(panicErrMsg)
				ctx.GetWriter().WriteArrayEnd()
				errMsg = ctx.GetWriter().Flush()
				if errMsg!=""{
					return
				}
			}
		},
	})
	return sh.Close
}
`)
}

func (ctx *generateCtx) generateClient() {
	ctx.goFileBuf.WriteString_(`func ` + ctx.req.Prefix + `_NewClient(addr string) (*` + ctx.req.Prefix + `_Client){
	c:=udwRpc2.NewClientHub(udwRpc2.ClientReq{
		Addr: addr,
	})
	return &Demo_Client{
		ch: c,
	}
}
type ` + ctx.req.Prefix + `_Client struct{
	ch *udwRpc2.ClientHub
}
`)
	for _, api := range ctx.req.RpcDefine.List {
		ctx.goFileBuf.WriteString_(`func (c *` + ctx.req.Prefix + `_Client) ` + api.Name + `(`)
		for _, parameter := range api.InputParameterList {
			if parameter.IsClientIp {
				continue
			}
			goTypeString := mustRpcTypeToGoTypeName(ctx.goFileWriter, &parameter.Type)
			ctx.goFileBuf.WriteString_(parameter.Name + ` ` + goTypeString + `,`)
		}
		ctx.goFileBuf.WriteString_(`)(`)
		for _, parameter := range api.OutputParameterList {
			goTypeString := mustRpcTypeToGoTypeName(ctx.goFileWriter, &parameter.Type)
			ctx.goFileBuf.WriteString_(parameter.Name + ` ` + goTypeString + `,`)
		}
		ctx.goFileBuf.WriteString_(`RpcErr *udwRpc2.RpcError){
	_networkErr:=c.ch.RequestCb(func(ctx *udwRpc2.ReqCtx){
		ctx.GetWriter().WriteUvarint(` + strconv.Itoa(api.Pos) + `)
`)
		for _, parameter := range api.InputParameterList {
			if parameter.IsClientIp {
				continue
			}
			ctx.goFileBuf.WriteString_(`		ctx.GetWriter().WriteValue(` + parameter.Name + `)
`)
		}
		ctx.goFileBuf.WriteString_(`		ctx.GetWriter().WriteArrayEnd()
		errMsg:=ctx.GetWriter().Flush()
		if errMsg!=""{
			RpcErr = udwRpc2.NewNetworkError("dehqx82rjj "+errMsg)
			return
		}
		var s string
		errMsg = ctx.GetReader().ReadValue(&s)
		if errMsg!=""{
			RpcErr = udwRpc2.NewNetworkError("ehtjkea4re "+errMsg)
			return
		}
		if s!=""{
			RpcErr = udwRpc2.NewOtherError(s)
			ctx.GetReader().ReadArrayEnd()
			return
		}
`)
		for _, parameter := range api.OutputParameterList {
			ctx.goFileBuf.WriteString_(`		errMsg = ctx.GetReader().ReadValue(&` + parameter.Name + `)
		if errMsg!=""{
			RpcErr = udwRpc2.NewNetworkError("kvkdcgtnk2 "+errMsg)
			return
		}
`)
		}
		ctx.goFileBuf.WriteString_(`		errMsg = ctx.GetReader().ReadArrayEnd()
		if errMsg!=""{
			RpcErr = udwRpc2.NewNetworkError("4b7rug5mf2 "+errMsg)
			return
		}
		RpcErr = nil
		return
	})
	if _networkErr!=""{
		RpcErr = udwRpc2.NewNetworkError("494fehebw6 "+_networkErr)
	}
	return
}
`)

	}
}
