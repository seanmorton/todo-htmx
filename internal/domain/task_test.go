package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/seanmorton/todo-htmx/pkg"
)

func TestNextRecurDate(t *testing.T) {
	tests := []struct {
		description string
		completedAt time.Time
		policy      *RecurPolicy
		want        *time.Time
	}{
		{
			description: "no recur policy",
			completedAt: time.Now().Local(),
			policy:      nil,
			want:        nil,
		},
		{
			description: "day of month recur policy",
			completedAt: pkg.LocalDate(2026, 1, 1),
			policy:      &RecurPolicy{Type: RPDayOfMonth, N: 5},
			want:        ptr(pkg.LocalDate(2026, 2, 5)),
		},
		{
			description: "day of month recur policy same day",
			completedAt: pkg.LocalDate(2026, 1, 5),
			policy:      &RecurPolicy{Type: RPDayOfMonth, N: 5},
			want:        ptr(pkg.LocalDate(2026, 2, 5)),
		},
		{
			description: "days after recur policy",
			completedAt: pkg.LocalDate(2026, 1, 1),
			policy:      &RecurPolicy{Type: RPDaysAfterComplete, N: 5},
			want:        ptr(pkg.LocalDate(2026, 1, 6)),
		},
		{
			description: "days after recur policy month boundary",
			completedAt: pkg.LocalDate(2026, 1, 31),
			policy:      &RecurPolicy{Type: RPDaysAfterComplete, N: 3},
			want:        ptr(pkg.LocalDate(2026, 2, 3)),
		},
		{
			description: "days after recur policy year boundary",
			completedAt: pkg.LocalDate(2026, 12, 25),
			policy:      &RecurPolicy{Type: RPDaysAfterComplete, N: 90},
			want:        ptr(pkg.LocalDate(2027, 3, 25)),
		},
		{
			description: "cal date recur policy",
			completedAt: pkg.LocalDate(2026, 1, 1),
			policy:      &RecurPolicy{Type: RPCalDate, N: 623},
			want:        ptr(pkg.LocalDate(2026, 6, 23)),
		},
		{
			description: "cal date recur policy next year",
			completedAt: pkg.LocalDate(2026, 6, 23),
			policy:      &RecurPolicy{Type: RPCalDate, N: 623},
			want:        ptr(pkg.LocalDate(2027, 6, 23)),
		},
		{
			description: "cal date recur policy double digit month single digit day",
			completedAt: pkg.LocalDate(2026, 1, 1),
			policy:      &RecurPolicy{Type: RPCalDate, N: 1101},
			want:        ptr(pkg.LocalDate(2026, 11, 1)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			task := Task{
				CompletedAt: &tt.completedAt,
				RecurPolicy: marshalPolicy(t, tt.policy),
			}
			got := task.NextRecurDate()

			if !timeEqual(tt.want, got) {
				t.Errorf("want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func marshalPolicy(t *testing.T, policy *RecurPolicy) []byte {
	t.Helper()
	if policy == nil {
		return nil
	}

	policyJson, err := json.Marshal(policy)
	if err != nil {
		t.Fatal("Failed marshalling test recur policy")
	}

	return policyJson
}

func timeEqual(a, b *time.Time) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equal(*b)
}

// TODO(go 1.26) replace with new()
func ptr[T any](v T) *T { return &v }
