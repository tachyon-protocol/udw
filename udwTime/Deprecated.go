package udwTime

import "time"

var DefaultNower tDefaultNower

func ParseAutoInLocal(sTime string) (t time.Time, err error) {
	return ParseAutoInLocation(sTime, time.Local)
}

func MustParseAutoInLocal(sTime string) (t time.Time) {
	t, err := ParseAutoInLocation(sTime, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}
