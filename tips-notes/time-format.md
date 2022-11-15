# Time format

- Input: Time variable.
- Output: Print year, month, and day.

How do you do it? Seems like [time.Format](https://pkg.go.dev/time?utm_source=godoc#Time.Format) can solve it. That's correct, but from the documentation, it seems that any layout can be used: 2006-01-02, 2022-02-03,...

However, interestingly, only layout 2006-01-02 returns the correct answer.

```go

package main

import (
	"fmt"
	"time"
)

func main() {
	tests := []struct {
		comment  string
		layout   string
		date     string
		expected string
	}{
		{
			comment:  "correct layout",
			date:     "2022-11-15T00:00:00Z",
			layout:   "2006-01-02",
			expected: "2022-11-15",
		},
		{
			comment:  "inserts month instead of day",
			date:     "2022-11-15T00:00:00Z",
			layout:   "2006-01-01",
			expected: "2022-11-15",
		},
		{
			comment:  "wrong year",
			date:     "2022-11-15T00:00:00Z",
			layout:   "1999-01-02",
			expected: "2022-11-15",
		},
		{
			comment:  "wrong year",
			date:     "2022-11-15T00:00:00Z",
			layout:   "2002-01-02",
			expected: "2022-11-15",
		},
	}

	for _, test := range tests {
		date, _ := time.Parse(time.RFC3339, test.date)
		fmt.Printf("Comment:  %s\nLayout:   %s \nDate:     %s\nExpected: %s \nActual:   %s\n\n",
			test.comment, test.layout, test.date, test.expected, date.Format(test.layout),
		)
	}
}
```

```shell
Comment:  correct layout
Layout:   2006-01-02
Date:     2022-11-15T00:00:00Z
Expected: 2022-11-15
Actual:   2022-11-15

Comment:  inserts month instead of day
Layout:   2006-01-01
Date:     2022-11-15T00:00:00Z
Expected: 2022-11-15
Actual:   2022-11-11

Comment:  wrong year
Layout:   1999-01-02
Date:     2022-11-15T00:00:00Z
Expected: 2022-11-15
Actual:   11999-11-15

Comment:  wrong year
Layout:   2002-01-02
Date:     2022-11-15T00:00:00Z
Expected: 2022-11-15
Actual:   15319-11-15
```

Quite interesting, right? Dig into [the documentation](https://pkg.go.dev/time?utm_source=godoc#example-Time.Format), I have found this:

```go
// The layout string used by the Parse function and Format method
// shows by example how the reference time should be represented.
// We stress that one must show how the reference time is formatted,
// not a time of the user's choosing. Thus each layout string is a
// representation of the time stamp,
//	Jan 2 15:04:05 2006 MST
// An easy way to remember this value is that it holds, when presented
// in this order, the values (lined up with the elements above):
//	  1 2  3  4  5    6  -7
// There are some wrinkles illustrated below.
```

I think I should accept this layout, and never think about it again.
