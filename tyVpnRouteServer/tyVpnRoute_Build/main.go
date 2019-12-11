package tyVpnRoute_Build

import (
	"github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Builder"
)

func DsnBuild() {
	udwRpc2Builder.Generate(udwRpc2Builder.GenerateReq{
		RpcDefine:        getRpcService(),
		FromObjName:      "serverRpcObj",
		FromPkgPath:      "github.com/tachyon-protocol/udw/tyVpnRouteServer",
		TargetPkgPath:    "github.com/tachyon-protocol/udw/tyVpnRouteServer",
		Prefix:           "Rpc",
		TargetFilePath:   "src/github.com/tachyon-protocol/udw/tyVpnRouteServer/rpc.go",
		GoFmt:            true,
		DisableGenClient: true,
	})
	udwRpc2Builder.Generate(udwRpc2Builder.GenerateReq{
		RpcDefine:        getRpcService(),
		FromObjName:      "serverRpcObj",
		FromPkgPath:      "github.com/tachyon-protocol/udw/tyVpnRouteServer",
		TargetPkgPath:    "github.com/tachyon-protocol/udw/tyVpnRouteServer/tyVpnRouteClient",
		Prefix:           "Rpc",
		TargetFilePath:   "src/github.com/tachyon-protocol/udw/tyVpnRouteServer/tyVpnRouteClient/rpc.go",
		GoFmt:            true,
		DisableGenServer: true,
	})
}

func getRpcService() udwRpc2Builder.RpcService {
	return udwRpc2Builder.RpcService{
		List: []udwRpc2Builder.RpcApi{
			{
				Name: "VpnNodeRegister",
				InputParameterList: []udwRpc2Builder.RpcParameter{
					{
						Type: udwRpc2Builder.RpcType{
							Kind:       udwRpc2Builder.RpcTypeKindNamedStruct,
							StructName: "PeerIp",
							GoPkg:      "github.com/tachyon-protocol/udw/udwRpc2",
						},
					},
					{
						Type: udwRpc2Builder.RpcType{
							Kind:       udwRpc2Builder.RpcTypeKindNamedStruct,
							StructName: "VpnNode",
							GoPkg:      "github.com/tachyon-protocol/udw/tyVpnRouteServer/tyVpnRouteClient",
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
			{
				Name: "VpnNodeList",
				OutputParameterList: []udwRpc2Builder.RpcParameter{
					{
						Type: udwRpc2Builder.RpcType{
							Kind: udwRpc2Builder.RpcTypeKindSlice,
							Elem: &udwRpc2Builder.RpcType{
								Kind:       udwRpc2Builder.RpcTypeKindNamedStruct,
								StructName: "VpnNode",
								GoPkg:      "github.com/tachyon-protocol/udw/tyVpnRouteServer/tyVpnRouteClient",
							},
						},
					},
				},
			},
			{
				Name: "Ping",
			},
		},
	}
}
