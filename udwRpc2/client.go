package udwRpc2

import (
	"github.com/tachyon-protocol/udw/udwClose"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwShm"
	"net"
	"sync"
	"time"
)

type ClientReq struct {
	Addr           string
	MaxOpenConnNum int
	MaxIdleTime    time.Duration
}
type ClientHub struct {
	req             ClientReq
	maxOpenConnChan chan struct{}
	mu              sync.Mutex
	freeConn        []*Conn
	closer          udwClose.Closer
	stat            ClientStat
}

func NewClientHub(req ClientReq) *ClientHub {
	if req.MaxOpenConnNum <= 0 {
		req.MaxOpenConnNum = 200
	}
	if req.MaxIdleTime <= 0 {
		req.MaxIdleTime = 2 * time.Minute
	}
	ch := &ClientHub{
		req: req,
	}
	ch.maxOpenConnChan = make(chan struct{}, req.MaxOpenConnNum)
	go func() {
		duration := ch.req.MaxIdleTime
		sleepDuration := duration / 2
		ch.closer.LoopUntilCloseFirstSleep(sleepDuration, func() {
			ch.mu.Lock()
			if len(ch.freeConn) == 0 {
				ch.mu.Unlock()
				return
			}
			deadlineTime := time.Now().Add(-duration)
			toCloseList := []*Conn{}
			newFreeList := []*Conn{}
			thisLen := len(ch.freeConn)
			for i := 0; i < thisLen; i++ {
				thisConn := ch.freeConn[i]
				isValid := thisConn.idleStartTime.After(deadlineTime)
				if isValid {
					newFreeList = append(newFreeList, thisConn)
					continue
				}
				toCloseList = append(toCloseList, thisConn)
			}
			if len(toCloseList) > 0 {
				ch.freeConn = newFreeList
			}
			ch.mu.Unlock()
			if len(toCloseList) > 0 {
				for _, conn := range toCloseList {
					conn.Close()
				}
			}
		})
	}()
	return ch
}

type ClientStat struct {
	GetConnNum    int
	NewConnNum    int
	InFreeConnNum int

	IsClose       bool
	CachedStmtNum int
}

func (ch *ClientHub) RequestCb(cb func(ctx *ReqCtx)) (errMsg string) {
	conn, errMsg := ch.getConnWithTimeout(time.Now().Add(time.Minute * 2))
	if errMsg != "" {
		return errMsg
	}
	ctx := &ReqCtx{
		conn: conn,
	}
	ctx.GetWriter().WriteArrayStart()
	errMsg = udwErr.PanicToErrorMsgWithStack(func() {
		cb(ctx)
	})
	if errMsg != "" {
		ctx.Close()
		return errMsg
	}
	if ctx.conn.closer.IsClose() == false {
		ch.putConn(conn)
	}
	return ""
}

func (ch *ClientHub) getConnWithTimeout(deadline time.Time) (mc *Conn, errMsg string) {
	var timeout time.Duration
	select {
	case ch.maxOpenConnChan <- struct{}{}:
	default:
		timeout = deadline.Sub(time.Now())
		select {
		case ch.maxOpenConnChan <- struct{}{}:
		case <-time.After(timeout):
			return nil, "nwatnax5mb i/o timeout"
		}
	}
	ch.mu.Lock()
	ch.stat.GetConnNum++
	if ch.closer.IsClose() {
		ch.mu.Unlock()
		ch.removeInUseConn()
		return nil, "vd8dpyuxx3"
	}
	if len(ch.freeConn) > 0 {
		thisMc := ch.freeConn[len(ch.freeConn)-1]
		ch.freeConn[len(ch.freeConn)-1] = nil
		ch.freeConn = ch.freeConn[:len(ch.freeConn)-1]
		ch.mu.Unlock()
		return thisMc, ""
	}
	ch.mu.Unlock()
	timeout = deadline.Sub(time.Now())
	if timeout <= 0 {
		ch.removeInUseConn()
		return nil, "r67y7ku3kw i/o timeout"
	}
	mc, errMsg = ch.newConnTimeout(timeout)
	if errMsg != "" {
		ch.removeInUseConn()
		return nil, errMsg
	}
	ch.mu.Lock()
	ch.stat.NewConnNum++
	ch.mu.Unlock()
	return mc, ""
}

func (ch *ClientHub) newConnTimeout(timeoutDur time.Duration) (mc *Conn, errMsg string) {
	mc = &Conn{}
	conn, err := net.DialTimeout("tcp", ch.req.Addr, timeoutDur)
	if err != nil {
		return nil, err.Error()
	}
	mc.wb = *udwShm.NewShmWriter(conn, 0)
	mc.rb = *udwShm.NewShmReader(conn, 0)
	mc.conn = conn
	return mc, ""
}

func (db *ClientHub) removeInUseConn() {
	<-db.maxOpenConnChan
}

func (db *ClientHub) putConn(mc *Conn) {
	db.mu.Lock()
	if db.closer.IsClose() {
		db.mu.Unlock()
		mc.Close()
		return
	}
	mc.idleStartTime = time.Now()
	db.freeConn = append(db.freeConn, mc)
	db.mu.Unlock()
	db.removeInUseConn()
}
