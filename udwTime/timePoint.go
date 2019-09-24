package udwTime

import "time"

type GetTimePointInBetweenRequest struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

func GetTimePointInBetween(req GetTimePointInBetweenRequest) int {
	if req.EndTime.After(req.StartTime) == false {
		return 0
	}
	if req.Duration == 0 {
		panic("s6v2qq85x6 req.Duration==0")
	}
	TruncateStart := req.StartTime.Truncate(req.Duration)
	TruncateEnd := req.EndTime.Truncate(req.Duration)
	out := int(TruncateEnd.Sub(TruncateStart) / req.Duration)
	return out
}
