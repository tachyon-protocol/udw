package udwErr

import (
	"bytes"
	"strings"
	"sync"
)

type ErrmsgList struct {
	List []string
}

func (l *ErrmsgList) AddErrMsg(errMsg string) {
	l.List = append(l.List, errMsg)
}

func (l *ErrmsgList) HasErrMsg() bool {
	return len(l.List) > 0
}

func (l *ErrmsgList) GetErrMsg() string {
	if l == nil || len(l.List) == 0 {
		return ""
	}
	buf := bytes.Buffer{}
	for i, errMsg := range l.List {
		buf.WriteString(errMsg)
		if i != len(l.List)-1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func (l *ErrmsgList) AddErrMsgList(errMsgList *ErrmsgList) {
	if len(errMsgList.List) == 0 {
		return
	}
	l.List = append(l.List, errMsgList.List...)
}

type ErrMsgListWithLocker struct {
	list   []string
	locker sync.RWMutex
}

func (l *ErrMsgListWithLocker) Add(errMsg string) {
	if errMsg == "" {
		return
	}
	l.locker.Lock()
	l.list = append(l.list, errMsg)
	l.locker.Unlock()
}

func (l *ErrMsgListWithLocker) Has() bool {
	l.locker.RLock()
	b := len(l.list) > 0
	l.locker.RUnlock()
	return b
}

func (l *ErrMsgListWithLocker) GetAll() string {
	l.locker.RLock()
	if l == nil || len(l.list) == 0 {
		l.locker.RUnlock()
		return ""
	}
	out := strings.Join(l.list, "\n")
	l.locker.RUnlock()
	return out
}
