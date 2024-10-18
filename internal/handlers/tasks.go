package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/seanmorton/todo-htmx/internal/domain"
	"github.com/seanmorton/todo-htmx/internal/serializers"
	"github.com/seanmorton/todo-htmx/internal/templates"
)

func (s *Server) tasks(w http.ResponseWriter, r *http.Request) *httpErr {
	tasks, params, err := s.fetchTasks(r)
	if err != nil {
		return &httpErr{"failed getting tasks", 500, err}
	}
	projects, err := s.db.ListProjects()
	if err != nil {
		return &httpErr{"failed getting projects", 500, err}
	}
	users, err := s.db.ListUsers()
	if err != nil {
		return &httpErr{"failed getting users", 500, err}
	}

	s.hxRender(w, r, templates.Tasks(tasks, projects, users, params))
	return nil
}

func (s *Server) taskRows(w http.ResponseWriter, r *http.Request) *httpErr {
	tasks, _, err := s.fetchTasks(r)
	if err != nil {
		return &httpErr{"failed getting tasks", 500, err}
	}

	s.hxRender(w, r, templates.TaskRows(tasks))
	return nil
}

func (s *Server) newTask(w http.ResponseWriter, r *http.Request) *httpErr {
	task := domain.Task{}
	serializers.ParseTaskForm(&task, r, s.tz)

	users, err := s.db.ListUsers()
	if err != nil {
		return &httpErr{"failed getting users", 500, err}
	}

	projects, err := s.db.ListProjects()
	if err != nil {
		return &httpErr{"failed getting projects", 500, err}
	}

	s.hxRender(w, r, templates.TaskForm(task, projects, users))
	return nil
}

func (s *Server) getTask(w http.ResponseWriter, r *http.Request) *httpErr {
	task, retrieveErr := s.fetchTask(r)
	if retrieveErr != nil {
		return retrieveErr
	}

	users, err := s.db.ListUsers()
	if err != nil {
		return &httpErr{"failed getting users", 500, err}
	}

	projects, err := s.db.ListProjects()
	if err != nil {
		return &httpErr{"failed getting projects", 500, err}
	}

	s.hxRender(w, r, templates.TaskForm(task, projects, users))
	return nil
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) *httpErr {
	task := domain.Task{}
	validationErr := serializers.ParseTaskForm(&task, r, s.tz)
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
	task, retrieveErr := s.fetchTask(r)
	if retrieveErr != nil {
		return retrieveErr
	}
	validationErr := serializers.ParseTaskForm(&task, r, s.tz)
	if validationErr != nil {
		return &httpErr{validationErr.Error(), 400, validationErr}
	}

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
	task, retrieveErr := s.fetchTask(r)
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

func (s *Server) incompleteTask(w http.ResponseWriter, r *http.Request) *httpErr {
	task, retrieveErr := s.fetchTask(r)
	if retrieveErr != nil {
		return retrieveErr
	}

	task.CompletedAt = nil
	_, err := s.db.UpdateTask(task)
	if err != nil {
		return &httpErr{"failed incompleting task", 500, err}
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

func (s *Server) fetchTasks(r *http.Request) (tasks []domain.Task, params map[string]any, err error) {
	params = map[string]any{}
	projectId := r.FormValue("projectId")
	if projectId != "" {
		params["projectId"], _ = strconv.ParseInt(projectId, 10, 64)
	}
	assigneeId := r.FormValue("assigneeId")
	if assigneeId != "" {
		params["assigneeId"], _ = strconv.ParseInt(assigneeId, 10, 64)
	}
	completed := r.FormValue("completed")
	if completed == "true" {
		params["completed_at"] = "NOT NULL"
	} else {
		params["completed_at"] = nil
	}
	nextMonthOnly := r.FormValue("nextMonthOnly")
	tasks, err = s.db.QueryTasks(params, nextMonthOnly != "false")

	return
}

func (s *Server) fetchTask(r *http.Request) (domain.Task, *httpErr) {
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
