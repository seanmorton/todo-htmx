package serializers

import (
	"encoding/json"
	"net/http"

	"github.com/seanmorton/todo-htmx/internal/domain"
)

func ParseTask(task *domain.Task, r *http.Request) error {
	var errs []string

	task.Title = parseString(r, "title", &errs)
	task.ProjectId = parseInt64(r, "projectId", &errs)
	task.AssigneeId = parseOptInt64(r, "assigneeId", &errs)
	task.Description = parseOptString(r, "description")
	task.DueDate = parseOptDate(r, "dueDate", &errs)
	task.RecurPolicy = parseRecurPolicy(r, &errs)

	return validationErr(errs)
}

func parseRecurPolicy(r *http.Request, errs *[]string) []byte {
	policyType := parseOptString(r, "recurPolicyType")
	n := parseOptInt64(r, "recurPolicyN", errs)
	if policyType == nil || n == nil {
		return nil
	}

	if *n < 1 {
		*errs = append(*errs, "days must be greater than 0")
	}
	if *policyType == domain.RPDayOfMonth && *n > 28 {
		*errs = append(*errs, "day of month cannot be greater than 28")
	}
	if *policyType == domain.RPCalDate {
		mm := *n / 100
		dd := int(*n % 100)
		if mm < 1 || mm > 12 {
			*errs = append(*errs, "MM must be within range 1-12")
		}
		if dd < 1 || dd > 28 {
			*errs = append(*errs, "DD must be within range 1-28")
		}
	}
	if len(*errs) > 0 {
		return nil
	}

	data, _ := json.Marshal(domain.RecurPolicy{Type: *policyType, N: *n})
	return data
}
