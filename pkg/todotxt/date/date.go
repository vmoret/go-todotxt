package date

import "time"

const (
	dateLayout = "2006-01-02"
)

// Date represents a date
type Date time.Time

// ZeroDate represents the zero date instant
var ZeroDate Date

// Now returns the current local date.
func Now() Date {
	return Date(time.Now())
}

// Parse parses date from string s.
func Parse(s string) Date {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return ZeroDate
	}
	return Date(t)
}

// IsZero reports whether d represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (d Date) IsZero() bool { return time.Time(d).IsZero() }

func (d Date) String() string {
	if time.Time(d).IsZero() {
		return ""
	}
	return time.Time(d).Format(dateLayout + " ")
}
