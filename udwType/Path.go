package udwType

import (
	"strings"
)

type Path []string

func (p Path) String() string {
	if p == nil {
		return ""
	}
	return strings.Join(p, ",")
}
func ParsePath(ps string) Path {
	ps = strings.TrimSpace(ps)
	if ps == "" {
		return nil
	}
	psa := strings.Split(ps, ",")
	pso := []string{}
	for _, v := range psa {
		v = strings.TrimSpace(v)
		pso = append(pso, v)
	}
	return Path(pso)
}
