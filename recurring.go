package recurring

import (
	"time"

	"github.com/icholy/recurring/timeutil"
)

// Find the next occurence of the temporal expression starting at t
func Next(t time.Time, te TemporalExpression) time.Time {
	day := 24 * time.Hour
	t = timeutil.BeginningOfDay(t)
	for !te.Includes(t) {
		t = t.Add(day)
	}
	return t
}

// FindN finds the next n occurences of the temportal expression starting at t
func NextN(t time.Time, te TemporalExpression, n int) []time.Time {
	ts := make([]time.Time, n)
	for i := 0; i < n; i++ {
		t = Next(t, te)
		ts[i] = t
		t = t.Add(24 * time.Hour)
	}
	return ts
}

// TemporalExpression matches a subset of time values
type TemporalExpression interface {

	// Include returns true when the provided time matches the temporal expression
	Includes(t time.Time) bool
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

// Include returns true when provided time's day matches the expressions
func (d Day) Includes(t time.Time) bool {
	return d.normalize(t) == t.Day()
}

// Days is a helper function that combines multiple Day temporal
// expressions with a logical OR operation
func Days(days ...int) TemporalExpression {
	expressions := make([]TemporalExpression, len(days))
	for i, d := range days {
		expressions[i] = Day(d)
	}
	return Or(expressions...)
}

// DayRange returns a temporal expression that matches all
// days between the start and end days
func DayRange(start, end int) DayRangeExpression {
	return DayRangeExpression{start, end}
}

// DayRange is a temporal expression that matches all
// days between the Start and End values
type DayRangeExpression struct {
	Start int
	End   int
}

// Includes returns true when the provided time's day falls
// between the range's Start and Stop values
func (dr DayRangeExpression) Includes(t time.Time) bool {
	d := t.Day()
	return Day(dr.Start).normalize(t) <= d && d <= Day(dr.End).normalize(t)
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

// Weeks is a helper function that combines multiple Week temporal
// expressions with a logical OR operation
func Weeks(weeks ...int) TemporalExpression {
	expressions := make([]TemporalExpression, len(weeks))
	for i, w := range weeks {
		expressions[i] = Week(w)
	}
	return Or(expressions...)
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

// Weekdays is a helper function that combines multiple Weekday
// temporal expressions using a local OR operation
func Weekdays(weekdays ...time.Weekday) TemporalExpression {
	expressions := make([]TemporalExpression, len(weekdays))
	for i, wd := range weekdays {
		expressions[i] = Weekday(wd)
	}
	return Or(expressions...)
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

// Include returns true when the provided time's date
// matches the temporal expression's
func (m Month) Includes(t time.Time) bool {
	return t.Month() == time.Month(m)
}

// Months is a helper function that combines multiple Month temporal
// expressions using a local OR operation
func Months(months ...time.Month) TemporalExpression {
	expressions := make([]TemporalExpression, len(months))
	for i, m := range months {
		expressions[i] = Month(m)
	}
	return Or(expressions...)
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

// Includes returns true when the provided time's month falls
// between the range's Start and Stop values
func (mr MonthRangeExpression) Includes(t time.Time) bool {
	m := t.Month()
	return mr.Start <= m && m <= mr.End
}

// Year is a temporal expression which matchese a year
type Year int

// Include returns true when the provided time's year
// matches the temporal expression's
func (y Year) Includes(t time.Time) bool {
	return t.Year() == int(y)
}

// Years is a helper function that combines multipe Year
// temporal expressions using a local OR operation
func Years(years ...int) TemporalExpression {
	expressions := make([]TemporalExpression, len(years))
	for i, y := range years {
		expressions[i] = Year(y)
	}
	return Or(expressions...)
}

// YearRange returns a temporal expression that matches all
// years between the start and end values
func YearRange(start, end int) YearRangeExpression {
	return YearRangeExpression{start, end}
}

// YearRangeExpression is a temporal expression that matches all
// years between the Start and End values
type YearRangeExpression struct {
	Start int
	End   int
}

// Includes returns true when the provided time's years falls
// between the range's Start and Stop values
func (yr YearRangeExpression) Includes(t time.Time) bool {
	year := t.Year()
	return yr.Start <= year && year <= yr.End
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

// Dates is a helper function that combines multiple Date temporal
// expressions using a logical OR operation
func Dates(dates ...time.Time) TemporalExpression {
	expressions := make([]TemporalExpression, len(dates))
	for i, d := range dates {
		expressions[i] = Date(d)
	}
	return Or(expressions...)
}

// Or combines multiple temporal expressions into one using
// a local Or operation
func Or(expressions ...TemporalExpression) OrExpression {
	return OrExpression{expressions}
}

// OrExpression is a temporal expression consisting of multiple
// temporal expressions combined using a logical OR operation
type OrExpression struct {
	expressions []TemporalExpression
}

// Or adds a temporal expression
func (oe *OrExpression) Or(te TemporalExpression) {
	oe.expressions = append(oe.expressions, te)
}

// Includes returns true when any of the underlying expressions
// match the provided time
func (oe OrExpression) Includes(t time.Time) bool {
	for _, te := range oe.expressions {
		if te.Includes(t) {
			return true
		}
	}
	return false
}

// And combines multiple temporal expressions into one using
// a local AND operation
func And(expressions ...TemporalExpression) AndExpression {
	return AndExpression{expressions}
}

// AndExpression is a temporal expressions consisting of mutliple
// temporal expressions combined with a local AND operation
type AndExpression struct {
	expressions []TemporalExpression
}

// And adds a temporal expression
func (ae *AndExpression) And(te TemporalExpression) {
	ae.expressions = append(ae.expressions, te)
}

// Includes return true when all the underlying temporal expressions
// match the provided time
func (ae AndExpression) Includes(t time.Time) bool {
	for _, te := range ae.expressions {
		if !te.Includes(t) {
			return false
		}
	}
	return true
}

// Not negates a temporal expression
func Not(te TemporalExpression) NotExpression {
	return NotExpression{te}
}

// NotExpression is a temporal expression with negates
// its underlying expression
type NotExpression struct {
	te TemporalExpression
}

// Include returns true when the underlying temporal expression
// does not match the provided time
func (ne NotExpression) Includes(t time.Time) bool {
	return !ne.te.Includes(t)
}
