package udwRpc2Tester

import (
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Builder"
)

func BuildAndTest() {
	udwRpc2Builder.Generate(udwRpc2Builder.GenerateReq{
		RpcDefine:      getRpcService(),
		TargetPkgPath:  "main",
		Prefix:         "Demo",
		TargetFilePath: "src/github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Tester/udwRpc2Demo/zzzig_Demo.go",
		GoFmt:          true,
	})
	udwCmd.MustRun("go run github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Tester/udwRpc2Demo")
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
