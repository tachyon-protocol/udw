package udwRpc2Demo

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwTest"
)

func RunTest() {
	runTest1()
	runTestNoServer()
}

func runTest1() {
	closer := Demo_RunServer(":8080")
	defer closer()
	c := Demo_NewClient("127.0.0.1:8080")
	rpcErr := c.SetName("1234")
	udwTest.Equal(rpcErr, nil)
	udwTest.Equal(Server{}.GetName(), "1234")

	b, rpcErr := c.GetName()
	udwTest.Equal(rpcErr, nil)
	udwTest.Equal(b, "1234")

	rpcErr = c.IncreaseInt()
	udwTest.Equal(rpcErr, nil)

	v, rpcErr := c.GetInt()
	udwTest.Equal(rpcErr, nil)
	udwTest.Equal(v, 1)

	checkFnp(c)

	rpcErr = c.Panic()
	udwTest.Ok(rpcErr != nil)
	udwTest.Ok(rpcErr.IsNetworkError() == false)
	udwTest.Ok(rpcErr.Error() == "jnp5gkkjfy")

	checkFnp(c)

	s, rpcErr := c.GetPeerIp("a1_", "a2_")
	if rpcErr != nil {
		fmt.Println(rpcErr.Error())
	}
	udwTest.Ok(rpcErr == nil)
	udwTest.Equal(s, "a1_a2_127.0.0.1")
}

func checkFnp(c *Demo_Client) {
	o1, o2, o3, rpcErr := c.FnP("a1", "a2", "a3", []Tstruct{{"1"}})
	if rpcErr != nil {
		fmt.Println("checkFnp fail", rpcErr.Error())
	}
	udwTest.Ok(rpcErr == nil)
	udwTest.Equal(o1, "a1")
	udwTest.Equal(o2, "a2")
	udwTest.Equal(o3, "a3_1")
}
func runTestNoServer() {
	c := Demo_NewClient("127.0.0.1:8081")
	_, rpcErr := c.GetName()
	udwTest.Ok(rpcErr != nil)
	udwTest.Ok(rpcErr.IsNetworkError() == true)
}
