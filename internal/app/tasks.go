package app

import (
	"net/http"
	"strconv"
	"time"

	"github.com/seanmorton/todo-htmx/internal/domain"
	"github.com/seanmorton/todo-htmx/internal/templates"
)

func getTaskFilters(r *http.Request) map[string]any {
	filters := map[string]any{
		"completed_at": nil,
	}
	projectId := r.URL.Query().Get("project_id")
	if projectId != "" {
		filters["project_id"], _ = strconv.ParseInt(projectId, 10, 64)
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

	s.hxRender(w, r, templates.TaskList(tasks))
	return nil
}

func (s *Server) newTask(w http.ResponseWriter, r *http.Request) *httpErr {
	projects, err := s.db.ListProjects()
	if err != nil {
		return &httpErr{"failed getting projects", 500, err}
	}

	s.hxRender(w, r, templates.TaskForm(domain.Task{}, projects))
	return nil
}

func (s *Server) getTask(w http.ResponseWriter, r *http.Request) *httpErr {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return &httpErr{"invalid id", 400, nil}
	}

	task, err := s.db.GetTask(id)
	if err != nil {
		return &httpErr{"failed getting task", 500, err}
	}
	if task == nil {
		return &httpErr{"task not found", 404, nil}
	}

	projects, _ := s.db.ListProjects()
	if err != nil {
		return &httpErr{"failed getting projects", 500, err}
	}

	s.hxRender(w, r, templates.TaskForm(*task, projects))
	return nil
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) *httpErr {
	task := domain.Task{}
	s.applyTaskReq(&task, r)

	task, err := s.db.CreateTask(task)
	if err != nil {
		return &httpErr{"failed creating task", 500, err}
	}

	s.hxRender(w, r, templates.TaskRow(task))
	return nil
}

func (s *Server) updateTask(w http.ResponseWriter, r *http.Request) *httpErr {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return &httpErr{"invalid id", 400, nil}
	}

	task := domain.Task{Id: id}
	s.applyTaskReq(&task, r)

	res, err := s.db.UpdateTask(task)
	if err != nil {
		return &httpErr{"failed updating task", 500, err}
	}
	if res == nil {
		return &httpErr{"task not found", 404, nil}
	}

	s.hxRender(w, r, templates.TaskRow(task))
	return nil
}

func (s *Server) completeTask(w http.ResponseWriter, r *http.Request) *httpErr {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return &httpErr{"invalid id", 400, nil}
	}
	task, err := s.db.GetTask(id)
	if err != nil {
		return &httpErr{"failed getting task", 500, err}
	}
	if task == nil {
		return &httpErr{"task not found", 404, nil}
	}

	now := time.Now()
	task.CompletedAt = &now
	_, err = s.db.UpdateTask(*task)
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

func (s *Server) applyTaskReq(task *domain.Task, r *http.Request) {
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

	task.Title = r.FormValue("title")

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

	// TODO add input support for recur policy
	//recurPolicy := domain.RecurPolicy{
	//	Type: domain.RPDaysAfterComplete,
	//	N:    30,
	//}
	//task.RecurPolicy, _ = json.Marshal(recurPolicy)
}
