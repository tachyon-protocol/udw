package udwRpc2Tester

import (
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Builder"
)

func BuildAndTest() {

	udwRpc2Builder.Generate(udwRpc2Builder.GenerateReq{
		RpcDefine:      getRpcService(),
		FromPkgPath:    "github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Tester/udwRpc2Demo",
		FromObjName:    "Server",
		TargetPkgPath:  "github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Tester/udwRpc2Demo",
		Prefix:         "Demo",
		TargetFilePath: "src/github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Tester/udwRpc2Demo/zzzig_Demo.go",
		GoFmt:          true,
	})

	udwRpc2Builder.Generate(udwRpc2Builder.GenerateReq{
		RpcDefine:        getRpcService(),
		FromPkgPath:      "github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Tester/udwRpc2Demo",
		FromObjName:      "Server",
		TargetPkgPath:    "main",
		Prefix:           "D2",
		TargetFilePath:   "src/github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Tester/udwRpc2DemoClient/zzzig_D2.go",
		GoFmt:            true,
		DisableGenServer: true,
	})
	udwCmd.MustRun("go run github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Tester/udwRpc2DemoClient")

}

func getRpcService() udwRpc2Builder.RpcService {
	return udwRpc2Builder.RpcService{
		List: []udwRpc2Builder.RpcApi{
			{
				Name: "SetName",
				InputParameterList: []udwRpc2Builder.RpcParameter{
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindString,
						},
					},
				},
			},
			{
				Name: "GetName",
				OutputParameterList: []udwRpc2Builder.RpcParameter{
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindString,
						},
					},
				},
			},
			{
				Name: "IncreaseInt",
			},
			{
				Name: "GetInt",
				OutputParameterList: []udwRpc2Builder.RpcParameter{
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindInt,
						},
					},
				},
			},
			{
				Name: "Panic",
			},
			{
				Name: "FnP",
				InputParameterList: []udwRpc2Builder.RpcParameter{
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindString,
						},
					},
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindString,
						},
					},
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindString,
						},
					},
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindSlice,
							Elem: &udwRpc2Builder.RpcType{
								Kind:       udwRpc2Builder.RpcTypeKindNamedStruct,
								StructName: "Tstruct",
								GoPkg:      "github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Tester/udwRpc2Demo",
							},
						},
					},
				},
				OutputParameterList: []udwRpc2Builder.RpcParameter{
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindString,
						},
					},
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindString,
						},
					},
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindString,
						},
					},
				},
			},
			{
				Name: "GetPeerIp",
				InputParameterList: []udwRpc2Builder.RpcParameter{
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindString,
						},
					},
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindString,
						},
					},
					{
						Type: udwRpc2Builder.RpcType{
							Kind:       udwRpc2Builder.RpcTypeKindNamedStruct,
							StructName: "PeerIp",
							GoPkg:      "github.com/tachyon-protocol/udw/udwRpc2",
						},
					},
				},
				OutputParameterList: []udwRpc2Builder.RpcParameter{
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindString,
						},
					},
				},
			},
		},
	}
}
