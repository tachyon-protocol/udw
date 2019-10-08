package udwTask

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
	"time"
)

func TestTimeoutRank(t *testing.T) {
	manager := NewTimeoutRankManager(NewTimeoutRankManagerRequest{
		Timeout:      time.Second * 5,
		ThreadNumber: 10,
		ResultLimit:  3,
	})
	manager.Add(func() interface{} {
		return 0
	}, func() {
		fmt.Println("close 0")
	})
	manager.Add(func() interface{} {
		time.Sleep(time.Second)
		return 1
	}, func() {
		fmt.Println("close 1")
	})
	manager.Add(func() interface{} {
		time.Sleep(time.Second * 2)
		return 2
	}, func() {
		fmt.Println("close 2")
	})
	manager.Add(func() interface{} {
		time.Sleep(time.Second * 4)
		return 1
	}, func() {
		fmt.Println("close 3")
	})
	resultList, err := manager.Run()
	udwTest.Ok(err == nil)
	udwTest.Ok(len(resultList) == 3)
	udwTest.Ok(resultList[0] == 0)
	udwTest.Ok(resultList[1] == 1)
	udwTest.Ok(resultList[2] == 2)
}
