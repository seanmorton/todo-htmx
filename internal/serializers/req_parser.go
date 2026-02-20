package serializers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func parseString(r *http.Request, field string, errs *[]string) string {
	val := r.FormValue(field)
	if val == "" {
		*errs = append(*errs, field + " is required")
	}
	return val
}

func parseInt64(r *http.Request, field string, errs *[]string) int64 {
	str := r.FormValue(field)
	if str == "" {
		*errs = append(*errs, field + " is required")
		return 0
	}
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, "invalid " + field)
	}
	return n
}

func parseOptString(r *http.Request, field string) *string {
	val := r.FormValue(field)
	if val == "" {
		return nil
	}
	return &val
}

func parseOptInt64(r *http.Request, field string, errs *[]string) *int64 {
	str := r.FormValue(field)
	if str == "" {
		return nil
	}
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, "invalid " + field)
		return nil
	}
	return &n
}

func parseOptDate(r *http.Request, field string, errs *[]string) *time.Time {
	str := r.FormValue(field)
	if str == "" {
		return nil
	}
	t, err := time.Parse(time.DateOnly, str)
	if err != nil {
		*errs = append(*errs, "invalid " + field)
		return nil
	}
	return &t
}

func validationErr(errs []string) error {
	if len(errs) == 0 {
		return nil
	}
	return errors.New(strings.Join(errs, "; "))
}
