package testPackage

import (
	"errors"
)

type Demo struct {
}

func (s *Demo) PostScoreInt(LbId string, Score int) (Info string, err error) {
	if Score == 1 {
		return LbId, nil
	} else {
		return "", errors.New("Score!=1")
	}
}
