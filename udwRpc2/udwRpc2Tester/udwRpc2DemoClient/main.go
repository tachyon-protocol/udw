package main

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwRpc2/udwRpc2Tester/udwRpc2Demo"
	"github.com/tachyon-protocol/udw/udwTest"
)

func main() {
	udwRpc2Demo.RunTest()
	closer := udwRpc2Demo.Demo_RunServer(":8080")
	defer closer()
	c := D2_NewClient("127.0.0.1:8080")
	s, rpcErr := c.GetPeerIp("a1_", "a2_")
	if rpcErr != nil {
		fmt.Println(rpcErr.Error())
	}
	udwTest.Ok(rpcErr == nil)
	udwTest.Equal(s, "a1_a2_127.0.0.1")
	fmt.Println("test pass")
}
