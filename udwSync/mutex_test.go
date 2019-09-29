package udwSync

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestBool_Get(t *testing.T) {
	b := NewBool(false)
	udwTest.Ok(b.Get() == false)
	b.Set(true)
	udwTest.Ok(b.Get() == true)
}

func TestMutexString(t *testing.T) {
	s := NewString("1")
	udwTest.Equal(s.Get(), "1")
	s.Set("2")
	udwTest.Equal(s.Get(), "2")
}

func TestMutexInt(t *testing.T) {
	i := NewInt(1)
	udwTest.Equal(i.Get(), 1)
	i.Set(2)
	udwTest.Equal(i.Get(), 2)
	i.Add(1)
	udwTest.Equal(i.Get(), 3)
	ret := i.AddAndReturnNew(1)
	udwTest.Equal(ret, 4)
}
