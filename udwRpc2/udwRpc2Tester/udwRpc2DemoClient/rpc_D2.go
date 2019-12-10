package main

import (
	"github.com/tachyon-protocol/udw/udwRpc2"
	"github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Tester/udwRpc2Demo"
)

func D2_NewClient(addr string) *D2_Client {
	c := udwRpc2.NewClientHub(udwRpc2.ClientReq{
		Addr: addr,
	})
	return &D2_Client{
		ch: c,
	}
}

type D2_Client struct {
	ch *udwRpc2.ClientHub
}

func (c *D2_Client) SetName(fi1 string) (RpcErr *udwRpc2.RpcError) {
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
func (c *D2_Client) GetName() (fo1 string, RpcErr *udwRpc2.RpcError) {
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
func (c *D2_Client) IncreaseInt() (RpcErr *udwRpc2.RpcError) {
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
func (c *D2_Client) GetInt() (fo1 int, RpcErr *udwRpc2.RpcError) {
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
func (c *D2_Client) Panic() (RpcErr *udwRpc2.RpcError) {
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
func (c *D2_Client) FnP(fi1 string, fi2 string, fi3 string, fi4 []udwRpc2Demo.Tstruct) (fo1 string, fo2 string, fo3 string, RpcErr *udwRpc2.RpcError) {
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
func (c *D2_Client) GetPeerIp(fi1 string, fi2 string) (fo1 string, RpcErr *udwRpc2.RpcError) {
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
