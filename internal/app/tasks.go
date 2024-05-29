package app

import (
	"net/http"
	"strconv"
	"time"

	"github.com/seanmorton/todo-htmx/internal/domain"
	"github.com/seanmorton/todo-htmx/internal/templates"
)

func (s *Server) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, _ := s.tasksDB.Query(map[string]any{
		"completed_at": nil,
	})
	s.hxRender(w, r, templates.Tasks(tasks))
}

func (s *Server) getTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	task, _ := s.tasksDB.Get(int64(id))

	s.hxRender(w, r, templates.TaskForm(task))
}

func (s *Server) newTask(w http.ResponseWriter, r *http.Request) {
	s.hxRender(w, r, templates.TaskForm(domain.Task{}))
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	task := domain.Task{}
	task.Title = r.FormValue("title")
	description := r.FormValue("description")
	dueDate := r.FormValue("dueDate")
	if description != "" {
		task.Description = &description
	}

	if dueDate != "" {
		parsed, _ := time.ParseInLocation(time.DateOnly, dueDate, s.tz)
		task.DueDate = &parsed
	}

	// TODO add input support for recur policy
	//recurPolicy := domain.RecurPolicy{
	//	Type: domain.RPDaysAfterComplete,
	//	N:    30,
	//}
	//task.RecurPolicy, _ = json.Marshal(recurPolicy)
	_, _ = s.tasksDB.Create(task)

	s.hxRedirect(w, r, "/tasks")
}

func (s *Server) updateTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	task, _ := s.tasksDB.Get(int64(id))
	task.Title = r.FormValue("title")

	description := r.FormValue("description")
	if description != "" {
		task.Description = &description
	}

	dueDate := r.FormValue("dueDate")
	if dueDate != "" {
		parsed, _ := time.ParseInLocation(time.DateOnly, dueDate, s.tz)
		task.DueDate = &parsed
	}

	// TODO add input support for recur policy
	// recurPolicy := domain.RecurPolicy{
	// 	Type: domain.RPDayOfMonth,
	// 	N:    15,
	// }
	// task.RecurPolicy, _ = json.Marshal(recurPolicy)

	_, _ = s.tasksDB.Update(task)

	s.hxRedirect(w, r, "/tasks")
}

func (s *Server) completeTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	task, _ := s.tasksDB.Get(int64(id))
	now := time.Now()
	task.CompletedAt = &now
	_, _ = s.tasksDB.Update(task)

	if task.RecurPolicy != nil {
		recurTask := domain.Task{
			Title:       task.Title,
			Description: task.Description,
			Assignee:    task.Assignee,
			RecurPolicy: task.RecurPolicy,
		}
		recurTask.DueDate = task.NextRecurDate()
		_, _ = s.tasksDB.Create(recurTask)
	}

	s.hxRedirect(w, r, "/tasks")
}

func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	_ = s.tasksDB.Delete(int64(id))

	s.hxRedirect(w, r, "/tasks")
}
