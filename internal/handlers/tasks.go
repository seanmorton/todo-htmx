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
	tasks, filter, err := s.fetchTasks(r)
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

	s.hxRender(w, r, templates.Tasks(tasks, projects, users, filter))
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

	// Preset fields from params, but ignore validation errors since
	// we're just setting up the form here
	serializers.ParseTask(&task, r)

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
	validationErr := serializers.ParseTask(&task, r)
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
	validationErr := serializers.ParseTask(&task, r)
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

	s.hxEvent(w, "taskChange")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) fetchTasks(r *http.Request) ([]domain.Task, domain.TaskFilters, error) {
	filter := domain.TaskFilters{
		Completed:     r.FormValue("completed") == "true",
		NextMonthOnly: r.FormValue("nextMonthOnly") != "false",
	}
	if projectId := r.FormValue("projectId"); projectId != "" {
		id, _ := strconv.ParseInt(projectId, 10, 64)
		filter.ProjectID = &id
	}
	if assigneeId := r.FormValue("assigneeId"); assigneeId != "" {
		id, _ := strconv.ParseInt(assigneeId, 10, 64)
		filter.AssigneeID = &id
	}
	if search := r.FormValue("q"); search != "" {
		filter.Search = &search
	}
	tasks, err := s.db.QueryTasks(filter)
	return tasks, filter, err
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
