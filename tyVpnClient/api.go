package tyVpnClient

import (
	"sync"
	"fmt"
	"github.com/tachyon-protocol/udw/udwConsole"
	"github.com/tachyon-protocol/udw/udwLog"
)

func CmdRun(req Config){
	SetConfig(req)
	SetOnChangeCallback("cmd",func(){
		udwLog.Log("SetOnChangeCallback",GetVpnStatus(),GetLastError())
	})
	Connect()
	udwConsole.WaitForExit()
	Disconnect()
}

func Toggle(){
	connectOrDisconnectL2(func()(isConnect bool,isDisconnect bool){
		switch gCsInnerCs {
		case innerCsDisconnected:
			gCsInnerCs = innerCsConnecting
			isConnect = true
		case innerCsConnecting,innerCsConnected,innerCsReconnecting:
			gCsInnerCs = innerCsDisconnectingThenDisconnect
			isDisconnect = true
		case innerCsDisconnectingThenDisconnect:
			gCsInnerCs = innerCsDisconnectingThenConnect
		case innerCsDisconnectingThenConnect:
			gCsInnerCs = innerCsDisconnectingThenDisconnect
		}
		return
	})
}

func Reconnect(){
	connectOrDisconnectL2(func()(isConnect bool,isDisconnect bool) {
		switch gCsInnerCs {
		case innerCsDisconnected:
			gCsInnerCs = innerCsConnecting
			isConnect = true
		case innerCsConnecting,innerCsConnected,innerCsReconnecting:
			gCsInnerCs = innerCsDisconnectingThenConnect
			isDisconnect = true
		case innerCsDisconnectingThenDisconnect,innerCsDisconnectingThenConnect:
			gCsInnerCs = innerCsDisconnectingThenConnect
		}
		return
	})
}

func Connect(){
	connectOrDisconnectL2(func()(isConnect bool,isDisconnect bool) {
		switch gCsInnerCs {
		case innerCsDisconnected:
			gCsInnerCs = innerCsConnecting
			isConnect = true
		case innerCsConnecting,innerCsConnected,innerCsReconnecting,innerCsDisconnectingThenConnect:
		case innerCsDisconnectingThenDisconnect:
			gCsInnerCs = innerCsDisconnectingThenConnect
		}
		return
	})
}

func Disconnect(){
	connectOrDisconnectL2(func()(isConnect bool,isDisconnect bool) {
		switch gCsInnerCs {
		case innerCsConnecting,innerCsConnected,innerCsReconnecting:
			gCsInnerCs = innerCsDisconnectingThenDisconnect
			isDisconnect = true
		case innerCsDisconnectingThenDisconnect,innerCsDisconnected:
		case innerCsDisconnectingThenConnect:
			gCsInnerCs = innerCsDisconnectingThenDisconnect
		}
		return
	})
}

func SetOnChangeCallback(name string,fn func()){
	gOnChangeCallbackListLocker.Lock()
	for i:=range gOnChangeCallbackList{
		elem:=gOnChangeCallbackList[i]
		if elem.name==name{
			gOnChangeCallbackList[i].fn = fn
			gOnChangeCallbackListLocker.Unlock()
			return
		}
	}
	gOnChangeCallbackList = append(gOnChangeCallbackList,onChangeCallbackElem{
		name: name,
		fn: fn,
	})
	gOnChangeCallbackListLocker.Unlock()
	fn()
}

func SetOnChangeCallbackFilterSame(name string,fn func(vpnStatus string,lastErr string)) {
	gLocker:=sync.Mutex{}
	gLastStatus := ""
	SetOnChangeCallback(name,func(){
		vpnStatus:=GetVpnStatus()
		lastError:=GetLastError()
		hasChange:=false
		thisStatusS := vpnStatus+"_"+lastError
		gLocker.Lock()
		if gLastStatus!=thisStatusS{
			gLastStatus = thisStatusS
			hasChange = true
		}
		gLocker.Unlock()
		if hasChange{
			fn(vpnStatus,lastError)
		}
	})
}

const (
	Disconnected = "Disconnected"
	Connecting = "Connecting"
	Connected = "Connected"
	Reconnecting = "Reconnecting"
)

func GetVpnStatus() string{
	gCsLocker.Lock()
	innerCs:=gCsInnerCs
	gCsLocker.Unlock()
	//udwLog.Log("GetVpnStatus",innerCs)
	return getVpnStatusFromInnerCs(innerCs)
}

func GetLastError() string{
	gLastErrorLocker.Lock()
	s:=gLastError
	gLastErrorLocker.Unlock()
	return s
}

func IsInnerCsDisconnected() bool{
	gCsLocker.Lock()
	innerCs:=gCsInnerCs
	gCsLocker.Unlock()
	return innerCs==innerCsDisconnected
}

const (
	innerCsDisconnected = "Disconnected"
	innerCsConnecting = "Connecting"
	innerCsConnected = "Connected"
	innerCsReconnecting = "Reconnecting"
	innerCsDisconnectingThenDisconnect = "DisconnectingThenDisconnect"
	innerCsDisconnectingThenConnect = "DisconnectingThenConnect"
)

var gCsLocker sync.Mutex
var gCsCmdId = 0
var gCsInnerCs = innerCsDisconnected
var gOnChangeCallbackListLocker sync.Mutex
var gOnChangeCallbackList []onChangeCallbackElem
var gLastError = ""
var gLastErrorLocker sync.Mutex

type onChangeCallbackElem struct{
	name string
	fn func()
}

func getVpnStatusFromInnerCs(innerCs string) string{
	switch innerCs {
	case innerCsDisconnected,innerCsDisconnectingThenDisconnect:
		return Disconnected
	case innerCsConnecting,innerCsDisconnectingThenConnect:
		return Connecting
	case innerCsConnected:
		return Connected
	case innerCsReconnecting:
		return Reconnecting
	}
	return Disconnected
}

func fireOnChangeEvent(){
	go func(){
		list:=[]func(){}
		gOnChangeCallbackListLocker.Lock()
		for _,elem:=range gOnChangeCallbackList{
			list = append(list,elem.fn)
		}
		gOnChangeCallbackListLocker.Unlock()
		for _,fn:=range list{
			fn()
		}
	}()
}

func setLastError(errMsg string){
	gLastErrorLocker.Lock()
	gLastError = errMsg
	gLastErrorLocker.Unlock()
	fireOnChangeEvent()
}

func connectOrDisconnectL2(cbFn func()(isConnect bool,isDisconnect bool)){
	thisCsCmdId:=0
	gCsLocker.Lock()
	isConnect,isDisconnect:=cbFn()
	if (isConnect || isDisconnect) {
		gCsCmdId++;
		thisCsCmdId = gCsCmdId;
	}
	gCsLocker.Unlock()
	go func(){
		fireOnChangeEvent()
		connectOrDisconnectL1(isConnect,isDisconnect,thisCsCmdId)
	}()
}

func connectOrDisconnectL1(isConnect bool,isDisconnect bool,thisCsCmdId int){
	if isConnect{
		vpnConnectL1(thisCsCmdId)
	}else if isDisconnect{
		isConnect = false
		vpnDisconnectL1(thisCsCmdId)
		gCsLocker.Lock()
		isVpnStatusChange :=false
		if gCsInnerCs==innerCsDisconnectingThenDisconnect{
			gCsInnerCs = innerCsDisconnected
		}else if gCsInnerCs==innerCsDisconnectingThenConnect{
			gCsInnerCs = innerCsConnecting
			isVpnStatusChange = true
			isConnect = true
			gCsCmdId++;
			thisCsCmdId = gCsCmdId;
		}else{
			fmt.Println("5xkps4fanh "+gCsInnerCs)
		}
		gCsLocker.Unlock()
		if isVpnStatusChange {
			fireOnChangeEvent()
		}
		if isConnect{
			vpnConnectL1(thisCsCmdId)
		}
	}
}

var gClient *Client
var gClientLocker sync.Mutex

func vpnConnectL1(thisCsCmdId int){
	// can be close by vpnDisconnectL1
	c:=&Client{
		thisCsCmdId: thisCsCmdId,
	}
	c.rcInc()
	gClientLocker.Lock()
	if c.isCsCmdIdValid()==false{
		gClientLocker.Unlock()
		return
	}
	gClient = c
	gClientLocker.Unlock()
	c.connectL1(getConfig())
}
func vpnDisconnectL1(thisCsCmdId int){
	//udwLog.Log("vpnDisconnectL1 1",thisCsCmdId)
	// close running vpnConnectL1
	gClientLocker.Lock()
	c:=gClient
	if c==nil{
		gClientLocker.Unlock()
		//udwLog.Log("vpnDisconnectL1 2",thisCsCmdId)
		return
	}
	gClientLocker.Unlock()
	//udwLog.Log("vpnDisconnectL1 3",thisCsCmdId)
	c.closer.Close()
	// wait vpnConnectL1 close finish
	c.rcWg.Wait()
	//udwLog.Log("vpnDisconnectL1 4",thisCsCmdId)
	gClientLocker.Lock()
	gClient = nil
	gClientLocker.Unlock()
}

func (c *Client) isCsCmdIdValid() bool {
	gCsLocker.Lock()
	currentCsCmdId:=gCsCmdId
	gCsLocker.Unlock()
	return currentCsCmdId==c.thisCsCmdId
}

func (c *Client) setInnerCsIfCsCmdIdValid(innerCs string) {
	hasChange:=false
	gCsLocker.Lock()
	if gCsCmdId==c.thisCsCmdId{
		gCsInnerCs = innerCs
		hasChange = true
	}
	gCsLocker.Unlock()
	if hasChange{
		fireOnChangeEvent()
	}
	return
}

func (c *Client) errorDurationConnecting(errMsg string){
	c.setInnerCsIfCsCmdIdValid(innerCsDisconnected)
	setLastError(errMsg)
	c.closer.Close()
}

func (c *Client) rcInc(){
	c.rcLocker.Lock()
	c.rc++
	if c.rc==1{
		c.rcWg.Add(1)
	}
	//udwLog.Log("rcInc",c.rc)
	c.rcLocker.Unlock()
}

func (c *Client) rcDec(){
	c.rcLocker.Lock()
	c.rc--
	if c.rc==0{
		c.rcWg.Done()
	}
	//udwLog.Log("rcDec",c.rc)
	c.rcLocker.Unlock()
}