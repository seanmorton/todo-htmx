package domain

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Assignee    *string    `json:"assignee"`
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
	RecurPolicy []byte     `json:"recur_policy"`
	CreatedAt   time.Time  `json:"created_at"`
}

const (
	RPDayOfMonth        = "DoM"
	RPDaysAfterComplete = "DaC"
)

type RecurPolicy struct {
	Type string `json:"type"`
	N    int    `json:"n"`
}

func (t *Task) DueDateStr() string {
	if t.DueDate == nil {
		return ""
	}
	return t.DueDate.Format(time.DateOnly)
}

func (t *Task) Done() bool {
	return t.CompletedAt != nil
}

func (t *Task) NextRecurDate() *time.Time {
	if t.RecurPolicy == nil || t.CompletedAt == nil {
		return nil
	}
	rp := RecurPolicy{}
	_ = json.Unmarshal(t.RecurPolicy, &rp)

	var next time.Time
	switch rp.Type {
	case RPDaysAfterComplete:
		next = t.CompletedAt.AddDate(0, 0, rp.N)
	case RPDayOfMonth:
		next = t.CompletedAt.AddDate(0, 1, (-t.CreatedAt.Day() + rp.N))
	}
	return &next
}
