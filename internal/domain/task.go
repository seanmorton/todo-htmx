package domain

import (
	"encoding/json"
	"time"
)

type Task struct {
	Id          int64      `json:"id"`
	ProjectId   int64      `json:"ProjectId"`
	AssigneeId  *int64     `json:"assignee_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
	RecurPolicy []byte     `json:"recur_policy"`
	CreatedAt   time.Time  `json:"created_at"`
}

const (
	RPDayOfMonth        = "DoM"
	RPDaysAfterComplete = "DaC"
)

type DueDateFilter string

const (
	DueToday        DueDateFilter = "TODAY"
	Due7Days        DueDateFilter = "7_DAYS"
	Due30Days       DueDateFilter = "30_DAYS"
	MissingDueDate  DueDateFilter = "MISSING"
	NoDueDateFilter DueDateFilter = ""
)

func ParseDueDateFilter(in string) DueDateFilter {
	switch in {
	case string(DueToday), string(Due7Days), string(Due30Days), string(MissingDueDate):
		return DueDateFilter(in)
	default:
		return NoDueDateFilter
	}
}

type TaskFilters struct {
	ProjectID  *int64
	AssigneeID *int64
	Search     *string
	Completed  bool
	DueDate    DueDateFilter
}

type RecurPolicy struct {
	Type string `json:"type"`
	N    int64  `json:"n"`
}

func (t *Task) Done() bool {
	return t.CompletedAt != nil
}

func (t *Task) GetRecurPolicy() *RecurPolicy {
	if len(t.RecurPolicy) == 0 {
		return nil
	}
	rp := RecurPolicy{}
	_ = json.Unmarshal(t.RecurPolicy, &rp)
	return &rp
}

func (t *Task) NextRecurDate() *time.Time {
	if t.CompletedAt == nil {
		return nil
	}
	rp := t.GetRecurPolicy()
	if rp == nil {
		return nil
	}

	var next time.Time
	switch rp.Type {
	case RPDaysAfterComplete:
		next = t.CompletedAt.AddDate(0, 0, int(rp.N))
	case RPDayOfMonth:
		next = t.CompletedAt.AddDate(0, 1, (-t.CompletedAt.Day() + int(rp.N)))
	}
	return &next
}
