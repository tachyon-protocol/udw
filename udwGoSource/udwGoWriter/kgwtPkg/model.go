package kgwtPkg

import (
	"encoding/json"
	"time"
)

type T1 struct {
	Now time.Time
	M   map[string]string
}

func (t1 *T1) String() string {
	t1.Now = t1.Now.UTC()
	b, err := json.Marshal(t1)
	if err != nil {
		panic(err)
	}
	return string(b)
}
