package udwTcpPing

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwLog"
	"net"
	"strconv"
	"time"
)

type RunClientRequest struct {
	Ip     string
	Count  int
	Logger func(a ...interface{}) `json:"-"`
}

func RunClient(req RunClientRequest) error {
	if req.Logger == nil {
		req.Logger = udwLog.Log
	}
	conn, err := net.Dial("tcp", req.Ip+":"+strconv.Itoa(port))
	if err != nil {

		return errors.New("[qymys9ajgp] " + err.Error())
	}
	type response struct {
		latency time.Duration
		err     error
	}
	var (
		responsePipe = make(chan response)
	)
	go func() {
		buf := make([]byte, 1)
		for {
			err := conn.SetDeadline(time.Now().Add(timeout))
			if err != nil {

				responsePipe <- response{
					err: errors.New("[6tjmhr3m29] " + err.Error()),
				}
				return
			}
			start := time.Now()
			_, err = conn.Write(buf)
			if err != nil {

				responsePipe <- response{
					err: errors.New("[6t4sbqyruh] write err" + err.Error()),
				}
				return
			}
			_, err = conn.Read(buf)
			if err != nil {
				responsePipe <- response{
					err: errors.New("[2sec9em8cs] read err" + err.Error()),
				}
				return
			}
			latency := time.Now().Sub(start)
			responsePipe <- response{
				latency: latency,
			}
			time.Sleep(interval - latency)
		}
	}()
	var (
		i             = 0
		intervalTimer = time.NewTimer(interval * 2)
	)
	for {
		i++
		if req.Count > 0 && i >= req.Count {
			return nil
		}
		intervalTimer.Reset(interval * 2)
		select {
		case resp := <-responsePipe:
			if resp.err != nil {

				err := errors.New("[p3a3srufqz] TCP connection is broken! " + resp.err.Error())
				req.Logger(err)
				return err
			}
			req.Logger("✔", i, resp.latency)
		case <-intervalTimer.C:
			req.Logger("✘", i)
		}
	}
}

func Ping(ip string) {
	err := RunClient(RunClientRequest{
		Ip: ip,
	})
	udwErr.PanicIfError(err)
}
