package timeutil

import "time"

func BeginningOfMinute(t time.Time) time.Time {
	return t.Truncate(time.Minute)
}

func BeginningOfHour(t time.Time) time.Time {
	return t.Truncate(time.Hour)
}

func BeginningOfDay(t time.Time) time.Time {
	d := time.Duration(-t.Hour()) * time.Hour
	return BeginningOfHour(t).Add(d)
}

func BeginningOfWeek(t time.Time) time.Time {
	t = BeginningOfDay(t)
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekday = weekday - 1
	d := time.Duration(-weekday) * 24 * time.Hour
	return t.Add(d)
}

func BeginningOfMonth(t time.Time) time.Time {
	t = BeginningOfDay(t)
	d := time.Duration(-int(t.Day())+1) * 24 * time.Hour
	return t.Add(d)
}

func BeginningOfQuarter(t time.Time) time.Time {
	month := BeginningOfMonth(t)
	offset := (int(month.Month()) - 1) % 3
	return month.AddDate(0, -offset, 0)
}

func BeginningOfYear(t time.Time) time.Time {
	t = BeginningOfDay(t)
	d := time.Duration(-int(t.YearDay())+1) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func EndOfMinute(t time.Time) time.Time {
	return BeginningOfMinute(t).Add(time.Minute - time.Nanosecond)
}

func EndOfHour(t time.Time) time.Time {
	return BeginningOfHour(t).Add(time.Hour - time.Nanosecond)
}

func EndOfDay(t time.Time) time.Time {
	return BeginningOfDay(t).Add(24*time.Hour - time.Nanosecond)
}

func EndOfWeek(t time.Time) time.Time {
	return BeginningOfWeek(t).AddDate(0, 0, 7).Add(-time.Nanosecond)
}

func EndOfMonth(t time.Time) time.Time {
	return BeginningOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

func EndOfQuarter(t time.Time) time.Time {
	return BeginningOfQuarter(t).AddDate(0, 3, 0).Add(-time.Nanosecond)
}

func EndOfYear(t time.Time) time.Time {
	return BeginningOfYear(t).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

func Monday(t time.Time) time.Time {
	t = BeginningOfDay(t)
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := time.Duration(-weekday+1) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func Sunday(t time.Time) time.Time {
	t = BeginningOfDay(t)
	weekday := int(t.Weekday())
	if weekday == 0 {
		return t
	} else {
		d := time.Duration(7-weekday) * 24 * time.Hour
		return t.Truncate(time.Hour).Add(d)
	}
}

func EndOfSunday(t time.Time) time.Time {
	return Sunday(t).Add(24*time.Hour - time.Nanosecond)
}

func WeekOfMonth(t time.Time) int {
	_, firstWeek := BeginningOfMonth(t).ISOWeek()
	_, thisWeek := t.ISOWeek()
	return 1 + thisWeek - firstWeek
}
