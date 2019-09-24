package udwTime

import "time"

type MonthDate uint32

func (m MonthDate) IsValid() bool {
	if m.Month() <= 0 || m.Month() > 12 {
		return false
	}
	if m.DayNum() <= 0 {
		return false
	}
	return true
}
func (m MonthDate) Year() int {
	return int(m / 100)
}
func (m MonthDate) Month() time.Month {
	return time.Month(m % 100)
}
func MonthDateFromTime(t time.Time) MonthDate {
	return MonthDate(uint32(t.Year())*100 + uint32(t.Month()))
}

func (m MonthDate) DayNum() int {
	month := m.Month()
	year := m.Year()
	if month == time.February {
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			return 29
		}
		return 28
	}

	if month <= 7 {
		month++
	}
	if month&0x0001 == 0 {
		return 31
	}
	return 30
}
