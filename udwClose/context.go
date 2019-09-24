package udwClose

import "context"

func (c *Closer) GetCloseContext() context.Context {
	ctx2, cancelFunc := context.WithCancel(context.Background())
	c.AddOnClose(cancelFunc)
	return ctx2
}
