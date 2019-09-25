package udwCmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (c *Cmd) MustSetEnv(key string, value string) *Cmd {
	err := SetCmdEnv(c.cmd, key, value)
	if err != nil {
		panic(err)
	}
	return c
}

func (c *Cmd) MustSetEnvMap(m map[string]string) *Cmd {
	env := &Env{
		values: m,
	}
	c.cmd.Env = env.ToArray()
	return c
}

func (c *Cmd) MustSetEnvMapMergeDefault(m map[string]string) *Cmd {
	if len(c.cmd.Env) == 0 {
		c.cmd.Env = os.Environ()
	}
	env, err := NewEnvFromArray(c.cmd.Env)
	if err != nil {
		panic(err)
	}

	for k, v := range m {
		env.values[k] = v
	}

	c.cmd.Env = env.ToArray()
	return c
}

func SetCmdEnv(cmd *exec.Cmd, key string, value string) error {
	if len(cmd.Env) == 0 {
		cmd.Env = os.Environ()
	}
	env, err := NewEnvFromArray(cmd.Env)
	if err != nil {
		return err
	}
	env.values[key] = value
	cmd.Env = env.ToArray()
	return nil
}

func MustGetEnvMapFromSystem() (m map[string]string) {
	envSList := os.Environ()
	obj, err := NewEnvFromArray(envSList)
	if err != nil {
		panic(err)
	}
	return obj.ToMap()
}

type Env struct {
	values map[string]string
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
