package udwSingleFlight

type One struct {
	group Group
}

func (one *One) Do(fn func() interface{}) interface{} {
	obj, _ := one.group.Do("one", func() (interface{}, error) {
		return fn(), nil
	})
	return obj
}

func (one *One) DoNoReturn(fn func()) {
	one.group.Do("one", func() (interface{}, error) {
		fn()
		return nil, nil
	})
	return
}
func (one *One) IsRunning() bool {
	return one.group.IsRunning("one")
}
