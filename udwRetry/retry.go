package udwRetry

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwDebug"
	"github.com/tachyon-protocol/udw/udwErr"
	"time"
)

func Run(max int, job func()) error {
	logPrefix := "[udwRetry] "
	allLog := ""
	log := func(s string) {
		s = logPrefix + s
		allLog += s + "\n"
	}
	ok := false
	_max := max
	if _max <= 0 {
		max = 1
	}
	for i := 0; i < max; i++ {
		if _max <= 0 {
			i = -1
		}
		if ok {
			break
		}
		err := udwErr.PanicToErrorWithStackAndLog(job)
		if err == nil {
			ok = true
			break
		} else {
			log(err.Error())
			time.Sleep(time.Second)
		}
	}
	if ok {
		return nil
	}
	log("failed finally!\n" + string(udwDebug.GetCurrentAllStackString(1)))
	return errors.New(allLog)
}

func MustRun(max int, job func()) {
	err := Run(max, job)
	if err != nil {
		panic(err)
	}
}

func AsyncRunInterval(interval time.Duration, job func()) {
	go func() {
		for {
			MustRun(0, job)
			time.Sleep(interval)
		}
	}()
}
