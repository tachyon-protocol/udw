package udwSingleFlight

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwErr"
	"sync"
)

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu        sync.Mutex
	m         map[string]*call
	isRunning bool
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.isRunning = true
	g.mu.Unlock()
	errS := udwErr.PanicToErrorMsgWithStack(func() {
		c.val, c.err = fn()
	})
	if errS != "" {
		c.err = errors.New(errS)
	}
	c.wg.Done()
	g.mu.Lock()
	g.isRunning = false
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}

func (g *Group) IsRunning(key string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.m == nil {
		return false
	}
	if g.m[key] == nil {
		return false
	} else {
		return true
	}
}
