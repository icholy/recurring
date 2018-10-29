package recurring

import (
	"time"

	"github.com/icholy/recurring/timeutil"
)

// Next finds the next occurence of the temporal expression starting at t
func Next(t time.Time, te TemporalExpression) time.Time {
	t = timeutil.BeginningOfDay(t)
	for !te.Includes(t) {
		t = t.Add(24 * time.Hour)
	}
	return t
}

// NextN finds the next n occurences of the temportal expression starting at t
func NextN(t time.Time, te TemporalExpression, n int) []time.Time {
	tt := make([]time.Time, n)
	for i := 0; i < n; i++ {
		t = Next(t, te)
		tt[i] = t
		t = t.Add(24 * time.Hour)
	}
	return tt
}

// TemporalExpression matches a subset of time values
type TemporalExpression interface {

	// Include returns true when the provided time matches the temporal expression
	Includes(t time.Time) bool

	// Next returns the first available time after t that matches the expression
	// if the resulting value is greater than max, return a zero time
	Next(t, max time.Time) time.Time
}

// Day is a temporal expression that matches a day of the month starting at 1
// negative numbers start at the end of the month and move backwards
type Day int

func (d Day) normalize(t time.Time) int {
	day := int(d)
	if day < 0 {
		day = timeutil.EndOfMonth(t).Day() + day + 1
	}
	return day
}

// Includes returns true when provided time's day matches the expressions
func (d Day) Includes(t time.Time) bool {
	return d.normalize(t) == t.Day()
}

// Next returns the first available time after t that matches the expression
// if the resulting value is greater than max, return a zero time
func (d Day) Next(t, max time.Time) time.Time {
	for t.Before(max) {
		if d.Includes(t) {
			return t
		}
		if d.normalize(t) < t.Day() {
			t = timeutil.BeginningOfMonth(t.AddDate(0, 1, 0))
		} else {
			t = timeutil.BeginningOfMonth(t)
		}
		t = t.AddDate(0, 0, d.normalize(t)-1)
	}
	return time.Time{}
}

// Days is a helper function that combines multiple Day temporal
// expressions with a logical OR operation
func Days(days ...int) TemporalExpression {
	ee := make([]TemporalExpression, len(days))
	for i, d := range days {
		ee[i] = Day(d)
	}
	return Or(ee...)
}

// DayRange returns a temporal expression that matches all
// days between the start and end days
func DayRange(start, end int) DayRangeExpression {
	return DayRangeExpression{Day(start), Day(end)}
}

// DayRangeExpression is a temporal expression that matches all
// days between the Start and End values
type DayRangeExpression struct {
	Start Day
	End   Day
}

// Includes returns true when the provided time's day falls
// between the range's Start and Stop values
func (dr DayRangeExpression) Includes(t time.Time) bool {
	d := t.Day()
	return dr.Start.normalize(t) <= d && d <= dr.End.normalize(t)
}

// Next returns the first available time after t that matches the expression
// if the resulting value is greater than max, return a zero time
func (dr DayRangeExpression) Next(t, max time.Time) time.Time {
	for t.Before(max) {
		if dr.Includes(t) {
			return t
		}
		t = dr.Start.Next(t, max)
	}
	return time.Time{}
}

// Week is a temporal expression that matches a week in a month starting at 1
// negative numbers start at the end of the month and move backwards
type Week int

func (w Week) normalize(t time.Time) int {
	week := int(w)
	if week < 0 {
		week = timeutil.WeekOfMonth(timeutil.EndOfMonth(t)) + week + 1
	}
	return week
}

// Includes returns true when the provided time's week matches the expression's
func (w Week) Includes(t time.Time) bool {
	return timeutil.WeekOfMonth(t) == w.normalize(t)
}

// Next returns the first available time after t that matches the expression
// if the resulting value is greater than max, return a zero time
func (w Week) Next(t, max time.Time) time.Time {
	for t.Before(max) {
		if w.Includes(t) {
			return t
		}
		t = timeutil.BeginningOfWeek(t.AddDate(0, 0, 7))
	}
	return time.Time{}
}

// Weeks is a helper function that combines multiple Week temporal
// expressions with a logical OR operation
func Weeks(weeks ...int) TemporalExpression {
	ee := make([]TemporalExpression, len(weeks))
	for i, w := range weeks {
		ee[i] = Week(w)
	}
	return Or(ee...)
}

// Weekday is a temporal expression that matches a day of the week
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

// Includes returns true if the provided time's day of the week
// matches the expression's
func (wd Weekday) Includes(t time.Time) bool {
	return t.Weekday() == time.Weekday(wd)
}

// Next returns the first available time after t that matches the expression
// if the resulting value is greater than max, return a zero time
func (wd Weekday) Next(t, max time.Time) time.Time {
	for t.Before(max) {
		if wd.Includes(t) {
			return t
		}
		t = timeutil.BeginningOfDay(t.AddDate(0, 0, 1))
	}
	return time.Time{}
}

// Weekdays is a helper function that combines multiple Weekday
// temporal expressions using a local OR operation
func Weekdays(weekdays ...time.Weekday) TemporalExpression {
	ee := make([]TemporalExpression, len(weekdays))
	for i, wd := range weekdays {
		ee[i] = Weekday(wd)
	}
	return Or(ee...)
}

// WeekdayRange returns a temporal expression that matches all
// days between the start and end values
func WeekdayRange(start, end time.Weekday) WeekdayRangeExpression {
	return WeekdayRangeExpression{start, end}
}

// WeekdayRangeExpression is a temporal expression that matches all
// days between the Start and End values
type WeekdayRangeExpression struct {
	Start time.Weekday
	End   time.Weekday
}

// Includes returns true when the provided time's weekday falls
// between the range's Start and Stop values
func (wr WeekdayRangeExpression) Includes(t time.Time) bool {
	w := t.Weekday()
	return wr.Start <= w && w <= wr.End
}

// Next returns the first available time after t that matches the expression
// if the resulting value is greater than max, return a zero time
func (wr WeekdayRangeExpression) Next(t, max time.Time) time.Time {
	for t.Before(max) {
		if wr.Includes(t) {
			return t
		}
		t = Weekday(wr.Start).Next(t, max)
	}
	return time.Time{}
}

// Month is a temporal expression which matches a month
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

// Includes returns true when the provided time's date
// matches the temporal expression's
func (m Month) Includes(t time.Time) bool {
	return t.Month() == time.Month(m)
}

// Next returns the first available time after t that matches the expression
// if the resulting value is greater than max, return a zero time
func (m Month) Next(t, max time.Time) time.Time {
	return time.Time{}
}

// Months is a helper function that combines multiple Month temporal
// expressions using a local OR operation
func Months(months ...time.Month) TemporalExpression {
	ee := make([]TemporalExpression, len(months))
	for i, m := range months {
		ee[i] = Month(m)
	}
	return Or(ee...)
}

// MonthRange returns a temporal expression that matches all
// months between the start and end values
func MonthRange(start, end time.Month) MonthRangeExpression {
	return MonthRangeExpression{start, end}
}

// MonthRangeExpression is a temporal expression that matches all
// months between the Start and End values
type MonthRangeExpression struct {
	Start time.Month
	End   time.Month
}

// Next returns the first available time after t that matches the expression
// if the resulting value is greater than max, return a zero time
func (mr MonthRangeExpression) Next(t, max time.Time) time.Time {
	return time.Time{}
}

// Includes returns true when the provided time's month falls
// between the range's Start and Stop values
func (mr MonthRangeExpression) Includes(t time.Time) bool {
	m := t.Month()
	return mr.Start <= m && m <= mr.End
}

// Year is a temporal expression which matchese a year
type Year int

// Includes returns true when the provided time's year
// matches the temporal expression's
func (y Year) Includes(t time.Time) bool {
	return t.Year() == int(y)
}

// Next returns the first available time after t that matches the expression
// if the resulting value is greater than max, return a zero time
func (y Year) Next(t, max time.Time) time.Time {
	return time.Time{}
}

// Years is a helper function that combines multipe Year
// temporal expressions using a local OR operation
func Years(years ...int) TemporalExpression {
	ee := make([]TemporalExpression, len(years))
	for i, y := range years {
		ee[i] = Year(y)
	}
	return Or(ee...)
}

// YearRange returns a temporal expression that matches all
// years between the start and end values
func YearRange(start, end int) YearRangeExpression {
	return YearRangeExpression{Year(start), Year(end)}
}

// YearRangeExpression is a temporal expression that matches all
// years between the Start and End values
type YearRangeExpression struct {
	Start Year
	End   Year
}

// Includes returns true when the provided time's years falls
// between the range's Start and Stop values
func (yr YearRangeExpression) Includes(t time.Time) bool {
	year := t.Year()
	return int(yr.Start) <= year && year <= int(yr.End)
}

// Next returns the first available time after t that matches the expression
// if the resulting value is greater than max, return a zero time
func (yr YearRangeExpression) Next(t, max time.Time) time.Time {
	return time.Time{}
}

// Date is temporal function that matches the year, month, and day
type Date time.Time

// Includes returns true when the provide time's year, month, and
// day match the temporal expression's
func (d Date) Includes(t time.Time) bool {
	y0, m0, d0 := t.Date()
	y1, m1, d1 := time.Time(d).Date()
	return y0 == y1 && m0 == m1 && d0 == d1
}

// Next returns the first available time after t that matches the expression
// if the resulting value is greater than max, return a zero time
func (d Date) Next(t, max time.Time) time.Time {
	return time.Time{}
}

// Dates is a helper function that combines multiple Date temporal
// expressions using a logical OR operation
func Dates(dates ...time.Time) TemporalExpression {
	ee := make([]TemporalExpression, len(dates))
	for i, d := range dates {
		ee[i] = Date(d)
	}
	return Or(ee...)
}

// Or combines multiple temporal expressions into one using
// a local Or operation
func Or(ee ...TemporalExpression) OrExpression {
	return OrExpression{ee}
}

// OrExpression is a temporal expression consisting of multiple
// temporal expressions combined using a logical OR operation
type OrExpression struct {
	ee []TemporalExpression
}

// Or adds a temporal expression
func (oe *OrExpression) Or(e TemporalExpression) {
	oe.ee = append(oe.ee, e)
}

// Includes returns true when any of the underlying expressions
// match the provided time
func (oe OrExpression) Includes(t time.Time) bool {
	for _, e := range oe.ee {
		if e.Includes(t) {
			return true
		}
	}
	return false
}

// Next returns the first available time after t that matches the expression
// if the resulting value is greater than max, return a zero time
func (oe OrExpression) Next(t, max time.Time) time.Time {
	return time.Time{}
}

// And combines multiple temporal expressions into one using
// a local AND operation
func And(ee ...TemporalExpression) AndExpression {
	return AndExpression{ee}
}

// AndExpression is a temporal expressions consisting of mutliple
// temporal expressions combined with a local AND operation
type AndExpression struct {
	ee []TemporalExpression
}

// And adds a temporal expression
func (ae *AndExpression) And(e TemporalExpression) {
	ae.ee = append(ae.ee, e)
}

// Includes return true when all the underlying temporal expressions
// match the provided time
func (ae AndExpression) Includes(t time.Time) bool {
	for _, e := range ae.ee {
		if !e.Includes(t) {
			return false
		}
	}
	return true
}

// Next returns the first available time after t that matches the expression
// if the resulting value is greater than max, return a zero time
func (ae AndExpression) Next(t, max time.Time) time.Time {
	return time.Time{}
}

// Not negates a temporal expression
func Not(e TemporalExpression) NotExpression {
	return NotExpression{e}
}

// NotExpression is a temporal expression with negates
// its underlying expression
type NotExpression struct {
	e TemporalExpression
}

// Includes returns true when the underlying temporal expression
// does not match the provided time
func (ne NotExpression) Includes(t time.Time) bool {
	return !ne.e.Includes(t)
}

// Next returns the first available time after t that matches the expression
// if the resulting value is greater than max, return a zero time
func (ne NotExpression) Next(t, max time.Time) time.Time {
	return time.Time{}
}
