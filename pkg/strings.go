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
