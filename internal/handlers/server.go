package handlers

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/seanmorton/todo-htmx/internal/data"
	"github.com/seanmorton/todo-htmx/internal/templates"
)

type Server struct {
	db data.DB
	tz *time.Location
}

type httpErr struct {
	Message string
	Code    int
	Cause   error
}

type handler func(http.ResponseWriter, *http.Request) *httpErr

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		slog.Error(err.Message, "code", err.Code, "cause", err.Cause)
		w.WriteHeader(err.Code)
		templates.Error(err.Message).Render(r.Context(), w)
	}
}

func NewServer(db data.DB, tz *time.Location) Server {
	return Server{db: db, tz: tz}
}

func (s *Server) Start(port string, publicDir embed.FS) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.index)

	mux.Handle("/public/", http.FileServer(http.FS(publicDir)))

	mux.Handle("GET /projects", handler(s.projects))
	mux.Handle("GET /projects/rows", handler(s.projectRows))
	mux.Handle("POST /projects", handler(s.createProject))
	mux.Handle("DELETE /projects/{id}", handler(s.deleteProject))

	mux.Handle("GET /tasks", handler(s.tasks))
	mux.Handle("GET /tasks/rows", handler(s.taskRows))
	mux.Handle("GET /tasks/new", handler(s.newTask))
	mux.Handle("GET /tasks/{id}", handler(s.getTask))
	mux.Handle("POST /tasks", handler(s.createTask))
	mux.Handle("POST /tasks/{id}/complete", handler(s.completeTask))
	mux.Handle("POST /tasks/{id}/incomplete", handler(s.incompleteTask))
	mux.Handle("PUT /tasks/{id}", handler(s.updateTask))
	mux.Handle("DELETE /tasks/{id}", handler(s.deleteTask))

	return http.ListenAndServe(port, s.loggingMiddleware(mux))
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

func (s *Server) hxEvent(w http.ResponseWriter, eventName string) {
	w.Header().Set("HX-Trigger", eventName)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(fmt.Sprintf("%s %s", r.Method, r.RequestURI), "hx", (r.Header.Get("HX-Request") == "true"))
		next.ServeHTTP(w, r)
	})
}
