# Recurring

> This is an Implementation of Martin Fowler's [Recurring Events
for Calendars](https://martinfowler.com/apsupp/recurring.pdf)

[![](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/icholy/recurring)


## Example:

``` go

import (
	"time"
	"fmt"
	"github.com/icholy/recurring"
)

func main() {
	
	newyears := recurring.And(
		recurring.January,
		recurring.Day(1)
	)

	halloween := recurring.And(
		recurring.October,
		recurring.Day(31),
	)

	holidays := recurring.Or(newyears, halloween)

	weekends := recurring.Weekdays(
		time.Saturday,
		time.Sunday,
	)

	workdays := recurring.And(
		recurring.Not(weekends),
		recurring.Not(holidays),
	)

	for _, t := range recurring.NextN(t, expr, 10) {
		fmt.Println(t)
	}

}

```
