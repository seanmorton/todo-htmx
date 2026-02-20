package serializers

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/seanmorton/todo-htmx/internal/domain"
)

func TestParseTask(t *testing.T) {
	tests := []struct {
		name      string
		form      url.Values
		expectErr bool
		check     func(t *testing.T, task domain.Task)
	}{
		{
			name: "all fields",
			form: url.Values{
				"title":       {"Buy groceries"},
				"projectId":   {"5"},
				"assigneeId":  {"10"},
				"description": {"milk and eggs"},
				"dueDate":     {"2025-04-01"},
			},
			check: func(t *testing.T, task domain.Task) {
				if task.Title != "Buy groceries" {
					t.Errorf("Title = %q", task.Title)
				}
				if task.ProjectId != 5 {
					t.Errorf("ProjectId = %d", task.ProjectId)
				}
				if task.AssigneeId == nil || *task.AssigneeId != 10 {
					t.Errorf("AssigneeId = %v", task.AssigneeId)
				}
				if task.Description == nil || *task.Description != "milk and eggs" {
					t.Errorf("Description = %v", task.Description)
				}
				if task.DueDate == nil || task.DueDate.Format("2006-01-02") != "2025-04-01" {
					t.Errorf("DueDate = %v", task.DueDate)
				}
			},
		},
		{
			name: "required only",
			form: url.Values{
				"title":     {"Do laundry"},
				"projectId": {"1"},
			},
			check: func(t *testing.T, task domain.Task) {
				if task.Title != "Do laundry" {
					t.Errorf("Title = %q", task.Title)
				}
				if task.ProjectId != 1 {
					t.Errorf("ProjectId = %d", task.ProjectId)
				}
				if task.AssigneeId != nil {
					t.Errorf("AssigneeId = %v, expected nil", task.AssigneeId)
				}
				if task.Description != nil {
					t.Errorf("Description = %v, expected nil", task.Description)
				}
				if task.DueDate != nil {
					t.Errorf("DueDate = %v, expected nil", task.DueDate)
				}
				if task.RecurPolicy != nil {
					t.Errorf("RecurPolicy = %v, expected nil", task.RecurPolicy)
				}
			},
		},
		{
			name:    "missing title",
			form:    url.Values{"projectId": {"1"}},
			expectErr: true,
		},
		{
			name:    "missing projectId",
			form:    url.Values{"title": {"Do laundry"}},
			expectErr: true,
		},
		{
			name:    "empty form",
			form:    url.Values{},
			expectErr: true,
		},
		{
			name:    "invalid projectId",
			form:    url.Values{"title": {"X"}, "projectId": {"abc"}},
			expectErr: true,
		},
		{
			name:    "invalid assigneeId",
			form:    url.Values{"title": {"X"}, "projectId": {"1"}, "assigneeId": {"abc"}},
			expectErr: true,
		},
		{
			name:    "invalid dueDate",
			form:    url.Values{"title": {"X"}, "projectId": {"1"}, "dueDate": {"not-a-date"}},
			expectErr: true,
		},
		{
			name: "multiple errors joined",
			form: url.Values{"projectId": {"abc"}, "assigneeId": {"xyz"}},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var task domain.Task
			err := ParseTask(&task, formRequest(tt.form))

			if err == nil && tt.expectErr || err != nil && !tt.expectErr {
				t.Errorf("got err: %v, expected an err: %t", err, tt.expectErr)
			}
			if tt.check != nil {
				tt.check(t, task)
			}
		})
	}
}

func TestParseTask_RecurPolicy(t *testing.T) {
	tests := []struct {
		name       string
		form       url.Values
		expected *domain.RecurPolicy
		expectErr  bool
	}{
		{
			name: "days after complete",
			form: url.Values{
				"title": {"X"}, "projectId": {"1"},
				"recurPolicyType": {domain.RPDaysAfterComplete},
				"recurPolicyN":    {"7"},
			},
			expected: &domain.RecurPolicy{Type: domain.RPDaysAfterComplete, N: 7},
		},
		{
			name: "day of month",
			form: url.Values{
				"title": {"X"}, "projectId": {"1"},
				"recurPolicyType": {domain.RPDayOfMonth},
				"recurPolicyN":    {"15"},
			},
			expected: &domain.RecurPolicy{Type: domain.RPDayOfMonth, N: 15},
		},
		{
			name: "no recurrence fields",
			form: url.Values{
				"title": {"X"}, "projectId": {"1"},
			},
			expected: nil,
		},
		{
			name: "only type no n",
			form: url.Values{
				"title": {"X"}, "projectId": {"1"},
				"recurPolicyType": {domain.RPDaysAfterComplete},
			},
			expected: nil,
		},
		{
			name: "only n no type",
			form: url.Values{
				"title": {"X"}, "projectId": {"1"},
				"recurPolicyN": {"7"},
			},
			expected: nil,
		},
		{
			name: "invalid n",
			form: url.Values{
				"title": {"X"}, "projectId": {"1"},
				"recurPolicyType": {domain.RPDaysAfterComplete},
				"recurPolicyN":    {"abc"},
			},
			expectErr: true,
		},
		{
			name: "n less than 1",
			form: url.Values{
				"title": {"X"}, "projectId": {"1"},
				"recurPolicyType": {domain.RPDaysAfterComplete},
				"recurPolicyN":    {"0"},
			},
			expectErr: true,
		},
		{
			name: "day of month greater than 28",
			form: url.Values{
				"title": {"X"}, "projectId": {"1"},
				"recurPolicyType": {domain.RPDayOfMonth},
				"recurPolicyN":    {"29"},
			},
			expectErr: true,
		},
		{
			name: "day of month at 28 is ok",
			form: url.Values{
				"title": {"X"}, "projectId": {"1"},
				"recurPolicyType": {domain.RPDayOfMonth},
				"recurPolicyN":    {"28"},
			},
			expected: &domain.RecurPolicy{Type: domain.RPDayOfMonth, N: 28},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var task domain.Task
			err := ParseTask(&task, formRequest(tt.form))

			if err == nil && tt.expectErr || err != nil && !tt.expectErr {
				t.Errorf("got err: %v, expected an err: %t", err, tt.expectErr)
			}
			if tt.expected == nil {
				if task.RecurPolicy != nil {
					t.Errorf("RecurPolicy = %s, expected nil", task.RecurPolicy)
				}
				return
			}
			var res domain.RecurPolicy
			if err := json.Unmarshal(task.RecurPolicy, &res); err != nil {
				t.Fatalf("unmarshal RecurPolicy: %v", err)
			}
			if res != *tt.expected {
				t.Errorf("RecurPolicy = %+v, expected %+v", res, *tt.expected)
			}
		})
	}
}
