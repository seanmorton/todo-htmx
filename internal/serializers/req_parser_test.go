package serializers

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

func formRequest(vals url.Values) *http.Request {
	return &http.Request{Form: vals}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		name      string
		val       string
		expected  string
		expectErr bool
	}{
		{"present", "hello", "hello", false},
		{"missing", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs []string
			r := formRequest(url.Values{"f": {tt.val}})
			res := parseString(r, "f", &errs)

			if len(errs) == 0 && tt.expectErr || len(errs) > 0 && !tt.expectErr {
				t.Errorf("got err: %v, expected an err: %t", validationErr(errs), tt.expectErr )
			}
			if res != tt.expected {
				t.Errorf("got %q, expected %q", res, tt.expected)
			}
		})
	}
}

func TestParseInt64(t *testing.T) {
	tests := []struct {
		name      string
		val       string
		expected  int64
		expectErr bool
	}{
		{"valid", "42", 42, false},
		{"negative", "-1", -1, false},
		{"NaN", "abc", 0, true},
		{"missing", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs []string
			r := formRequest(url.Values{"f": {tt.val}})
			res := parseInt64(r, "f", &errs)

			if len(errs) == 0 && tt.expectErr || len(errs) > 0 && !tt.expectErr {
				t.Errorf("got err: %q, expected an err: %t", validationErr(errs), tt.expectErr)
			}
			if res != tt.expected {
				t.Errorf("got %d, expected %d", res, tt.expected)
			}
		})
	}
}

func TestParseOptString(t *testing.T) {
	tests := []struct {
		name     string
		val      string
		expected *string
	}{
		{"present", "hello", strPtr("hello")},
		{"empty", "", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := formRequest(url.Values{"f": {tt.val}})
			res := parseOptString(r, "f")

			if res == nil && tt.expected != nil {
				t.Errorf("got nil, expected %q", *tt.expected)
			}
			if res != nil && (tt.expected == nil || *res != *tt.expected) {
				t.Errorf("got %q, expected %v", *res, tt.expected)
			}
		})
	}
}

func TestParseOptInt64(t *testing.T) {
	tests := []struct {
		name      string
		val       string
		expected  *int64
		expectErr bool
	}{
		{"valid", "42", int64Ptr(42), false},
		{"NaN", "abc", nil, true},
		{"empty", "", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs []string
			r := formRequest(url.Values{"f": {tt.val}})
			res := parseOptInt64(r, "f", &errs)

			if len(errs) == 0 && tt.expectErr || len(errs) > 0 && !tt.expectErr {
				t.Errorf("got err: %v, expected an err: %t", validationErr(errs), tt.expectErr)
			}
			if res == nil && tt.expected != nil {
				t.Errorf("got nil, expected %d", *tt.expected)
			}
			if res != nil && (tt.expected == nil || *res != *tt.expected) {
				t.Errorf("got %d, expected %v", *res, tt.expected)
			}
		})
	}
}

func TestParseOptDate(t *testing.T) {
	validDate := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		val       string
		expected  *time.Time
		expectErr bool
	}{
		{"valid", "2025-03-15", &validDate, false},
		{"bad format", "03/15/2025", nil, true},
		{"empty", "", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs []string
			r := formRequest(url.Values{"f": {tt.val}})
			res := parseOptDate(r, "f", &errs)

			if len(errs) == 0 && tt.expectErr || len(errs) > 0 && !tt.expectErr {
				t.Errorf("got err: %v, expected an err: %t", validationErr(errs), tt.expectErr)
			}
			if res == nil && tt.expected != nil {
				t.Errorf("got nil, expected %v", *tt.expected)
			}
			if res != nil && (tt.expected == nil || !res.Equal(*tt.expected)) {
				t.Errorf("got %v, expected %v", *res, tt.expected)
			}
		})
	}
}

func strPtr(s string) *string   { return &s }
func int64Ptr(n int64) *int64   { return &n }
