package pkg

import (
	"time"
)

func LocalDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func LocalNowMidnight() time.Time {
	now := time.Now().Local()
	return LocalDate(now.Year(), now.Month(), now.Day())
}

func ParseOptLocalDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}

	t, err := time.ParseInLocation(time.DateOnly, s, time.Local)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func DateStr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.DateOnly)
}

func DateStrShort(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("01/02")
}
