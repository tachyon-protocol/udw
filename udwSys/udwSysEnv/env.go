package udwSysEnv

import (
	"fmt"
	"os"
	"strings"
)

type Env struct {
	values map[string]string
}

func MustNewEnvFromCurrentProcess() (envObj Env) {
	envObj, err := NewEnvFromArray(os.Environ())
	if err != nil {
		panic(err)
	}
	return envObj
}

func NewEnvFromArray(env []string) (envObj Env, err error) {
	envObj = Env{values: make(map[string]string)}
	for _, v1 := range env {
		pos := strings.IndexRune(v1, '=')
		if pos == -1 {
			return Env{}, fmt.Errorf("[NewEnvFromArray] input string not have =, string: %s", v1)
		}
		key := v1[:pos]
		v2 := v1[pos+1:]

		_, ok := envObj.values[key]
		if ok {
			continue
		}
		envObj.values[key] = v2
	}
	return
}

func (env Env) ToArray() []string {
	output := make([]string, 0, len(env.values))
	for k, v := range env.values {
		output = append(output, k+"="+v)
	}
	return output
}

func (env Env) ToMap() map[string]string {
	return env.values
}
