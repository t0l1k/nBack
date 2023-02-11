package data

import (
	"time"
)

const (
	Day Period = iota
	Week
	Month
	Year
	All
)

type Period int

func NewPeriod(p *Period) *Period {
	return p
}

func (p *Period) Next() {
	*p++
	if *p > All {
		*p = Day
	}
}

func (p *Period) Prev() {
	*p--
	if *p < Day {
		*p = All
	}
}

func (p Period) Len() int {
	result := 0
	switch p {
	case Day:
		result = 1
	case Week:
		result = 7
	case Month:
		result = 31
	case Year:
		result = 365
	case All:
		result = 0
	}
	return result
}

func (p Period) String() string {
	s := "all"
	switch p {
	case Day:
		s = "day"
	case Week:
		s = "week"
	case Month:
		s = "month"
	case Year:
		s = "year"
	default:
		s = "all"
	}
	return s
}

func NextYear(k int) (string, string, bool) {
	if k < 0 {
		return "", "", false
	}
	firstDt := GetDb().GetFirstDate()
	dt := time.Now().AddDate(-k, 0, 0)
	dt = time.Date(dt.Year(), 1, 1, 0, 0, 0, 0, dt.Location())
	dtTo := time.Date(dt.Year()+1, 1, 1, 0, 0, 0, 0, dt.Location())
	if dtTo.Before(firstDt) {
		return "", "", false
	}
	return dt.Format("2006-01-02"), dtTo.Format("2006-01-02"), true
}

func NextMonth(k int) (string, string, bool) {
	if k < 0 {
		return "", "", false
	}
	firstDt := GetDb().GetFirstDate()
	dt := time.Now().AddDate(0, -k, 0)
	dtFrom := time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, dt.Location())
	dtTo := time.Date(dtFrom.Year(), dtFrom.Month()+1, 1, 0, 0, 0, 0, dtFrom.Location())
	if dtTo.Before(firstDt) {
		return "", "", false
	}
	return dtFrom.Format("2006-01-02"), dtTo.Format("2006-01-02"), true
}

func NextWeek(k int) (string, string, bool) {
	if k < 0 {
		return "", "", false
	}
	firstDt := GetDb().GetFirstDate()
	dt := time.Now().AddDate(0, 0, -k*6)
	mn := k * 7
	for dt.Weekday() != time.Monday {
		dt = time.Now().AddDate(0, 0, -mn)
		mn++
	}
	dtTo := dt.AddDate(0, 0, 7)
	if dtTo.Before(firstDt) {
		return "", "", false
	}
	return dt.Format("2006-01-02"), dtTo.Format("2006-01-02"), true
}

func NextDay(k int) (string, string, bool) {
	if k < 0 {
		return "", "", false
	}
	firstDt := GetDb().GetFirstDate()
	dt := time.Now().AddDate(0, 0, -k)
	if dt.Before(firstDt) {
		return "", "", false
	}
	return dt.Format("2006-01-02"), dt.Format("2006-01-02"), true
}
