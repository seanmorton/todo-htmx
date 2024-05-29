package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/seanmorton/todo-htmx/internal/templates"
)

type Server struct {
	tasksDB TasksDB
	tz      *time.Location
}

// TODO error handling
// TODO SSE for new tasks
func NewServer(taskDB TasksDB, tz *time.Location) Server {
	return Server{tasksDB: taskDB, tz: tz}
}

func (s *Server) RegisterRoutes() {
	http.HandleFunc("/", s.index)
	http.HandleFunc("GET /tasks", s.listTasks)
	http.HandleFunc("GET /tasks/new", s.newTask)
	http.HandleFunc("GET /tasks/{id}", s.getTask)
	http.HandleFunc("POST /tasks", s.createTask)
	http.HandleFunc("POST /tasks/{id}/complete", s.completeTask)
	http.HandleFunc("PUT /tasks/{id}", s.updateTask)
	http.HandleFunc("DELETE /tasks/{id}", s.deleteTask)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tasks", http.StatusTemporaryRedirect)
}

func (s *Server) hxRender(w http.ResponseWriter, r *http.Request, content templ.Component) {
	if r.Header.Get("Hx-Request") == "true" {
		content.Render(r.Context(), w)
	} else {
		templates.Index(content).Render(r.Context(), w)
	}
}

func (s *Server) hxRedirect(w http.ResponseWriter, r *http.Request, location string) {
	if r.Header.Get("Hx-Request") == "true" {
		w.Header().Add("HX-Location", fmt.Sprintf(`{"path":"%s", "target":"main"}`, location))
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, location, http.StatusSeeOther)
	}
}
