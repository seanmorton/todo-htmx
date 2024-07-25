package app

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/seanmorton/todo-htmx/internal/domain"
	"github.com/seanmorton/todo-htmx/internal/templates"
)

func (s *Server) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.db.QueryTasks(map[string]any{
		"completed_at": nil,
	})
	fmt.Println(err)
	s.hxRender(w, r, templates.Tasks(tasks))
}

func (s *Server) getTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	task, _ := s.db.GetTask(id)
	projects, _ := s.db.ListProjects()

	s.hxRender(w, r, templates.TaskForm(task, projects))
}

func (s *Server) newTask(w http.ResponseWriter, r *http.Request) {
	projects, _ := s.db.ListProjects()
	s.hxRender(w, r, templates.TaskForm(domain.Task{}, projects))
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	task := domain.Task{}
	s.applyTaskReq(&task, r)

	_, err := s.db.CreateTask(task)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}

	s.hxRedirect(w, r, "/tasks")
}

func (s *Server) updateTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	task, _ := s.db.GetTask(id)
	s.applyTaskReq(&task, r)

	_, _ = s.db.UpdateTask(task)

	s.hxRedirect(w, r, "/tasks")
}

func (s *Server) completeTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	task, _ := s.db.GetTask(id)
	now := time.Now()
	task.CompletedAt = &now
	_, _ = s.db.UpdateTask(task)

	if task.RecurPolicy != nil {
		recurTask := domain.Task{
			ProjectId:   task.ProjectId,
			AssigneeId:  task.AssigneeId,
			Title:       task.Title,
			Description: task.Description,
			RecurPolicy: task.RecurPolicy,
		}
		recurTask.DueDate = task.NextRecurDate()
		_, _ = s.db.CreateTask(recurTask)
	}

	s.hxRedirect(w, r, "/tasks")
}

func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	_ = s.db.DeleteTask(id)

	s.hxRedirect(w, r, "/tasks")
}

func (s *Server) applyTaskReq(task *domain.Task, r *http.Request) {
	task.Title = r.FormValue("title")

	description := r.FormValue("description")
	if description != "" {
		task.Description = &description
	} else {
		task.Description = nil
	}

	projectIdStr := r.FormValue("projectId")
	if projectIdStr != "" {
		projectId, _ := strconv.ParseInt(projectIdStr, 10, 64)
		task.ProjectId = &projectId
	} else {
		task.ProjectId = nil
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
