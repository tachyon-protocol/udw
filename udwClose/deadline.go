package udwClose

import "time"

func (c *Closer) SetTimeoutFromStart(timeout time.Duration) {
	timeout = timeout - time.Now().Sub(c.startTime)
	if timeout <= 0 {
		c.Close()
		return
	}
	c.fieldLocker.Lock()
	if c.timeoutTimer == nil {
		c.timeoutTimer = time.AfterFunc(timeout, c.Close)
	} else {
		c.timeoutTimer.Reset(timeout)
	}
	c.fieldLocker.Unlock()
}

func (c *Closer) DurationFromStart() time.Duration {
	return time.Now().Sub(c.startTime)
}

func (c *Closer) ClearTimeout() {
	c.fieldLocker.Lock()
	if c.timeoutTimer != nil {
		c.timeoutTimer.Stop()
		c.timeoutTimer = nil
	}
	c.fieldLocker.Unlock()
}
