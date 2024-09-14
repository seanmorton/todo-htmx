package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/seanmorton/todo-htmx/internal/domain"
	"github.com/seanmorton/todo-htmx/internal/templates"
)

func getTaskFilters(r *http.Request) map[string]any {
	filters := map[string]any{}
	projectId := r.URL.Query().Get("projectId")
	if projectId != "" {
		filters["projectId"], _ = strconv.ParseInt(projectId, 10, 64)
	}
	completed := r.URL.Query().Get("completed")
	if completed == "true" {
		filters["completed_at"] = "NOT NULL"
	} else {
		filters["completed_at"] = nil
	}

	return filters
}

func (s *Server) tasks(w http.ResponseWriter, r *http.Request) *httpErr {
	filters := getTaskFilters(r)
	tasks, err := s.db.QueryTasks(filters)
	if err != nil {
		return &httpErr{"failed getting tasks", 500, err}
	}
	projects, err := s.db.ListProjects()
	if err != nil {
		return &httpErr{"failed getting projects", 500, err}
	}

	s.hxRender(w, r, templates.Tasks(tasks, projects, filters))
	return nil
}

func (s *Server) taskList(w http.ResponseWriter, r *http.Request) *httpErr {
	filters := getTaskFilters(r)
	tasks, err := s.db.QueryTasks(filters)
	if err != nil {
		return &httpErr{"failed getting tasks", 500, err}
	}

	s.hxRender(w, r, templates.TaskRows(tasks))
	return nil
}

func (s *Server) newTask(w http.ResponseWriter, r *http.Request) *httpErr {
	task := domain.Task{}
	s.applyTaskReq(&task, r)

	projects, err := s.db.ListProjects()
	if err != nil {
		return &httpErr{"failed getting projects", 500, err}
	}

	s.hxRender(w, r, templates.TaskForm(task, projects))
	return nil
}

func (s *Server) getTask(w http.ResponseWriter, r *http.Request) *httpErr {
	task, retrieveErr := s.retrieveTask(r)
	if retrieveErr != nil {
		return retrieveErr
	}

	projects, err := s.db.ListProjects()
	if err != nil {
		return &httpErr{"failed getting projects", 500, err}
	}

	s.hxRender(w, r, templates.TaskForm(task, projects))
	return nil
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) *httpErr {
	task := domain.Task{}
	validationErr := s.applyTaskReq(&task, r)
	if validationErr != nil {
		return &httpErr{validationErr.Error(), 400, validationErr}
	}

	task, err := s.db.CreateTask(task)
	if err != nil {
		return &httpErr{"failed creating task", 500, err}
	}

	s.hxEvent(w, "taskChange")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) updateTask(w http.ResponseWriter, r *http.Request) *httpErr {
	task, retrieveErr := s.retrieveTask(r)
	if retrieveErr != nil {
		return retrieveErr
	}

	s.applyTaskReq(&task, r)
	res, err := s.db.UpdateTask(task)
	if err != nil {
		return &httpErr{"failed updating task", 500, err}
	}
	if res == nil {
		return &httpErr{"task not found", 404, nil}
	}

	s.hxEvent(w, "taskChange")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) completeTask(w http.ResponseWriter, r *http.Request) *httpErr {
	task, retrieveErr := s.retrieveTask(r)
	if retrieveErr != nil {
		return retrieveErr
	}

	now := time.Now()
	task.CompletedAt = &now
	_, err := s.db.UpdateTask(task)
	if err != nil {
		return &httpErr{"failed completing task", 500, err}
	}

	if task.RecurPolicy != nil {
		recurTask := domain.Task{
			ProjectId:   task.ProjectId,
			AssigneeId:  task.AssigneeId,
			Title:       task.Title,
			Description: task.Description,
			RecurPolicy: task.RecurPolicy,
		}
		recurTask.DueDate = task.NextRecurDate()
		_, err = s.db.CreateTask(recurTask)
		if err != nil {
			return &httpErr{"failed creating next recurring task", 500, err}
		}
	}

	s.hxEvent(w, "taskChange")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) *httpErr {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return &httpErr{"invalid id", 400, nil}
	}

	deleted, err := s.db.DeleteTask(id)
	if err != nil {
		return &httpErr{"failed deleting task", 500, err}
	}
	if !deleted {
		return &httpErr{"task not found", 404, nil}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) retrieveTask(r *http.Request) (domain.Task, *httpErr) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return domain.Task{}, &httpErr{"invalid id", 400, nil}
	}
	task, err := s.db.GetTask(id)
	if err != nil {
		return domain.Task{}, &httpErr{"failed getting task", 500, err}
	}
	if task == nil {
		return domain.Task{}, &httpErr{"task not found", 404, nil}
	}
	return *task, nil
}

func (s *Server) applyTaskReq(task *domain.Task, r *http.Request) error {
	title := r.FormValue("title")
	if title == "" {
		return errors.New("title is required")
	}
	task.Title = title

	projectIdStr := r.FormValue("projectId")
	if projectIdStr != "" {
		projectId, _ := strconv.ParseInt(projectIdStr, 10, 64)
		task.ProjectId = &projectId
	} else {
		task.ProjectId = nil
	}

	assigneeIdStr := r.FormValue("assigneeId")
	if assigneeIdStr != "" {
		assigneeId, _ := strconv.ParseInt(assigneeIdStr, 10, 64)
		task.AssigneeId = &assigneeId
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
		parsed, _ := time.ParseInLocation(time.DateOnly, dueDate, s.tz)
		task.DueDate = &parsed
	} else {
		task.DueDate = nil
	}

	recurPolicyType := r.FormValue("recurPolicyType")
	recurPolicyNStr := r.FormValue("recurPolicyN")
	if recurPolicyType != "" && recurPolicyNStr != "" {
		recurPolicyN, _ := strconv.ParseInt(recurPolicyNStr, 10, 64)
		if recurPolicyN < 1 {
			return errors.New("days must be greater than 0")
		} else if recurPolicyType == domain.RPDayOfMonth && recurPolicyN > 28 {
			return errors.New("day of month cannot be greater than 28")
		}

		recurPolicy := domain.RecurPolicy{
			Type: recurPolicyType,
			N:    recurPolicyN,
		}
		task.RecurPolicy, _ = json.Marshal(recurPolicy)
	} else {
		task.RecurPolicy = nil
	}

	return nil
}
