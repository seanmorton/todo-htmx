package pkg

import (
	"time"
)

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

func ParseOptDateStr(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}

	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
