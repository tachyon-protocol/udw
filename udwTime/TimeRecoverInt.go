package udwTime

import (
	"math"
	"time"
)

type TimeRecoverInt struct {
	Num             int
	Max             int
	LastRecoverTime time.Time
	AddDuration     time.Duration
}

func (t *TimeRecoverInt) Sync(now time.Time) {
	timeDuring := now.Sub(t.LastRecoverTime)
	staminaTimes := float64(timeDuring) / float64(t.AddDuration)
	addStamina := int(math.Floor(staminaTimes))

	if (addStamina + t.Num) >= t.Max {
		t.LastRecoverTime = now
		t.Num = t.Max
		return
	}

	if addStamina < 0 {
		t.LastRecoverTime = now
		return
	}

	if addStamina == 0 {
		return
	}
	t.Num = addStamina + t.Num
	t.LastRecoverTime = t.LastRecoverTime.Add(t.AddDuration * time.Duration(addStamina))
	return
}

func (t *TimeRecoverInt) Full(now time.Time) {
	t.Num = t.Max
	t.LastRecoverTime = now
}
