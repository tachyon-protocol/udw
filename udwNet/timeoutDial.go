package udwNet

import (
	"context"
	"errors"
	"github.com/tachyon-protocol/udw/udwClose"
	"net"
	"time"
)

type WrapDialTimeoutReq struct {
	Dial    Dialer
	Timeout time.Duration
	Closer  *udwClose.Closer
}

func WrapDialTimeout(req WrapDialTimeoutReq) Dialer {
	return func(network, address string) (net.Conn, error) {
		if req.Closer == nil {
			req.Closer = udwClose.NewCloser()
		}
		if req.Timeout == 0 {
			req.Timeout = time.Second * 5
		}
		var (
			conn net.Conn
			err  error
		)
		if req.Closer.IsClose() {

			return nil, errors.New("kbja9xqyxr be closed before dial")
		}
		dialDoneChan := make(chan struct{})
		go func() {
			conn, err = req.Dial(network, address)
			close(dialDoneChan)
		}()
		select {
		case <-req.Closer.GetCloseChan():
			return nil, errors.New("377kd87mjm dial be closed")
		case <-time.After(req.Timeout):

			return nil, errors.New("zsrg6aypxz dial timeout")
		case <-dialDoneChan:
		}
		if err != nil {
			return nil, err
		}
		req.Closer.AddOnClose(func() {
			conn.Close()
		})
		if req.Closer.IsClose() {

			return nil, errors.New("rcqt8aawf4")
		}
		return conn, nil
	}
}

func NewRetractableNetDial(closer *udwClose.Closer, timeout time.Duration) Dialer {
	dialer := net.Dialer{
		Timeout: timeout,
	}
	dialCtx, cancelFunc := context.WithCancel(context.Background())
	if closer == nil {
		panic("wv3kyjwg42 closer == nil")
	}
	closer.AddOnClose(cancelFunc)
	return func(network, address string) (net.Conn, error) {
		return dialer.DialContext(dialCtx, network, address)
	}
}

type WrapTcpDialRequest struct {
	ParentDial                  Dialer
	ParentDialTimeout           time.Duration
	ParentDialAlreadyHasTimeout bool

	Closer   *udwClose.Closer
	ConnWrap func(conn net.Conn) (_conn net.Conn, err error)
}

func WrapTcpDial(req WrapTcpDialRequest) Dialer {
	return func(network, address string) (net.Conn, error) {
		if network != "tcp" {
			return nil, errors.New("gghp9ac4yc NewDialer only support tcp, get " + network)
		}
		if req.Closer == nil {
			req.Closer = udwClose.NewCloser()
		}
		if !req.ParentDialAlreadyHasTimeout {
			req.ParentDial = WrapDialTimeout(WrapDialTimeoutReq{
				Dial:    req.ParentDial,
				Timeout: req.ParentDialTimeout,
				Closer:  req.Closer,
			})
		}
		if req.Closer.IsClose() {

			return nil, errors.New("v32ffkjbku")
		}
		conn, err := req.ParentDial(network, address)
		if err != nil {
			return nil, err
		}
		req.Closer.AddOnClose(func() {
			conn.Close()
		})
		if req.Closer.IsClose() {

			return nil, errors.New("ujddt5k9yn")
		}
		return req.ConnWrap(conn)
	}
}
