package serializers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seanmorton/todo-htmx/internal/domain"
)

func ParseTaskForm(task *domain.Task, r *http.Request) error {
	var errMessages []string
	title := r.FormValue("title")
	if title == "" {
		errMessages = append(errMessages, "title is required")
	} else {
		task.Title = title
	}

	projectIdStr := r.FormValue("projectId")
	if projectIdStr == "" {
		errMessages = append(errMessages, "project is required")
	} else {
		projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
		if err != nil {
			errMessages = append(errMessages, "invalid project id")
		} else {
			task.ProjectId = projectId
		}
	}

	assigneeIdStr := r.FormValue("assigneeId")
	if assigneeIdStr != "" {
		assigneeId, err := strconv.ParseInt(assigneeIdStr, 10, 64)
		if err != nil {
			errMessages = append(errMessages, "invalid assignee id")
		} else {
			task.AssigneeId = &assigneeId
		}
	} else {
		task.AssigneeId = nil
	}

	description := r.FormValue("description")
	if description != "" {
		task.Description = &description
	} else {
		task.Description = nil
	}

	dueDate := r.FormValue("dueDate")
	if dueDate != "" {
		parsed, err := time.Parse(time.DateOnly, dueDate)
		if err != nil {
			errMessages = append(errMessages, "invalid due date")
		} else {
			task.DueDate = &parsed
		}
	} else {
		task.DueDate = nil
	}

	recurPolicyType := r.FormValue("recurPolicyType")
	recurPolicyNStr := r.FormValue("recurPolicyN")
	if recurPolicyType != "" && recurPolicyNStr != "" {
		recurPolicyN, err := strconv.ParseInt(recurPolicyNStr, 10, 64)
		if err != nil {
			errMessages = append(errMessages, "invalid recurrence value")
		} else if recurPolicyN < 1 {
			errMessages = append(errMessages, "days must be greater than 0")
		} else if recurPolicyType == domain.RPDayOfMonth && recurPolicyN > 28 {
			errMessages = append(errMessages, "day of month cannot be greater than 28")
		} else {
			recurPolicy := domain.RecurPolicy{
				Type: recurPolicyType,
				N:    recurPolicyN,
			}
			task.RecurPolicy, _ = json.Marshal(recurPolicy)
		}
	} else {
		task.RecurPolicy = nil
	}

	if errMessages != nil {
		return errors.New(strings.Join(errMessages, "; "))
	}

	return nil
}
