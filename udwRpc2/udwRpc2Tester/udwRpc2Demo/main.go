package main

import (
	"github.com/tachyon-protocol/udw/udwRpc2"
	"github.com/tachyon-protocol/udw/udwSync"
)

type Server struct{}

var gName udwSync.String

func (Server) SetName(v string) {
	gName.Set(v)
}

func (Server) GetName() string {
	return gName.Get()
}

var gInt udwSync.Int

func (Server) IncreaseInt() {
	gInt.Add(1)
	return
}

func (Server) GetInt() int {
	return gInt.Get()
}

func (Server) Panic() {
	panic("jnp5gkkjfy")
}

func (Server) FnP(i1 string, i2 string, i3 string) (o1 string, o2 string, o3 string) {
	return i1, i2, i3
}

func (Server) GetPeerIp(p1 string, p2 string, ClientIp udwRpc2.PeerIp) string {
	return p1 + p2 + ClientIp.Ip
}
