package recurring

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestNext(t *testing.T) {

	// yyyy/mm/dd
	layout := "2006/01/02"

	tests := []struct {
		name     string
		expr     TemporalExpression
		input    string
		expected string
	}{
		{
			name:     "Day/After",
			expr:     Day(2),
			input:    "2012/12/01",
			expected: "2012/12/02",
		},
		{
			name:     "Day/Same",
			expr:     Day(1),
			input:    "2012/12/01",
			expected: "2012/12/01",
		},
		{
			name:     "Day/Before",
			expr:     Day(1),
			input:    "2012/12/02",
			expected: "2013/01/01",
		},
		{
			name:     "Day/Rare",
			expr:     Day(31),
			input:    "2018/09/30",
			expected: "2018/10/31",
		},
		{
			name:     "Day/Negative",
			expr:     Day(-2),
			input:    "2018/09/30",
			expected: "2018/10/30",
		},
		{
			name:     "DayRange/After",
			expr:     DayRange(2, 5),
			input:    "2012/01/01",
			expected: "2012/01/02",
		},
		{
			name:     "DayRange/Same",
			expr:     DayRange(2, 5),
			input:    "2012/01/03",
			expected: "2012/01/03",
		},
		{
			name:     "DayRange/Before",
			expr:     DayRange(2, 5),
			input:    "2012/01/06",
			expected: "2012/02/02",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := time.Parse(layout, tt.input)
			assert.NilError(t, err)
			actual := tt.expr.Next(input, input.AddDate(1, 0, 0))
			assert.Assert(t, !actual.IsZero())
			assert.Equal(t, actual.Format(layout), tt.expected)
		})
	}
}

func TestIncludes(t *testing.T) {

	// yyyy/mm/dd
	layout := "2006/01/02"

	tests := []struct {
		name    string
		expr    TemporalExpression
		match   []string
		nomatch []string
	}{
		{
			name:    "Day",
			expr:    Day(1),
			match:   []string{"2012/01/01", "2012/12/01", "2014/02/01"},
			nomatch: []string{"2012/01/02", "2016/07/20", "2018/08/08"},
		},
		{
			name:    "Negative Day",
			expr:    Day(-1),
			match:   []string{"2018/10/31", "2018/11/30", "2018/02/28"},
			nomatch: []string{"2018/10/30", "2018/11/28", "2018/02/13"},
		},
		{
			name:    "DayRange",
			expr:    DayRange(5, 7),
			match:   []string{"2018/10/05", "2018/10/06", "2018/10/07"},
			nomatch: []string{"2018/10/04", "2018/10/08"},
		},
		{
			name:    "DayRange Same Day",
			expr:    DayRange(1, 1),
			match:   []string{"2018/10/01", "2018/12/01", "2020/01/01"},
			nomatch: []string{"2018/10/07", "2018/12/10", "2020/01/18"},
		},
		{
			name:  "DayRange Whole Month",
			expr:  DayRange(0, -1),
			match: []string{"2018/10/31", "2018/11/30", "2018/02/28"},
		},
		{
			name:    "Week",
			expr:    Week(1),
			match:   []string{"2018/10/01", "2018/10/03", "2018/10/07"},
			nomatch: []string{"2018/10/08", "2018/10/20"},
		},
		{
			name:    "Weekday",
			expr:    Tuesday,
			match:   []string{"2018/10/02", "2018/10/16"},
			nomatch: []string{"2018/10/01", "2018/10/18"},
		},
		{
			name:    "WeekdayRange",
			expr:    WeekdayRange(time.Tuesday, time.Thursday),
			match:   []string{"2018/10/02", "2018/10/03", "2018/10/04"},
			nomatch: []string{"2018/10/01", "2018/10/05"},
		},
		{
			name:    "Month",
			expr:    October,
			match:   []string{"2018/10/02", "2018/10/03", "2018/10/04"},
			nomatch: []string{"2018/11/02", "2018/12/03", "2018/02/04"},
		},
		{
			name:    "MonthRange",
			expr:    MonthRange(time.January, time.February),
			match:   []string{"2018/01/02", "2018/01/03", "2018/02/04"},
			nomatch: []string{"2018/11/02", "2018/12/03", "2018/03/04"},
		},
		{
			name:    "Year",
			expr:    Year(2018),
			match:   []string{"2018/01/02", "2018/01/03", "2018/02/04"},
			nomatch: []string{"2017/01/02", "2019/01/03", "2012/02/04"},
		},
		{
			name:    "YearRange",
			expr:    YearRange(2016, 2019),
			match:   []string{"2016/01/02", "2017/01/03", "2019/02/04"},
			nomatch: []string{"2015/01/02", "2020/01/03", "2012/02/04"},
		},
		{
			name:    "Or",
			expr:    Or(Year(2012), Day(1)),
			match:   []string{"2012/01/02", "2012/01/12", "2019/02/01"},
			nomatch: []string{"2015/01/02", "2020/01/03", "2013/02/04"},
		},
		{
			name:    "And",
			expr:    And(Year(2012), Day(1)),
			match:   []string{"2012/01/01", "2012/12/01", "2012/04/01"},
			nomatch: []string{"2015/01/01", "2012/01/03", "2013/02/04"},
		},
		{
			name:    "Not",
			expr:    Not(MonthRange(time.January, time.February)),
			match:   []string{"2018/11/02", "2018/12/03", "2018/03/04"},
			nomatch: []string{"2018/01/02", "2018/01/03", "2018/02/04"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, st := range tt.match {
				mt, err := time.Parse(layout, st)
				assert.NilError(t, err)
				assert.Assert(t, tt.expr.Includes(mt))
			}
			for _, st := range tt.nomatch {
				mt, err := time.Parse(layout, st)
				assert.NilError(t, err)
				assert.Assert(t, !tt.expr.Includes(mt))
			}
		})
	}
}
