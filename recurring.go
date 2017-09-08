package recurring

import (
	"time"

	"github.com/icholy/recurring/timeutil"
)

func Next(t time.Time, te TemporalExpression) time.Time {
	day := 24 * time.Hour
	t = timeutil.BeginningOfDay(t)
	for !te.Includes(t) {
		t = t.Add(day)
	}
	return t
}

func NextN(t time.Time, te TemporalExpression, n int) []time.Time {
	ts := make([]time.Time, n)
	for i := 0; i < n; i++ {
		t = Next(t, te)
		ts[i] = t
		t = t.Add(24 * time.Hour)
	}
	return ts
}

type TemporalExpression interface {
	Includes(t time.Time) bool
}

type Day int

func (d Day) normalize(t time.Time) int {
	day := int(d)
	if day < 0 {
		day = timeutil.EndOfMonth(t).Day() + day + 1
	}
	return day
}

func (d Day) Includes(t time.Time) bool {
	return d.normalize(t) == t.Day()
}

func Days(days ...int) TemporalExpression {
	expressions := make([]TemporalExpression, len(days))
	for i, d := range days {
		expressions[i] = Day(d)
	}
	return Or(expressions...)
}

func DayRange(start, end int) DayRangeExpression {
	return DayRangeExpression{start, end}
}

type DayRangeExpression struct {
	Start int
	End   int
}

func (dr DayRangeExpression) Includes(t time.Time) bool {
	d := t.Day()
	return Day(dr.Start).normalize(t) <= d && d <= Day(dr.End).normalize(t)
}

type Week int

func (w Week) normalize(t time.Time) int {
	week := int(w)
	if week < 0 {
		week = timeutil.WeekOfMonth(timeutil.EndOfMonth(t)) + week + 1
	}
	return week
}

func (w Week) Includes(t time.Time) bool {
	return timeutil.WeekOfMonth(t) == w.normalize(t)
}

func Weeks(weeks ...int) TemporalExpression {
	expressions := make([]TemporalExpression, len(weeks))
	for i, w := range weeks {
		expressions[i] = Week(w)
	}
	return Or(expressions...)
}

type Weekday time.Weekday

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func (wd Weekday) Includes(t time.Time) bool {
	return t.Weekday() == time.Weekday(wd)
}

func Weekdays(weekdays ...time.Weekday) TemporalExpression {
	expressions := make([]TemporalExpression, len(weekdays))
	for i, wd := range weekdays {
		expressions[i] = Weekday(wd)
	}
	return Or(expressions...)
}

func WeekdayRange(start, end time.Weekday) WeekdayRangeExpression {
	return WeekdayRangeExpression{start, end}
}

type WeekdayRangeExpression struct {
	Start time.Weekday
	End   time.Weekday
}

func (wr WeekdayRangeExpression) Includes(t time.Time) bool {
	w := t.Weekday()
	return wr.Start <= w && w <= wr.End
}

type Month time.Month

const (
	January Month = 1 + iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

func (m Month) Includes(t time.Time) bool {
	return t.Month() == time.Month(m)
}

func Months(months ...time.Month) TemporalExpression {
	expressions := make([]TemporalExpression, len(months))
	for i, m := range months {
		expressions[i] = Month(m)
	}
	return Or(expressions...)
}

func MonthRange(start, end time.Month) MonthRangeExpression {
	return MonthRangeExpression{start, end}
}

type MonthRangeExpression struct {
	Start time.Month
	End   time.Month
}

func (mr MonthRangeExpression) Includes(t time.Time) bool {
	m := t.Month()
	return mr.Start <= m && m <= mr.End
}

type Year int

func (y Year) Includes(t time.Time) bool {
	return t.Year() == int(y)
}

func Years(years ...int) TemporalExpression {
	expressions := make([]TemporalExpression, len(years))
	for i, y := range years {
		expressions[i] = Year(y)
	}
	return Or(expressions...)
}

func YearRange(start, end int) YearRangeExpression {
	return YearRangeExpression{start, end}
}

type YearRangeExpression struct {
	Start int
	End   int
}

func (yr YearRangeExpression) Includes(t time.Time) bool {
	year := t.Year()
	return yr.Start <= year && year <= yr.End
}

type Date time.Time

func (d Date) Includes(t time.Time) bool {
	y0, m0, d0 := t.Date()
	y1, m1, d1 := time.Time(d).Date()
	return y0 == y1 && m0 == m1 && d0 == d1
}

func Dates(dates ...time.Time) TemporalExpression {
	expressions := make([]TemporalExpression, len(dates))
	for i, d := range dates {
		expressions[i] = Date(d)
	}
	return Or(expressions...)
}

func Or(expressions ...TemporalExpression) OrExpression {
	return OrExpression{expressions}
}

type OrExpression struct {
	expressions []TemporalExpression
}

func (oe *OrExpression) Or(te TemporalExpression) {
	oe.expressions = append(oe.expressions, te)
}

func (oe OrExpression) Includes(t time.Time) bool {
	for _, te := range oe.expressions {
		if te.Includes(t) {
			return true
		}
	}
	return false
}

func And(expressions ...TemporalExpression) AndExpression {
	return AndExpression{expressions}
}

type AndExpression struct {
	expressions []TemporalExpression
}

func (ae *AndExpression) And(te TemporalExpression) {
	ae.expressions = append(ae.expressions, te)
}

func (ae AndExpression) Includes(t time.Time) bool {
	for _, te := range ae.expressions {
		if !te.Includes(t) {
			return false
		}
	}
	return true
}

func Not(te TemporalExpression) NotExpression {
	return NotExpression{te}
}

type NotExpression struct {
	te TemporalExpression
}

func (ne NotExpression) Includes(t time.Time) bool {
	return !ne.te.Includes(t)
}
