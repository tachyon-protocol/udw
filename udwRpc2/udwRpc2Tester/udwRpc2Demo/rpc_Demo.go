package udwRpc2Demo

import (
	"github.com/tachyon-protocol/udw/udwRpc2"
)

func Demo_RunServer(addr string) (closer func()) {
	s := Server{}
	sh := udwRpc2.NewServerHub(udwRpc2.ServerReq{
		Addr: addr,
		Handler: func(ctx *udwRpc2.ReqCtx) {
			var fnId uint64
			var errMsg string
			fnId, errMsg = ctx.GetReader().ReadUvarint()
			if errMsg != "" {
				return
			}
			panicErrMsg := udwRpc2.PanicToErrMsg(func() {
				switch fnId {
				case 1:
					var tmp_1 string
					errMsg = ctx.GetReader().ReadValue(&tmp_1)
					if errMsg != "" {
						return
					}
					errMsg = ctx.GetReader().ReadArrayEnd()
					if errMsg != "" {
						return
					}
					s.SetName(tmp_1)
					ctx.GetWriter().WriteString("")
					ctx.GetWriter().WriteArrayEnd()
					errMsg = ctx.GetWriter().Flush()
					if errMsg != "" {
						return
					}
				case 2:
					errMsg = ctx.GetReader().ReadArrayEnd()
					if errMsg != "" {
						return
					}
					tmp_2 := s.GetName()
					ctx.GetWriter().WriteString("")
					errMsg = ctx.GetWriter().WriteValue(tmp_2)
					if errMsg != "" {
						return
					}
					ctx.GetWriter().WriteArrayEnd()
					errMsg = ctx.GetWriter().Flush()
					if errMsg != "" {
						return
					}
				case 3:
					errMsg = ctx.GetReader().ReadArrayEnd()
					if errMsg != "" {
						return
					}
					s.IncreaseInt()
					ctx.GetWriter().WriteString("")
					ctx.GetWriter().WriteArrayEnd()
					errMsg = ctx.GetWriter().Flush()
					if errMsg != "" {
						return
					}
				case 4:
					errMsg = ctx.GetReader().ReadArrayEnd()
					if errMsg != "" {
						return
					}
					tmp_3 := s.GetInt()
					ctx.GetWriter().WriteString("")
					errMsg = ctx.GetWriter().WriteValue(tmp_3)
					if errMsg != "" {
						return
					}
					ctx.GetWriter().WriteArrayEnd()
					errMsg = ctx.GetWriter().Flush()
					if errMsg != "" {
						return
					}
				case 5:
					errMsg = ctx.GetReader().ReadArrayEnd()
					if errMsg != "" {
						return
					}
					s.Panic()
					ctx.GetWriter().WriteString("")
					ctx.GetWriter().WriteArrayEnd()
					errMsg = ctx.GetWriter().Flush()
					if errMsg != "" {
						return
					}
				case 6:
					var tmp_4 string
					errMsg = ctx.GetReader().ReadValue(&tmp_4)
					if errMsg != "" {
						return
					}
					var tmp_5 string
					errMsg = ctx.GetReader().ReadValue(&tmp_5)
					if errMsg != "" {
						return
					}
					var tmp_6 string
					errMsg = ctx.GetReader().ReadValue(&tmp_6)
					if errMsg != "" {
						return
					}
					var tmp_7 []Tstruct
					errMsg = ctx.GetReader().ReadValue(&tmp_7)
					if errMsg != "" {
						return
					}
					errMsg = ctx.GetReader().ReadArrayEnd()
					if errMsg != "" {
						return
					}
					tmp_8, tmp_9, tmp_10 := s.FnP(tmp_4, tmp_5, tmp_6, tmp_7)
					ctx.GetWriter().WriteString("")
					errMsg = ctx.GetWriter().WriteValue(tmp_8)
					if errMsg != "" {
						return
					}
					errMsg = ctx.GetWriter().WriteValue(tmp_9)
					if errMsg != "" {
						return
					}
					errMsg = ctx.GetWriter().WriteValue(tmp_10)
					if errMsg != "" {
						return
					}
					ctx.GetWriter().WriteArrayEnd()
					errMsg = ctx.GetWriter().Flush()
					if errMsg != "" {
						return
					}
				case 7:
					var tmp_11 string
					errMsg = ctx.GetReader().ReadValue(&tmp_11)
					if errMsg != "" {
						return
					}
					var tmp_12 string
					errMsg = ctx.GetReader().ReadValue(&tmp_12)
					if errMsg != "" {
						return
					}
					tmp_13 := udwRpc2.PeerIp{ctx.GetPeerIp()}
					errMsg = ctx.GetReader().ReadArrayEnd()
					if errMsg != "" {
						return
					}
					tmp_14 := s.GetPeerIp(tmp_11, tmp_12, tmp_13)
					ctx.GetWriter().WriteString("")
					errMsg = ctx.GetWriter().WriteValue(tmp_14)
					if errMsg != "" {
						return
					}
					ctx.GetWriter().WriteArrayEnd()
					errMsg = ctx.GetWriter().Flush()
					if errMsg != "" {
						return
					}
				default:
				}
			})
			if panicErrMsg != "" {
				ctx.GetWriter().WriteString(panicErrMsg)
				ctx.GetWriter().WriteArrayEnd()
				errMsg = ctx.GetWriter().Flush()
				if errMsg != "" {
					return
				}
			}
		},
	})
	return sh.Close
}
func Demo_NewClient(addr string) *Demo_Client {
	c := udwRpc2.NewClientHub(udwRpc2.ClientReq{
		Addr: addr,
	})
	return &Demo_Client{
		ch: c,
	}
}

type Demo_Client struct {
	ch *udwRpc2.ClientHub
}

func (c *Demo_Client) SetName(fi1 string) (RpcErr *udwRpc2.RpcError) {
	_networkErr := c.ch.RequestCb(func(ctx *udwRpc2.ReqCtx) {
		ctx.GetWriter().WriteUvarint(1)
		ctx.GetWriter().WriteValue(fi1)
		ctx.GetWriter().WriteArrayEnd()
		errMsg := ctx.GetWriter().Flush()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("dehqx82rjj " + errMsg)
			ctx.Close()
			return
		}
		var s string
		errMsg = ctx.GetReader().ReadValue(&s)
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("ehtjkea4re " + errMsg)
			ctx.Close()
			return
		}
		if s != "" {
			RpcErr = udwRpc2.NewOtherError(s)
			ctx.GetReader().ReadArrayEnd()
			return
		}
		errMsg = ctx.GetReader().ReadArrayEnd()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("4b7rug5mf2 " + errMsg)
			ctx.Close()
			return
		}
		RpcErr = nil
		return
	})
	if _networkErr != "" {
		RpcErr = udwRpc2.NewNetworkError("494fehebw6 " + _networkErr)
	}
	return
}
func (c *Demo_Client) GetName() (fo1 string, RpcErr *udwRpc2.RpcError) {
	_networkErr := c.ch.RequestCb(func(ctx *udwRpc2.ReqCtx) {
		ctx.GetWriter().WriteUvarint(2)
		ctx.GetWriter().WriteArrayEnd()
		errMsg := ctx.GetWriter().Flush()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("dehqx82rjj " + errMsg)
			ctx.Close()
			return
		}
		var s string
		errMsg = ctx.GetReader().ReadValue(&s)
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("ehtjkea4re " + errMsg)
			ctx.Close()
			return
		}
		if s != "" {
			RpcErr = udwRpc2.NewOtherError(s)
			ctx.GetReader().ReadArrayEnd()
			return
		}
		errMsg = ctx.GetReader().ReadValue(&fo1)
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("kvkdcgtnk2 " + errMsg)
			ctx.Close()
			return
		}
		errMsg = ctx.GetReader().ReadArrayEnd()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("4b7rug5mf2 " + errMsg)
			ctx.Close()
			return
		}
		RpcErr = nil
		return
	})
	if _networkErr != "" {
		RpcErr = udwRpc2.NewNetworkError("494fehebw6 " + _networkErr)
	}
	return
}
func (c *Demo_Client) IncreaseInt() (RpcErr *udwRpc2.RpcError) {
	_networkErr := c.ch.RequestCb(func(ctx *udwRpc2.ReqCtx) {
		ctx.GetWriter().WriteUvarint(3)
		ctx.GetWriter().WriteArrayEnd()
		errMsg := ctx.GetWriter().Flush()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("dehqx82rjj " + errMsg)
			ctx.Close()
			return
		}
		var s string
		errMsg = ctx.GetReader().ReadValue(&s)
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("ehtjkea4re " + errMsg)
			ctx.Close()
			return
		}
		if s != "" {
			RpcErr = udwRpc2.NewOtherError(s)
			ctx.GetReader().ReadArrayEnd()
			return
		}
		errMsg = ctx.GetReader().ReadArrayEnd()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("4b7rug5mf2 " + errMsg)
			ctx.Close()
			return
		}
		RpcErr = nil
		return
	})
	if _networkErr != "" {
		RpcErr = udwRpc2.NewNetworkError("494fehebw6 " + _networkErr)
	}
	return
}
func (c *Demo_Client) GetInt() (fo1 int, RpcErr *udwRpc2.RpcError) {
	_networkErr := c.ch.RequestCb(func(ctx *udwRpc2.ReqCtx) {
		ctx.GetWriter().WriteUvarint(4)
		ctx.GetWriter().WriteArrayEnd()
		errMsg := ctx.GetWriter().Flush()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("dehqx82rjj " + errMsg)
			ctx.Close()
			return
		}
		var s string
		errMsg = ctx.GetReader().ReadValue(&s)
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("ehtjkea4re " + errMsg)
			ctx.Close()
			return
		}
		if s != "" {
			RpcErr = udwRpc2.NewOtherError(s)
			ctx.GetReader().ReadArrayEnd()
			return
		}
		errMsg = ctx.GetReader().ReadValue(&fo1)
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("kvkdcgtnk2 " + errMsg)
			ctx.Close()
			return
		}
		errMsg = ctx.GetReader().ReadArrayEnd()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("4b7rug5mf2 " + errMsg)
			ctx.Close()
			return
		}
		RpcErr = nil
		return
	})
	if _networkErr != "" {
		RpcErr = udwRpc2.NewNetworkError("494fehebw6 " + _networkErr)
	}
	return
}
func (c *Demo_Client) Panic() (RpcErr *udwRpc2.RpcError) {
	_networkErr := c.ch.RequestCb(func(ctx *udwRpc2.ReqCtx) {
		ctx.GetWriter().WriteUvarint(5)
		ctx.GetWriter().WriteArrayEnd()
		errMsg := ctx.GetWriter().Flush()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("dehqx82rjj " + errMsg)
			ctx.Close()
			return
		}
		var s string
		errMsg = ctx.GetReader().ReadValue(&s)
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("ehtjkea4re " + errMsg)
			ctx.Close()
			return
		}
		if s != "" {
			RpcErr = udwRpc2.NewOtherError(s)
			ctx.GetReader().ReadArrayEnd()
			return
		}
		errMsg = ctx.GetReader().ReadArrayEnd()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("4b7rug5mf2 " + errMsg)
			ctx.Close()
			return
		}
		RpcErr = nil
		return
	})
	if _networkErr != "" {
		RpcErr = udwRpc2.NewNetworkError("494fehebw6 " + _networkErr)
	}
	return
}
func (c *Demo_Client) FnP(fi1 string, fi2 string, fi3 string, fi4 []Tstruct) (fo1 string, fo2 string, fo3 string, RpcErr *udwRpc2.RpcError) {
	_networkErr := c.ch.RequestCb(func(ctx *udwRpc2.ReqCtx) {
		ctx.GetWriter().WriteUvarint(6)
		ctx.GetWriter().WriteValue(fi1)
		ctx.GetWriter().WriteValue(fi2)
		ctx.GetWriter().WriteValue(fi3)
		ctx.GetWriter().WriteValue(fi4)
		ctx.GetWriter().WriteArrayEnd()
		errMsg := ctx.GetWriter().Flush()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("dehqx82rjj " + errMsg)
			ctx.Close()
			return
		}
		var s string
		errMsg = ctx.GetReader().ReadValue(&s)
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("ehtjkea4re " + errMsg)
			ctx.Close()
			return
		}
		if s != "" {
			RpcErr = udwRpc2.NewOtherError(s)
			ctx.GetReader().ReadArrayEnd()
			return
		}
		errMsg = ctx.GetReader().ReadValue(&fo1)
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("kvkdcgtnk2 " + errMsg)
			ctx.Close()
			return
		}
		errMsg = ctx.GetReader().ReadValue(&fo2)
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("kvkdcgtnk2 " + errMsg)
			ctx.Close()
			return
		}
		errMsg = ctx.GetReader().ReadValue(&fo3)
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("kvkdcgtnk2 " + errMsg)
			ctx.Close()
			return
		}
		errMsg = ctx.GetReader().ReadArrayEnd()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("4b7rug5mf2 " + errMsg)
			ctx.Close()
			return
		}
		RpcErr = nil
		return
	})
	if _networkErr != "" {
		RpcErr = udwRpc2.NewNetworkError("494fehebw6 " + _networkErr)
	}
	return
}
func (c *Demo_Client) GetPeerIp(fi1 string, fi2 string) (fo1 string, RpcErr *udwRpc2.RpcError) {
	_networkErr := c.ch.RequestCb(func(ctx *udwRpc2.ReqCtx) {
		ctx.GetWriter().WriteUvarint(7)
		ctx.GetWriter().WriteValue(fi1)
		ctx.GetWriter().WriteValue(fi2)
		ctx.GetWriter().WriteArrayEnd()
		errMsg := ctx.GetWriter().Flush()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("dehqx82rjj " + errMsg)
			ctx.Close()
			return
		}
		var s string
		errMsg = ctx.GetReader().ReadValue(&s)
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("ehtjkea4re " + errMsg)
			ctx.Close()
			return
		}
		if s != "" {
			RpcErr = udwRpc2.NewOtherError(s)
			ctx.GetReader().ReadArrayEnd()
			return
		}
		errMsg = ctx.GetReader().ReadValue(&fo1)
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("kvkdcgtnk2 " + errMsg)
			ctx.Close()
			return
		}
		errMsg = ctx.GetReader().ReadArrayEnd()
		if errMsg != "" {
			RpcErr = udwRpc2.NewNetworkError("4b7rug5mf2 " + errMsg)
			ctx.Close()
			return
		}
		RpcErr = nil
		return
	})
	if _networkErr != "" {
		RpcErr = udwRpc2.NewNetworkError("494fehebw6 " + _networkErr)
	}
	return
}
